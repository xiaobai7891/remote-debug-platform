import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'

const instance = axios.create({
  baseURL: '/api',
  timeout: 30000
})

// 请求拦截器
instance.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => Promise.reject(error)
)

// 响应拦截器
instance.interceptors.response.use(
  response => response.data,
  error => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      router.push('/login')
      ElMessage.error('登录已过期，请重新登录')
    } else {
      ElMessage.error(error.response?.data?.message || '请求失败')
    }
    return Promise.reject(error)
  }
)

export default {
  auth: {
    login: (username, password) => instance.post('/auth/login', { username, password }),
    register: (username, password, email) => instance.post('/auth/register', { username, password, email }),
    logout: () => instance.post('/auth/logout')
  },
  user: {
    getProfile: () => instance.get('/user/profile'),
    updateProfile: (data) => instance.put('/user/profile', data),
    updatePassword: (oldPassword, newPassword) => instance.put('/user/password', { old_password: oldPassword, new_password: newPassword }),
    getOrders: (page = 1, pageSize = 20) => instance.get('/user/orders', { params: { page, page_size: pageSize } }),
    getDashboard: () => instance.get('/user/dashboard')
  },
  ports: {
    list: () => instance.get('/ports'),
    create: (data) => instance.post('/ports', data),
    get: (id) => instance.get(`/ports/${id}`),
    update: (id, data) => instance.put(`/ports/${id}`, data),
    delete: (id) => instance.delete(`/ports/${id}`),
    renew: (id, days) => instance.post(`/ports/${id}/renew`, { days }),
    change: (id) => instance.post(`/ports/${id}/change`),
    regenerateToken: (id) => instance.post(`/ports/${id}/regenerate-token`)
  },
  devices: {
    list: () => instance.get('/devices'),
    get: (id) => instance.get(`/devices/${id}`),
    getLogs: (id, limit = 100) => instance.get(`/devices/${id}/logs`, { params: { limit } }),
    exec: (id, code) => instance.post(`/devices/${id}/exec`, { code }),
    screenshot: (id) => instance.post(`/devices/${id}/screenshot`)
  },
  admin: {
    getStats: () => instance.get('/admin/stats'),
    getUsers: (page = 1, pageSize = 20, keyword = '') => instance.get('/admin/users', { params: { page, page_size: pageSize, keyword } }),
    rechargeUser: (id, amount, description) => instance.post(`/admin/users/${id}/recharge`, { amount, description }),
    updateUserStatus: (id, status) => instance.put(`/admin/users/${id}/status`, { status }),
    updateUser: (id, data) => instance.put(`/admin/users/${id}`, data),
    getPorts: (page = 1, pageSize = 20) => instance.get('/admin/ports', { params: { page, page_size: pageSize } }),
    stopPort: (id) => instance.post(`/admin/ports/${id}/stop`),
    startPort: (id) => instance.post(`/admin/ports/${id}/start`),
    getDevices: (page = 1, pageSize = 20, status = '') => instance.get('/admin/devices', { params: { page, page_size: pageSize, status } }),
    getOrders: (page = 1, pageSize = 20, type = '') => instance.get('/admin/orders', { params: { page, page_size: pageSize, type } }),
    getConfig: () => instance.get('/admin/config'),
    updateConfig: (data) => instance.put('/admin/config', data)
  }
}
