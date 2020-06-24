module.exports = {
  purge: [
    "./app/views/**/*.html.erb",
    "./app/webpack/**/*.js",
    "./app/webpack/**/*.jsx",
  ],
  theme: {
    extend: {
      padding: {
        1.5: "0.375rem",
      },
      maxHeight: {
        "0": "0",
      },
      minHeight: {
        "16": "4rem",
        "1": "1px",
      },
      spacing: {
        "68": "17rem",
        "72": "18rem",
        "84": "21rem",
        "144": "36rem",
      },
      colors: {
        lightgray: "#E6E7ED",
      },
      fontFamily: {
        tempbugfix: [
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
      },
    },
  },
  variants: {
    borderColor: ["responsive", "first", "hover", "focus", "focus-within"],
    cursor: ["responsive", "hover", "focus"],
  },
  plugins: [],
};
