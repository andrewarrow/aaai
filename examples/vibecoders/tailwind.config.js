/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./templates/**/*.html",
    "./static/src/**/*.js",
  ],
  theme: {
    extend: {
      colors: {
        'primary': '#3b82f6',
        'secondary': '#6366f1',
        'dark': '#1e293b',
        'darker': '#0f172a',
      },
    },
  },
  plugins: [],
}