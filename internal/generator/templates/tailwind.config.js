/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./views/**/*.templ",
        "./views/**/*.go",
    ],
    theme: {
        extend: {
            fontFamily: {
                sans: ['Inter', 'system-ui', 'sans-serif'],
            },
        },
    },
    plugins: [],
    // DaisyUI is loaded via CDN, so no plugin needed here
}
