import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import topLevelAwait from 'vite-plugin-top-level-await';

export default defineConfig({
	plugins: [sveltekit(), topLevelAwait()],
	server: {},
	build: {
		rollupOptions: {
			external: ['a'],
		},
	},
});
