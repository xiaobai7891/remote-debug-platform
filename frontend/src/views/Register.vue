<template>
  <div class="auth-container">
    <!-- Animated background -->
    <div class="cyber-grid"></div>
    <div class="floating-particles">
      <div v-for="n in 20" :key="n" class="particle" :style="getParticleStyle(n)"></div>
    </div>

    <!-- Register card -->
    <div class="auth-card">
      <div class="card-glow"></div>

      <!-- Logo section -->
      <div class="auth-header">
        <div class="logo-icon">
          <svg viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M20 4L36 12V28L20 36L4 28V12L20 4Z" stroke="currentColor" stroke-width="2" fill="none"/>
            <path d="M20 4V36" stroke="currentColor" stroke-width="1.5" opacity="0.5"/>
            <path d="M4 12L36 28" stroke="currentColor" stroke-width="1.5" opacity="0.5"/>
            <path d="M36 12L4 28" stroke="currentColor" stroke-width="1.5" opacity="0.5"/>
            <circle cx="20" cy="20" r="6" stroke="currentColor" stroke-width="2"/>
            <path d="M17 20H23M20 17V23" stroke="currentColor" stroke-width="2"/>
          </svg>
        </div>
        <h1 class="auth-title">
          <span class="title-text">CREATE</span>
          <span class="title-accent">ACCOUNT</span>
        </h1>
        <p class="auth-subtitle">// 创建新的调试终端</p>
      </div>

      <!-- Form section -->
      <el-form ref="formRef" :model="form" :rules="rules" @submit.prevent="handleRegister" class="auth-form">
        <div class="input-group">
          <label class="input-label">
            <span class="label-icon">&#62;</span>
            USER_ID
          </label>
          <el-form-item prop="username">
            <el-input
              v-model="form.username"
              placeholder="设置用户名"
              size="large"
              :prefix-icon="User"
            />
          </el-form-item>
        </div>

        <div class="input-group">
          <label class="input-label">
            <span class="label-icon">&#62;</span>
            EMAIL_ADDRESS
          </label>
          <el-form-item prop="email">
            <el-input
              v-model="form.email"
              placeholder="输入邮箱（可选）"
              size="large"
              :prefix-icon="Message"
            />
          </el-form-item>
        </div>

        <div class="input-group">
          <label class="input-label">
            <span class="label-icon">&#62;</span>
            ACCESS_KEY
          </label>
          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="设置密码"
              size="large"
              show-password
              :prefix-icon="Lock"
            />
          </el-form-item>
        </div>

        <div class="input-group">
          <label class="input-label">
            <span class="label-icon">&#62;</span>
            CONFIRM_KEY
          </label>
          <el-form-item prop="confirmPassword">
            <el-input
              v-model="form.confirmPassword"
              type="password"
              placeholder="确认密码"
              size="large"
              show-password
              :prefix-icon="Lock"
            />
          </el-form-item>
        </div>

        <el-button
          type="primary"
          native-type="submit"
          :loading="loading"
          size="large"
          class="submit-btn"
        >
          <span class="btn-text">{{ loading ? '正在创建...' : '创建账号' }}</span>
          <span class="btn-icon" v-if="!loading">
            <el-icon><Right /></el-icon>
          </span>
        </el-button>

        <div class="auth-links">
          <span class="link-text">已有账号？</span>
          <router-link to="/login" class="link-action">
            返回登录 <el-icon><ArrowRight /></el-icon>
          </router-link>
        </div>
      </el-form>

      <!-- Decorative elements -->
      <div class="corner-decor top-left"></div>
      <div class="corner-decor top-right"></div>
      <div class="corner-decor bottom-left"></div>
      <div class="corner-decor bottom-right"></div>
    </div>

    <!-- Status bar -->
    <div class="status-bar">
      <div class="status-item">
        <span class="status-dot online"></span>
        <span>SYSTEM_ONLINE</span>
      </div>
      <div class="status-item">
        <span>v1.0.0</span>
      </div>
      <div class="status-item">
        <span>{{ currentTime }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Message, Right, ArrowRight } from '@element-plus/icons-vue'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref()
const loading = ref(false)
const currentTime = ref('')

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const validatePassword = (rule, value, callback) => {
  if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度为 3-20 个字符', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 32, message: '密码长度为 6-32 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validatePassword, trigger: 'blur' }
  ]
}

function getParticleStyle(n) {
  return {
    '--delay': `${Math.random() * 5}s`,
    '--duration': `${15 + Math.random() * 10}s`,
    '--x': `${Math.random() * 100}%`,
    '--size': `${2 + Math.random() * 4}px`
  }
}

function updateTime() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour12: false })
}

let timeInterval

onMounted(() => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
})

onUnmounted(() => {
  clearInterval(timeInterval)
})

async function handleRegister() {
  await formRef.value.validate()
  loading.value = true
  try {
    const res = await userStore.register(form.username, form.password, form.email)
    if (res.code === 0) {
      ElMessage.success('账号创建成功')
      router.push('/')
    } else {
      ElMessage.error(res.message)
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-void);
  position: relative;
  overflow: hidden;
}

/* Cyberpunk grid background */
.cyber-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(168, 85, 247, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(168, 85, 247, 0.03) 1px, transparent 1px);
  background-size: 50px 50px;
  animation: grid-move 20s linear infinite;
}

