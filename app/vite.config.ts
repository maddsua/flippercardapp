import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
	plugins: [vue()],
	server: {
		host: '0.0.0.0',
		proxy: {
			'/api': {
				target: 'http://localhost:8280/api',
				rewrite: (path) => path.replace(/^\/api\//, '/'),
			},
			'/media': {
				target: 'http://localhost:8280/media',
				rewrite: (path) => path.replace(/^\/media\//, '/'),
			}
		}
	}
});
