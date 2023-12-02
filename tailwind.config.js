const defaultTheme = require("tailwindcss/defaultTheme");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./pkg/web_assets/*.html"],
  theme: {
    extend: {
      fontFamily: {
        'sans': ['GeistMono', ...defaultTheme.fontFamily.sans],
      },
    },
  },
  plugins: [],
}

