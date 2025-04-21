/** @type {import('tailwindcss').Config} */
export default {
	content: ["./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}"],
	theme: {
		extend: {
			colors: {
				primary: {
					50:  "#eaf1ff",   // very light blue
					100: "#d6e3ff",   // light pastel
					200: "#adc8ff",   // soft sky blue
					300: "#7aaaff",   // mid-tone blue
					400: "#478cff",   // medium
					500: "#1f6fff",   // strong blue
					600: "#004fd6",   // deep blue
					700: "#003cb3",   // dark blue
					800: "#002c8a",   // navy tone
					900: "#001a66",   // very dark blue
					950: "#000d33",   // almost black with blue
				},
				secondary: {
					50:  "#fff2ed",   // very light
					100: "#ffdcd1",   // soft coral
					200: "#ffb4a3",   // warm mid-light
					300: "#ff8166",   // energetic coral
					400: "#ff5733",   // bold + readable
					500: "#e6461f",   // slightly deeper
					600: "#c73714",   // punchy but solid
					700: "#a82b10",   // deep orange-red
					800: "#88220d",   // dark red-orange
					900: "#5c1608",   // bold dark base
					950: "#310800",   // nearly black with warm tint
				},
			},
			typography: {
				DEFAULT: {
					css: {
						color: "rgb(229, 231, 235)",
						"h1, h2, h3, h4": {
							color: "#a8cbff",
						},
						a: {
							color: "rgb(96, 165, 250)",
							"&:hover": {
								color: "rgb(147, 197, 253)",
							},
						},
						code: {
							color: "rgb(243, 244, 246)",
							backgroundColor: "rgb(31, 41, 55)",
						},
						blockquote: {
							color: "rgb(229, 231, 235)",
							borderLeftColor: "rgb(75, 85, 99)",
						},
						strong: {
							color: "#ffaa71",
						},
						th: {
							color: "#ffaa71",
						},
					},
				},
			},
		},
	},
	plugins: [require("@tailwindcss/typography")],
};
