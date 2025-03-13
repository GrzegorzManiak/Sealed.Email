import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import topLevelAwait from 'vite-plugin-top-level-await';
import { nodePolyfills } from 'vite-plugin-node-polyfills'
export default defineConfig({
	plugins: [
		sveltekit(),
		topLevelAwait(),
		nodePolyfills({
			globals: {
				Buffer: true,
				global: true,
				process: true,
			},
			protocolImports: true,
		})
	],
	server: {},
	build: {
		rollupOptions: {
			external: ['a'],
		},
	}
});
