// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: ['@nuxt/eslint'],
  app: {
    head: {
      htmlAttrs: { lang: 'en', class: 'light' },
      title: 'OmniChat Inbox Dashboard',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1.0' }
      ],
      link: [
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700;800&display=swap'
        },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap'
        }
      ],
      script: [
        { src: 'https://cdn.tailwindcss.com?plugins=forms,container-queries' },
        {
          id: 'tailwind-config',
          innerHTML: `tailwind.config = {darkMode: "class", theme: {extend: {colors: {primary: "#a900a7", "on-secondary-fixed": "#380037", "inverse-on-surface": "#ffebf7", "on-surface-variant": "#564051", surface: "#fff7f9", "surface-variant": "#f4dcec", "on-background": "#251722", secondary: "#a2209f", "surface-container-high": "#fae1f2", "surface-tint": "#a900a7", "surface-container-low": "#ffeff8", "surface-container": "#ffe7f7", "on-secondary-container": "#750073", "surface-container-lowest": "#ffffff", "on-secondary": "#ffffff", "error-container": "#ffdad6", "on-tertiary-fixed-variant": "#414b00", "secondary-fixed": "#ffd7f5", "on-tertiary-container": "#262d00", "inverse-primary": "#ffabf2", background: "#fff7f9", "on-primary-fixed": "#380037", "surface-dim": "#ebd3e3", "outline-variant": "#dcbed3", "primary-container": "#ff00fc", "on-tertiary-fixed": "#191e00", "surface-bright": "#fff7f9", "tertiary-fixed-dim": "#bcd143", "primary-fixed": "#ffd7f5", "on-secondary-fixed-variant": "#810080", outline: "#896f83", "on-surface": "#251722", "primary-fixed-dim": "#ffabf2", "on-primary-container": "#510050", "secondary-fixed-dim": "#ffabf2", "inverse-surface": "#3b2c38", error: "#ba1a1a", "on-error-container": "#93000a", "on-error": "#ffffff", "surface-container-highest": "#f4dcec", "on-tertiary": "#ffffff", "on-primary-fixed-variant": "#810080", "secondary-container": "#fe77f3", tertiary: "#586400", "tertiary-container": "#879a00", "on-primary": "#ffffff", "tertiary-fixed": "#d8ee5d"}, borderRadius: {DEFAULT: "0.25rem", lg: "0.5rem", xl: "0.75rem", full: "9999px"}, spacing: {"sidebar-width": "280px", xl: "40px", md: "16px", lg: "24px", "chat-list-width": "360px", sm: "8px", unit: "4px", xs: "4px"}, fontFamily: {"headline-md": ["Inter"], "display-lg-mobile": ["Inter"], "body-sm": ["Inter"], "label-caps": ["Inter"], "body-lg": ["Inter"], "status-label": ["Inter"], "display-lg": ["Inter"], headline: ["Inter"], display: ["Inter"], body: ["Inter"], label: ["Inter"]}, fontSize: {"headline-md": ["20px", {lineHeight: "28px", fontWeight: "600"}], "display-lg-mobile": ["24px", {lineHeight: "32px", letterSpacing: "-0.01em", fontWeight: "700"}], "body-sm": ["14px", {lineHeight: "20px", fontWeight: "400"}], "label-caps": ["12px", {lineHeight: "16px", letterSpacing: "0.05em", fontWeight: "600"}], "body-lg": ["16px", {lineHeight: "24px", fontWeight: "400"}], "status-label": ["11px", {lineHeight: "12px", fontWeight: "700"}], "display-lg": ["32px", {lineHeight: "40px", letterSpacing: "-0.02em", fontWeight: "700"}]}}}};`
        }
      ],
      style: [
        {
          innerHTML: `body { font-family: 'Inter', sans-serif; background-color: #faf8ff; }
        .material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; display: inline-block; line-height: 1; vertical-align: middle; }
        .message-item-active { border-left: 3px solid #0050cb; background-color: #d0e1fb; }
        .scrollbar-hide::-webkit-scrollbar { display: none; }
        .scrollbar-hide { -ms-overflow-style: none; scrollbar-width: none; }`
        }
      ]
    }
  }
})