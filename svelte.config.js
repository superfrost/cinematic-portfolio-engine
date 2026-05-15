import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
    preprocess: vitePreprocess(),
    kit: {
        adapter: adapter({
            pages: 'build',
            assets: 'build',
            fallback: '404.html', // Нужен для корректного роутинга на гите
            precompress: false,
            strict: true
        }),
        paths: {
            // Если твой сайт будет лежать по адресу https://username.github.io/portfolio-repo/,
            // то здесь должно быть указано имя репозитория:
            base: process.env.NODE_ENV === 'production' ? '/portfolio-repo' : ''
        }
    }
};

export default config;
