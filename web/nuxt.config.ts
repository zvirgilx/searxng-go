// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false,
  modules: ["@pinia/nuxt"],
  vite: {
    css: {
      preprocessorOptions: {
        less: {
          additionalData: `@import "@/assets/less/global.less";`,
        },
      },
    },
  },
  devtools: { enabled: true },
  css: ["@/assets/less/global.less"],
  // routeRules: {
  //   "/api_proxy/**": {
  //     cors: true,
  //     proxy: { to: "http://localhost:9999/**" },
  //   },
  // },
});
