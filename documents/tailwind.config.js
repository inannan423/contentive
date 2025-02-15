/** @type {import('tailwindcss').Config} */
export default {
    content: [
      './app/**/*.{js,jsx,ts,tsx,md,mdx}',
      './app/layout.jsx',
      './content/**/*.{md,mdx}',
   
      // Or if using `src` directory:
      './src/**/*.{js,jsx,ts,tsx,md,mdx}'
    ],
    theme: {
      extend: {}
    },
    darkMode: 'class',
    plugins: []
}