import {fileURLToPath, URL} from 'node:url'

import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig(({command, mode, ssrBuild}) => {
	const ret = {
		plugins: [vue()],
		resolve: {
			alias: {
				'@': fileURLToPath(new URL('./src', import.meta.url))
			}
		},
	};
	if (command === 'serve' && mode !== 'developement-external') {
		ret.define = {
			"__API_URL__": JSON.stringify("http://localhost:3000"),
		};
	} else if (mode === 'embedded') {
		ret.define = {
			"__API_URL__": JSON.stringify("/"),
		};
	} else {
		ret.define = {
			"__API_URL__": JSON.stringify("<your API URL>"),
		};
	}
	return ret;
})
