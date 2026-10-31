/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{svelte,js}'],
  theme: {
    extend: {
      colors: {
        ink: '#0D1B2A',
        navy: '#1B263B',
        steel: '#415A77',
        sage: '#778D7A',
        sand: '#D4C4A8',
        cream: '#F4F1DE'
      },
      fontFamily: {
        jp: ['"Noto Sans JP"', 'Meiryo', 'Yu Gothic', 'sans-serif']
      }
    }
  },
  plugins: []
}
