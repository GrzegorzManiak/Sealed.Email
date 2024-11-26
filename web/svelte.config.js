import adapter from '@sveltejs/adapter-auto';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		alias: {
			"@/*": "./src/lib/components/*",
			"&/*": "./src/static/*",
			"$shadcn": "./src/lib/components/ui",
			"$shadcn/*": "./src/lib/components/ui/*",
			"$local": "./src/lib/components/local",
			"$local/*": "./src/lib/components/local/*"
		},
		paths: {
			base: ''
		},
		adapter: adapter({
			fallback: null,
			strict: true
		}),
		prerender: {
			handleMissingId: 'warn'
		}
	}
};

export default config;
