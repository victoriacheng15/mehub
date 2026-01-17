/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./dist/**/*.html", "./internal/templates/**/*.html"],
  theme: {
    extend: {
      colors: {
        primary: {
          50:  "#eaf1ff",
          100: "#d6e3ff",
          200: "#adc8ff",
          300: "#7aaaff",
          400: "#478cff",
          500: "#1f6fff",
          600: "#004fd6",
          700: "#003cb3",
          800: "#002c8a",
          900: "#001a66",
          950: "#000d33",
        },
        secondary: {
          50:  "#fff2ed",
          100: "#ffdcd1",
          200: "#ffb4a3",
          300: "#ff8166",
          400: "#ff5733",
          500: "#e6461f",
          600: "#c73714",
          700: "#a82b10",
          800: "#88220d",
          900: "#5c1608",
          950: "#310800",
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
  plugins: [
    require('@tailwindcss/typography'),
  ],
}
