import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../api'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)

  function setAuth(newToken, newUser) {
    token.value = newToken
    user.value = newUser
    localStorage.setItem('token', newToken)
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  function clearAuth() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  async function login(username, password) {
    const res = await api.auth.login(username, password)
    if (res.code === 0) {
      setAuth(res.data.token, res.data.user)
    }
    return res
  }

  async function register(username, password, email) {
    const res = await api.auth.register(username, password, email)
    if (res.code === 0) {
      setAuth(res.data.token, res.data.user)
    }
    return res
  }

  async function logout() {
    await api.auth.logout()
    clearAuth()
  }

  async function refreshProfile() {
    const res = await api.user.getProfile()
    if (res.code === 0) {
      user.value = res.data.user
      localStorage.setItem('user', JSON.stringify(res.data.user))
    }
    return res
  }

  return {
    token,
    user,
    isLoggedIn,
    login,
    register,
    logout,
    refreshProfile,
    clearAuth
  }
})
