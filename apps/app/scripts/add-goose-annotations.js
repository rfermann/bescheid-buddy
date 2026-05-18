import { existsSync, mkdirSync, readdirSync, readFileSync, writeFileSync } from 'fs';
import { join } from 'path';
import { createHash } from 'crypto';

const drizzleDir = './drizzle';
const outputDir = '../../migrations';
const manifestFile = join(outputDir, '.goose-manifest.json');

try {
	mkdirSync(outputDir, { recursive: true });
} catch (e) {}

let manifest = {};
if (existsSync(manifestFile)) {
	try {
		manifest = JSON.parse(readFileSync(manifestFile, 'utf-8'));
	} catch (e) {
		console.warn('Could not read manifest, starting fresh');
	}
}

const files = readdirSync(drizzleDir).filter((f) => f.endsWith('.sql'));

files.sort();

for (let i = 0; i < files.length; i++) {
	const file = files[i];
	const inputPath = join(drizzleDir, file);

	let sql = readFileSync(inputPath, 'utf-8');

	sql = sql.replace(/--> statement-breakpoint/g, '');

	const hash = createHash('sha256').update(sql).digest('hex');

	const existingEntry = manifest[file];
	if (
		existingEntry &&
		existingEntry.hash === hash &&
		existsSync(join(outputDir, existingEntry.gooseFile))
	) {
		console.log(`Skipping ${file} - already processed as ${existingEntry.gooseFile}`);
		continue;
	}

	let migrationName;
	if (i === 0) {
		migrationName = 'init_auth';
	} else {
		migrationName = `update_auth_${i}`;
	}

	const timestamp = new Date()
		.toISOString()
		.replace(/[-:]/g, '')
		.replace(/\.\d+Z$/, '')
		.replace('T', '');
	const outputFileName = `${timestamp}_${migrationName}.sql`;
	const outputPath = join(outputDir, outputFileName);

	const statements = sql.split(';').filter((s) => s.trim());

	const upSql = statements.join(';\n');

	const downSql = generateDownMigration(statements);

	const gooseSql = `-- +goose Up
-- +goose StatementBegin
${upSql};
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
${downSql};
-- +goose StatementEnd
`;

	writeFileSync(outputPath, gooseSql);
	console.log(`Processed ${file} -> ${outputPath}`);

	manifest[file] = {
		hash,
		gooseFile: outputFileName
	};
}

writeFileSync(manifestFile, JSON.stringify(manifest, null, 2));
console.log(`Updated manifest at ${manifestFile}`);

function generateDownMigration(statements) {
	const reversed = [...statements].reverse();

	return reversed
		.map((stmt) => {
			const trimmed = stmt.trim();
			if (trimmed.startsWith('CREATE TABLE')) {
				const tableName = trimmed.match(/CREATE TABLE\s+"?(\w+)"?/i)?.[1];
				return `DROP TABLE IF EXISTS "${tableName}" CASCADE`;
			}
			if (trimmed.startsWith('CREATE INDEX')) {
				const indexName = trimmed.match(/CREATE INDEX\s+"?(\w+)"?/i)?.[1];
				return `DROP INDEX IF EXISTS "${indexName}"`;
			}
			if (trimmed.startsWith('ALTER TABLE')) {
				const tableName = trimmed.match(/ALTER TABLE\s+"?(\w+)"?/i)?.[1];
				const constraintMatch = trimmed.match(/ADD CONSTRAINT\s+"?(\w+)"?/i);
				if (constraintMatch) {
					const constraintName = constraintMatch[1];
					return `ALTER TABLE "${tableName}" DROP CONSTRAINT "${constraintName}"`;
				}
				const columnMatch = trimmed.match(/ADD COLUMN\s+"?(\w+)"?/i);
				if (columnMatch) {
					const columnName = columnMatch[1];
					return `ALTER TABLE "${tableName}" DROP COLUMN "${columnName}"`;
				}
			}
			throw new Error(`Unsupported SQL statement for down migration: ${trimmed}`);
		})
		.join(';\n');
}
