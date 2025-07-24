// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  
  // Modules
  modules: [
    '@nuxt/ui',
    '@nuxt/image', 
    '@nuxt/eslint'
  ],

  // Runtime config for API endpoints
  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.API_BASE_URL || 'http://localhost:8080',
      wsUrl: process.env.WS_URL || 'ws://localhost:8080'
    }
  },

  // App configuration
  app: {
    head: {
      title: 'KT Chat',
      meta: [
        { name: 'description', content: 'Real-time chat application' }
      ]
    }
  },





  // Development server
  devServer: {
    port: 3000
  }
})