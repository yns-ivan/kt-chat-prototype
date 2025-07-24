<template>
  <div class="min-h-screen" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<script setup>
// Global app setup
useHead({
  title: 'KT Chat',
  meta: [
    { name: 'description', content: 'Real-time chat application' }
  ]
})

// Initialize auth state
const { initAuth } = useAuth()

onMounted(() => {
  initAuth()
})
</script>

 