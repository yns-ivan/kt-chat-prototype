// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  // Update compatibility date
  compatibilityDate: '2025-07-25',
  
  // Development tools
  devtools: { enabled: true }, // Enable for development
  
  // Modules
  modules: [
    '@nuxt/ui',
    '@nuxt/image', 
    '@nuxt/eslint'
  ],

  // Runtime config for API endpoints
  runtimeConfig: {
    // Private keys (server-side only)
    // apiSecret: process.env.API_SECRET,
    
    // Public keys (client-side)
    public: {
      apiBaseUrl: process.env.API_BASE_URL || 'http://localhost:8080',
      wsUrl: process.env.WS_URL || 'ws://localhost:8080',
      nodeEnv: process.env.NODE_ENV || 'development'
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

  // Build configuration for SPA
  nitro: {
    static: true
  },

  // Development server
  devServer: {
    port: 3000
  }
})
