import adapter from '@sveltejs/adapter-node';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: { adapter: adapter(), experimental: { remoteFunctions: true } },
	vitePlugin: {
		dynamicCompileOptions: ({ filename }) =>
			filename.includes('node_modules') ? undefined : { runes: true }
	},
	compilerOptions: {
		experimental: {
			async: true
		}
	}
};

export default config;
