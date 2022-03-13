// import adapter from '@sveltejs/adapter-auto';
import sveltePreprocess from 'svelte-preprocess';
import svg from '@poppanator/sveltekit-svg';
import adapter from '@sveltejs/adapter-static';
import path from 'path';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		adapter: adapter({
			fallback: 'index.html'
		}),

		prerender: {
			enabled: false
		},

		vite: {
			resolve: {
				alias: {
					// these are the aliases and paths to them
					'@components': path.resolve('./src/components'),
					'@lib': path.resolve('./src/lib')
				}
			},
			plugins: [
				svg({
					includePaths: ['./src/assets/icons/'],
					svgoOptions: {
						multipass: true,
						plugins: [
							'preset-default',
							{ name: 'removeAttrs', params: { attrs: '(fill|stroke)' } }
						]
					}
				})
			]
		}
	},
	preprocess: sveltePreprocess({
		pug: true,
		scss: true
	})
};

export default config;
