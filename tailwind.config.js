/** @type import('tailwindcss').Config;*/
export default {
	content: ["./cmd/web/**/*.html", "./cmd/web/**/*.templ"],
	theme: {
		extend: {},
	},
	plugins: [require('daisyui')],
	daisyui: {
		themes: ["luxury"]
	}
}
