module.exports = {
  purge: [
    "./app/views/**/*.html.erb",
    "./app/webpack/**/*.js",
    "./app/webpack/**/*.jsx",
  ],
  theme: {
    extend: {
      spacing: {
        "84": "21rem",
      },
      fontFamily: {
        "tempbugfix": [
          "Inter",
          "system-ui",
          "-apple-system",
          "BlinkMacSystemFont",
          "Segoe UI",
          "Roboto",
          "Helvetica Neue",
          "Arial",
          "Noto Sans",
          "sans-serif",
          "Apple Color Emoji",
          "Segoe UI Emoji",
          "Segoe UI Symbol",
          "Noto Color Emoji",
        ],
      }
    },
  },
  variants: {borderColor: ['responsive', 'hover', 'focus', 'focus-within']},
  plugins: [],
};