@keyframes grid-move {
  0% { transform: perspective(500px) rotateX(60deg) translateY(0); }
  100% { transform: perspective(500px) rotateX(60deg) translateY(50px); }
}

/* Floating particles */
.floating-particles {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.particle {
  position: absolute;
  width: var(--size);
  height: var(--size);
  background: var(--neon-purple);
  border-radius: 50%;
  left: var(--x);
  bottom: -10px;
  opacity: 0.6;
  animation: float-up var(--duration) var(--delay) infinite linear;
  box-shadow: 0 0 10px var(--neon-purple);
}

@keyframes float-up {
  0% { transform: translateY(0) scale(1); opacity: 0; }
  10% { opacity: 0.6; }
  90% { opacity: 0.6; }
  100% { transform: translateY(-100vh) scale(0.5); opacity: 0; }
}

/* Auth card */
.auth-card {
  position: relative;
  width: 420px;
  padding: 40px 48px;
  background: linear-gradient(135deg, var(--bg-card) 0%, var(--bg-surface) 100%);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-xl);
  z-index: 10;
  animation: card-appear 0.6s ease-out;
}

@keyframes card-appear {
  from {
    opacity: 0;
    transform: translateY(30px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.card-glow {
  position: absolute;
  inset: -1px;
  border-radius: var(--radius-xl);
  background: linear-gradient(135deg, var(--neon-purple), var(--neon-cyan));
  opacity: 0.15;
  filter: blur(20px);
  z-index: -1;
}

/* Corner decorations */
.corner-decor {
  position: absolute;
  width: 20px;
  height: 20px;
  border-color: var(--neon-purple);
  opacity: 0.5;
}

.corner-decor.top-left {
  top: 12px;
  left: 12px;
  border-top: 2px solid;
  border-left: 2px solid;
}

.corner-decor.top-right {
  top: 12px;
  right: 12px;
  border-top: 2px solid;
  border-right: 2px solid;
}

.corner-decor.bottom-left {
  bottom: 12px;
  left: 12px;
  border-bottom: 2px solid;
  border-left: 2px solid;
}

.corner-decor.bottom-right {
  bottom: 12px;
  right: 12px;
  border-bottom: 2px solid;
  border-right: 2px solid;
}

/* Logo section */
.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo-icon {
  width: 56px;
  height: 56px;
  margin: 0 auto 16px;
  color: var(--neon-purple);
  animation: pulse-glow 2s ease-in-out infinite;
}

@keyframes pulse-glow {
  0%, 100% { filter: drop-shadow(0 0 10px var(--neon-purple)); }
  50% { filter: drop-shadow(0 0 20px var(--neon-purple)); }
}

.auth-title {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 800;
  letter-spacing: 2px;
  margin-bottom: 8px;
}

.title-text {
  color: var(--text-primary);
}

.title-accent {
  color: var(--neon-purple);
  margin-left: 8px;
}

.auth-subtitle {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

/* Form styles */
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.input-label {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 500;
  color: var(--text-muted);
  letter-spacing: 1px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.label-icon {
  color: var(--neon-purple);
}

.auth-form :deep(.el-form-item) {
  margin-bottom: 0;
}

.auth-form :deep(.el-input__wrapper) {
  padding: 10px 14px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
}

.auth-form :deep(.el-input__wrapper:hover) {
  border-color: var(--neon-purple-dim);
}

.auth-form :deep(.el-input__wrapper.is-focus) {
  border-color: var(--neon-purple);
  box-shadow: 0 0 0 3px rgba(168, 85, 247, 0.1), inset 0 0 20px rgba(168, 85, 247, 0.05);
}

.auth-form :deep(.el-input__inner) {
  font-family: var(--font-mono);
  font-size: 14px;
}

/* Submit button */
.submit-btn {
  height: 48px;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 1px;
  margin-top: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: linear-gradient(135deg, var(--neon-purple), var(--neon-purple-dim));
  border: none;
  transition: all 0.3s ease;
}

.submit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(168, 85, 247, 0.3);
}

.submit-btn:active {
  transform: translateY(0);
}

.btn-icon {
  transition: transform 0.3s ease;
}

.submit-btn:hover .btn-icon {
  transform: translateX(4px);
}

/* Links */
.auth-links {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 4px;
}

.link-text {
  font-size: 13px;
  color: var(--text-muted);
}

.link-action {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 500;
  color: var(--neon-purple);
  text-decoration: none;
  transition: all 0.2s ease;
}

.link-action:hover {
  color: var(--text-primary);
  text-shadow: 0 0 10px var(--neon-purple);
}

.link-action :deep(.el-icon) {
  transition: transform 0.2s ease;
}

.link-action:hover :deep(.el-icon) {
  transform: translateX(4px);
}

/* Status bar */
.status-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 32px;
  background: var(--bg-surface);
  border-top: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
  z-index: 100;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--neon-green);
  box-shadow: 0 0 8px var(--neon-green);
}

/* Responsive */
@media (max-width: 480px) {
  .auth-card {
    width: calc(100% - 32px);
    padding: 28px 24px;
    margin: 16px;
  }

  .auth-title {
    font-size: 22px;
  }
}
</style>
