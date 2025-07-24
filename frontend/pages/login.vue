<template>
  <div style="min-height: 100vh; display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 48px 16px;">
    <div style="max-width: 448px; width: 100%; display: flex; flex-direction: column; gap: 32px;">
      <div style="text-align: center;">
        <div style="width: 80px; height: 80px; background: linear-gradient(135deg, #3b82f6, #8b5cf6); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 24px;">
          <span style="color: white; font-size: 32px;">💬</span>
        </div>
        <h2 style="font-size: 32px; font-weight: bold; background: linear-gradient(135deg, #3b82f6, #8b5cf6); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text;">
          Welcome to KT Chat
        </h2>
        <p style="margin-top: 8px; color: #6b7280;">
          Sign in to your account or
          <button
            @click="showRegister = true"
            style="font-weight: 600; color: #3b82f6; background: none; border: none; cursor: pointer; transition: color 0.2s;"
          >
            create a new account
          </button>
        </p>
      </div>
      
      <div style="background: rgba(255, 255, 255, 0.9); backdrop-filter: blur(10px); border-radius: 16px; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1); border: 1px solid rgba(0, 0, 0, 0.1); padding: 32px;">
        <form style="display: flex; flex-direction: column; gap: 24px;" @submit.prevent="handleLogin">
          <div style="display: flex; flex-direction: column; gap: 16px;">
            <div>
              <label style="display: block; font-size: 14px; font-weight: 500; color: #374151; margin-bottom: 8px;">
                Username *
              </label>
              <input
                v-model="form.username"
                type="text"
                placeholder="Enter your username"
                required
                style="width: 100%; padding: 8px 12px; border: 1px solid #d1d5db; border-radius: 8px; background: white; color: #111827; font-size: 14px;"
              />
            </div>
            
            <div>
              <label style="display: block; font-size: 14px; font-weight: 500; color: #374151; margin-bottom: 8px;">
                Password *
              </label>
              <input
                v-model="form.password"
                type="password"
                placeholder="Enter your password"
                required
                style="width: 100%; padding: 8px 12px; border: 1px solid #d1d5db; border-radius: 8px; background: white; color: #111827; font-size: 14px;"
              />
            </div>
          </div>

          <div v-if="error" style="background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 12px;">
            <p style="color: #dc2626; font-size: 14px; text-align: center;">
              {{ error }}
            </p>
          </div>

          <div v-if="success" style="background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 8px; padding: 12px;">
            <p style="color: #16a34a; font-size: 14px; text-align: center;">
              {{ success }}
            </p>
          </div>

          <button
            type="submit"
            :disabled="loading"
            style="width: 100%; background: linear-gradient(135deg, #3b82f6, #8b5cf6); color: white; padding: 12px 16px; border-radius: 8px; border: none; font-size: 16px; font-weight: 600; cursor: pointer; box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1); transition: all 0.2s;"
          >
            🔐 Sign in
          </button>
        </form>
      </div>
    </div>

    <!-- Register Modal -->
    <div v-if="showRegister" style="position: fixed; inset: 0; background: rgba(0, 0, 0, 0.5); display: flex; align-items: center; justify-content: center; z-index: 50;">
      <div style="background: white; border-radius: 16px; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1); max-width: 448px; width: 100%; margin: 0 16px;">
        <div style="padding: 24px; border-bottom: 1px solid #e5e7eb;">
          <h3 style="font-size: 18px; font-weight: 600; color: #111827;">Create Account</h3>
        </div>

        <form @submit.prevent="handleRegister" style="padding: 24px; display: flex; flex-direction: column; gap: 16px;">
          <div>
            <label style="display: block; font-size: 14px; font-weight: 500; color: #374151; margin-bottom: 8px;">
              Username *
            </label>
            <input
              v-model="registerForm.username"
              type="text"
              placeholder="Choose a username"
              required
              style="width: 100%; padding: 8px 12px; border: 1px solid #d1d5db; border-radius: 8px; background: white; color: #111827; font-size: 14px;"
            />
          </div>
          
          <div>
            <label style="display: block; font-size: 14px; font-weight: 500; color: #374151; margin-bottom: 8px;">
              Email *
            </label>
            <input
              v-model="registerForm.email"
              type="email"
              placeholder="Enter your email"
              required
              style="width: 100%; padding: 8px 12px; border: 1px solid #d1d5db; border-radius: 8px; background: white; color: #111827; font-size: 14px;"
            />
          </div>
          
          <div>
            <label style="display: block; font-size: 14px; font-weight: 500; color: #374151; margin-bottom: 8px;">
              Password *
            </label>
            <input
              v-model="registerForm.password"
              type="password"
              placeholder="Choose a password"
              required
              style="width: 100%; padding: 8px 12px; border: 1px solid #d1d5db; border-radius: 8px; background: white; color: #111827; font-size: 14px;"
            />
          </div>
        </form>

        <div style="padding: 24px; border-top: 1px solid #e5e7eb; display: flex; justify-content: flex-end; gap: 8px;">
          <button
            @click="showRegister = false"
            style="padding: 8px 16px; color: #374151; background: transparent; border: none; border-radius: 8px; font-size: 14px; cursor: pointer; transition: background-color 0.2s;"
          >
            Cancel
          </button>
          <button
            @click="handleRegister"
            :disabled="registerLoading"
            style="padding: 8px 16px; background: #3b82f6; color: white; border: none; border-radius: 8px; font-size: 14px; cursor: pointer; transition: background-color 0.2s;"
          >
            {{ registerLoading ? 'Creating...' : 'Create Account' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Confirmation Modal -->
    <div v-if="showConfirmation" style="position: fixed; inset: 0; background: rgba(0, 0, 0, 0.5); display: flex; align-items: center; justify-content: center; z-index: 50;">
      <div style="background: white; border-radius: 16px; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1); max-width: 448px; width: 100%; margin: 0 16px;">
        <div style="padding: 24px; border-bottom: 1px solid #e5e7eb;">
          <h3 style="font-size: 18px; font-weight: 600; color: #111827;">Confirm Your Account</h3>
          <p style="margin-top: 8px; color: #6b7280; font-size: 14px;">
            Please check your email for a confirmation code and enter it below.
          </p>
        </div>

        <form @submit.prevent="handleConfirm" style="padding: 24px; display: flex; flex-direction: column; gap: 16px;">
          <div>
            <label style="display: block; font-size: 14px; font-weight: 500; color: #374151; margin-bottom: 8px;">
              Confirmation Code *
            </label>
            <input
              v-model="confirmationForm.code"
              type="text"
              placeholder="Enter 6-digit code"
              required
              maxlength="6"
              style="width: 100%; padding: 8px 12px; border: 1px solid #d1d5db; border-radius: 8px; background: white; color: #111827; font-size: 14px; text-align: center; letter-spacing: 2px;"
            />
          </div>

          <div v-if="confirmationError" style="background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 12px;">
            <p style="color: #dc2626; font-size: 14px; text-align: center;">
              {{ confirmationError }}
            </p>
          </div>

          <div v-if="confirmationSuccess" style="background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 8px; padding: 12px;">
            <p style="color: #16a34a; font-size: 14px; text-align: center;">
              {{ confirmationSuccess }}
            </p>
          </div>
        </form>

        <div style="padding: 24px; border-top: 1px solid #e5e7eb; display: flex; justify-content: space-between; align-items: center;">
          <button
            @click="handleResendCode"
            :disabled="resendLoading"
            style="padding: 8px 16px; color: #3b82f6; background: transparent; border: none; border-radius: 8px; font-size: 14px; cursor: pointer; transition: background-color 0.2s;"
          >
            {{ resendLoading ? 'Sending...' : 'Resend Code' }}
          </button>
          
          <div style="display: flex; gap: 8px;">
            <button
              @click="closeConfirmation"
              style="padding: 8px 16px; color: #374151; background: transparent; border: none; border-radius: 8px; font-size: 14px; cursor: pointer; transition: background-color 0.2s;"
            >
              Cancel
            </button>
            <button
              @click="handleConfirm"
              :disabled="confirmationLoading"
              style="padding: 8px 16px; background: #3b82f6; color: white; border: none; border-radius: 8px; font-size: 14px; cursor: pointer; transition: background-color 0.2s;"
            >
              {{ confirmationLoading ? 'Confirming...' : 'Confirm Account' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineOptions({
  name: 'LoginPage'
})

// Auth composable
const { login, register } = useAuth()

// Form data
const form = ref({
  username: '',
  password: ''
})

const registerForm = ref({
  username: '',
  email: '',
  password: ''
})

const confirmationForm = ref({
  code: ''
})

// UI state
const loading = ref(false)
const registerLoading = ref(false)
const confirmationLoading = ref(false)
const resendLoading = ref(false)
const showRegister = ref(false)
const showConfirmation = ref(false)
const error = ref('')
const success = ref('')
const confirmationError = ref('')
const confirmationSuccess = ref('')

// Store username for confirmation
const pendingConfirmationUsername = ref('')

// Handle login
const handleLogin = async () => {
  if (!form.value.username || !form.value.password) {
    error.value = 'Please fill in all fields'
    return
  }

  loading.value = true
  error.value = ''
  success.value = ''

  try {
    const result = await login(form.value.username, form.value.password)
    
    if (result.success) {
      // Navigate to chat
      navigateTo('/')
    } else {
      // Check if user needs confirmation using error code
      if (result.errorCode === 'USER_NOT_CONFIRMED') {
        pendingConfirmationUsername.value = form.value.username
        showConfirmation.value = true
        error.value = ''
      } else {
        error.value = result.error || 'Login failed'
      }
    }
  } catch {
    error.value = 'Login failed. Please try again.'
  } finally {
    loading.value = false
  }
}

// Handle register
const handleRegister = async () => {
  if (!registerForm.value.username || !registerForm.value.email || !registerForm.value.password) {
    error.value = 'Please fill in all fields'
    return
  }

  registerLoading.value = true
  error.value = ''

  try {
    const result = await register(
      registerForm.value.username,
      registerForm.value.email,
      registerForm.value.password
    )
    
    if (result.success) {
      showRegister.value = false
      // Auto-fill login form
      form.value.username = registerForm.value.username
      form.value.password = registerForm.value.password
      
      // Clear register form
      registerForm.value = { username: '', email: '', password: '' }
      
      // Show success message and prompt to login
      success.value = 'Registration successful! Please check your email for confirmation code and login.'
      error.value = ''
    } else {
      error.value = result.error || 'Registration failed'
    }
  } catch {
    error.value = 'Registration failed. Please try again.'
  } finally {
    registerLoading.value = false
  }
}

// Handle confirmation
const handleConfirm = async () => {
  if (!confirmationForm.value.code) {
    confirmationError.value = 'Please enter the confirmation code'
    return
  }

  confirmationLoading.value = true
  confirmationError.value = ''
  confirmationSuccess.value = ''

  try {
    await $fetch('/api/v1/auth/confirm', {
      baseURL: 'http://localhost:8080',
      method: 'POST',
      body: {
        username: pendingConfirmationUsername.value,
        confirmation_code: confirmationForm.value.code
      }
    })

    confirmationSuccess.value = 'Account confirmed successfully! You can now login.'
    
    // Close confirmation modal after 2 seconds
    setTimeout(() => {
      closeConfirmation()
      // Try to login automatically
      handleLogin()
    }, 2000)

  } catch (error) {
    console.error('Confirmation error:', error)
    let errorMessage = 'Confirmation failed'
    if (error && typeof error === 'object' && 'data' in error) {
      const errorData = error.data
      if (errorData && typeof errorData === 'object' && 'error' in errorData && errorData.error) {
        errorMessage = errorData.error.message
      }
    }
    confirmationError.value = errorMessage
  } finally {
    confirmationLoading.value = false
  }
}

// Handle resend confirmation code
const handleResendCode = async () => {
  resendLoading.value = true
  confirmationError.value = ''

  try {
    await $fetch('/api/v1/auth/resend-confirmation', {
      baseURL: 'http://localhost:8080',
      method: 'POST',
      body: {
        username: pendingConfirmationUsername.value
      }
    })

    confirmationSuccess.value = 'Confirmation code sent! Please check your email.'
    
    // Clear success message after 3 seconds
    setTimeout(() => {
      confirmationSuccess.value = ''
    }, 3000)

  } catch (error) {
    console.error('Resend error:', error)
    let errorMessage = 'Failed to resend code'
    if (error && typeof error === 'object' && 'data' in error) {
      const errorData = error.data
      if (errorData && typeof errorData === 'object' && 'error' in errorData && errorData.error) {
        errorMessage = errorData.error.message
      }
    }
    confirmationError.value = errorMessage
  } finally {
    resendLoading.value = false
  }
}

// Close confirmation modal
const closeConfirmation = () => {
  showConfirmation.value = false
  confirmationForm.value.code = ''
  confirmationError.value = ''
  confirmationSuccess.value = ''
  pendingConfirmationUsername.value = ''
}
</script> 