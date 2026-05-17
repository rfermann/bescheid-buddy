import { readFileSync, writeFileSync, readdirSync, existsSync, unlinkSync } from 'fs';
import { join } from 'path';

const drizzleDir = './drizzle';
const outputDir = '../../migrations';

// Ensure output directory exists
import { mkdirSync } from 'fs';
try {
	mkdirSync(outputDir, { recursive: true });
} catch (e) {}

// Read all SQL files from drizzle directory
const files = readdirSync(drizzleDir).filter((f) => f.endsWith('.sql'));

// Sort files to ensure consistent ordering
files.sort();

// Delete all existing auth migration files before creating new ones
if (existsSync(outputDir)) {
	const existingFiles = readdirSync(outputDir).filter(
		(f) => f.includes('init_auth') || f.includes('update_auth')
	);
	for (const existingFile of existingFiles) {
		const existingPath = join(outputDir, existingFile);
		unlinkSync(existingPath);
		console.log(`Deleted existing file: ${existingFile}`);
	}
}

for (let i = 0; i < files.length; i++) {
	const file = files[i];
	const inputPath = join(drizzleDir, file);

	let sql = readFileSync(inputPath, 'utf-8');

	// Remove Drizzle's statement-breakpoint markers
	sql = sql.replace(/--> statement-breakpoint/g, '');

	// Generate migration name based on file index
	let migrationName;
	if (i === 0) {
		migrationName = 'init_auth';
	} else {
		migrationName = `update_auth_${i}`;
	}

	// Generate timestamp-based filename with descriptive name
	const timestamp = new Date()
		.toISOString()
		.replace(/[-:]/g, '')
		.replace(/\.\d+Z$/, '')
		.replace('T', '');
	const outputFileName = `${timestamp}_${migrationName}.sql`;
	const outputPath = join(outputDir, outputFileName);

	// Split into individual statements
	const statements = sql.split(';').filter((s) => s.trim());

	// Build the up migration wrapped in StatementBegin/StatementEnd
	const upSql = statements.join(';\n');

	// Build the down migration (reverse of up)
	const downSql = generateDownMigration(statements);

	// Wrap with Goose annotations
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
}

function generateDownMigration(statements) {
	// Simple heuristic: reverse CREATE statements to DROP
	// This is basic - you may need to enhance based on your needs
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

function generateMigrationName(sql) {
	// Since this is for auth tables and will only be one file, use simple naming
	const createTableMatches = [...sql.matchAll(/CREATE TABLE\s+"?(\w+)"?/gi)];

	if (createTableMatches.length > 0) {
		const tables = createTableMatches.map((m) => m[1]).sort();
		const isInitialAuth =
			tables.includes('user') &&
			tables.includes('session') &&
			tables.includes('account') &&
			tables.includes('verification');

		if (isInitialAuth) {
			return 'init_auth';
		}

		return 'create_' + tables.join('_');
	}

	return 'auth_migration';
}
