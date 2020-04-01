import dotenv from 'dotenv';
import svelte from 'rollup-plugin-svelte';
import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import livereload from 'rollup-plugin-livereload';
import { terser } from 'rollup-plugin-terser';
import postcss from 'rollup-plugin-postcss';
import alias from '@rollup/plugin-alias';
import replace from '@rollup/plugin-replace';
import strip from '@rollup/plugin-strip'

// preload .env variables
dotenv.config()
const firebase = {
	apiKey: process.env.FIREBASE_APIKEY,
	authDomain: process.env.FIREBASE_AUTHDOMAIN,
	projectId: process.env.FIREBASE_PID,
	appId: process.env.FIREBASE_APPID
}
const gateway = {
    login: process.env.GW_LOGIN || 'http://localhost:8080/pgw/login',
    logout: process.env.GW_LOGOUT || 'http://localhost:8080/pgw/logout',
    topics: process.env.GW_TOPICS || 'http://localhost:8080/ggw/topics',
	bets: process.env.GW_BETS || 'http://localhost:8080/ggw/bets',
	profile: process.env.GW_PROFILE || 'http://localhost:8080/ggw/profile'
}
const production = !process.env.ROLLUP_WATCH;

const postcssOptions = () => ({
	extensions: ['.scss', '.sass'],
	extract: false,
	minimize: true,
	use: [
		['sass', {
		includePaths: [
			'./src/theme',
			'./node_modules',
		],
		}],
	],
});

export default {
	input: 'src/main.js',
	output: {
		sourcemap: true,
		format: 'iife',
		name: 'app',
		file: 'public/build/bundle.js'
	},
	plugins: [
		strip(),
		replace({
			'exclude': 'node_modules/**',
			'process.browser': true,
			'process.env.FIREBASE_APIKEY': JSON.stringify(firebase.apiKey),
			'process.env.FIREBASE_AUTHDOMAIN': JSON.stringify(firebase.authDomain),
			'process.env.FIREBASE_PID': JSON.stringify(firebase.projectId),
			'process.env.FIREBASE_APPID': JSON.stringify(firebase.appId),

			'process.env.GW_LOGIN': JSON.stringify(gateway.login),
			'process.env.GW_LOGOUT': JSON.stringify(gateway.logout),
			'process.env.GW_TOPICS': JSON.stringify(gateway.topics),
			'process.env.GW_BETS': JSON.stringify(gateway.bets),
			'process.env.GW_PROFILE': JSON.stringify(gateway.profile)
		}),
		svelte({
			// enable run-time checks when not in production
			dev: !production,
			// we'll extract any component CSS out into
			// a separate file - better for performance
			css: css => {
				css.write('public/build/bundle.css');
			}
		}),
		alias({
			entries: [
				{ find: 'components', replacement: '../../components' },
				{ find: 'firebase-conf', replacement: '../../config/firebase' },
				{ find: 'gateway-conf', replacement: '../../config/gateway' }
			]
		}),

		// If you have external dependencies installed from
		// npm, you'll most likely need these plugins. In
		// some cases you'll need additional configuration -
		// consult the documentation for details:
		// https://github.com/rollup/plugins/tree/master/packages/commonjs
		resolve({
			browser: true,
			dedupe: ['svelte']
		}),
		commonjs(),
		postcss(postcssOptions()),

		// In dev mode, call `npm run start` once
		// the bundle has been generated
		!production && serve(),

		// Watch the `public` directory and refresh the
		// browser on changes when not in production
		!production && livereload('public'),

		// If we're building for production (npm run build
		// instead of npm run dev), minify
		production && terser()
	],
	watch: {
		clearScreen: false
	}
};

function serve() {
	let started = false;

	return {
		writeBundle() {
			if (!started) {
				started = true;

				require('child_process').spawn('npm', ['run', 'start', '--', '--dev'], {
					stdio: ['ignore', 'inherit', 'inherit'],
					shell: true
				});
			}
		}
	};
}
