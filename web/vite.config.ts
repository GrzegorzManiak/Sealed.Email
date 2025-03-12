import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import topLevelAwait from 'vite-plugin-top-level-await';
import inject from '@rollup/plugin-inject'

export default defineConfig({
	plugins: [sveltekit(), topLevelAwait()],
	server: {},
	build: {
		rollupOptions: {
			external: ['a'],
			plugins: [inject({ Buffer: ['buffer/', 'Buffer'] })],
		},
	},
});
