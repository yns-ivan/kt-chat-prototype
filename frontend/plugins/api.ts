export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  
  return {
    provide: {
      api: config.public.apiBaseUrl
    }
  }
}) 