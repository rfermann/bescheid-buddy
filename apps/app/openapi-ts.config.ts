import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
	input: 'https://dev.api.bb.tricirc.de/openapi.yaml',
	output: {
		path: 'src/generated'
	},
	plugins: [
		{
			enums: 'typescript',
			name: '@hey-api/typescript'
		},
		{ name: '@hey-api/sdk', transformer: true },
		{
			dates: true,
			name: '@hey-api/transformers'
		},
		'valibot'
	]
});
