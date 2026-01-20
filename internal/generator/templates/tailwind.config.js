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
    plugins: [<!-- TAILWIND_PLUGIN -->]<!-- DAISYUI_CONFIG -->
}
