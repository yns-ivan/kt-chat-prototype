<template>
  <div v-if="isInitializing" class="auth-loading-overlay">
    <div class="auth-loading-content">
      <div class="auth-loading-spinner"></div>
      <p class="auth-loading-text">Loading...</p>
    </div>
  </div>
</template>

<script setup>
const { isInitializing } = useAuth()
</script>

<style scoped>
.auth-loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.auth-loading-content {
  text-align: center;
  color: white;
}

.auth-loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(255, 255, 255, 0.3);
  border-top: 4px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

.auth-loading-text {
  font-size: 16px;
  font-weight: 500;
  margin: 0;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style> 