<template>
  <div class="admin-layout">
    <!-- Sidebar -->
    <aside class="admin-sidebar">
      <!-- Logo -->
      <div class="sidebar-logo">
        <div class="logo-icon">
          <svg viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M20 4L36 12V28L20 36L4 28V12L20 4Z" stroke="currentColor" stroke-width="2" fill="none"/>
            <path d="M20 4V36" stroke="currentColor" stroke-width="1.5" opacity="0.5"/>
            <path d="M4 12L36 28" stroke="currentColor" stroke-width="1.5" opacity="0.5"/>
            <path d="M36 12L4 28" stroke="currentColor" stroke-width="1.5" opacity="0.5"/>
            <circle cx="20" cy="20" r="5" stroke="currentColor" stroke-width="2"/>
          </svg>
        </div>
        <div class="logo-text">
          <span class="logo-title">ADMIN</span>
          <span class="logo-subtitle">CONTROL</span>
        </div>
      </div>

      <!-- Admin Badge -->
      <div class="admin-badge">
        <span class="badge-icon">
          <el-icon><Warning /></el-icon>
        </span>
        <span class="badge-text">ROOT_ACCESS</span>
      </div>

      <!-- Navigation -->
      <nav class="sidebar-nav">
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{ active: isActive(item.path) }"
        >
          <span class="nav-indicator"></span>
          <span class="nav-icon">
            <el-icon><component :is="item.icon" /></el-icon>
          </span>
          <span class="nav-label">{{ item.label }}</span>
          <span class="nav-code">{{ item.code }}</span>
        </router-link>
      </nav>

      <!-- System Status -->
      <div class="system-status">
        <div class="status-header">
          <span class="status-dot"></span>
          <span>SYSTEM_STATUS</span>
        </div>
        <div class="status-items">
          <div class="status-item">
            <span class="item-label">CPU:</span>
            <span class="item-value">{{ systemLoad }}%</span>
          </div>
          <div class="status-item">
            <span class="item-label">MEM:</span>
            <span class="item-value">{{ memoryUsage }}%</span>
          </div>
        </div>
      </div>

      <!-- Back to Frontend -->
      <router-link to="/" class="back-link">
        <el-icon><Back /></el-icon>
        <span>EXIT_ADMIN</span>
      </router-link>
    </aside>

    <!-- Main Content -->
    <div class="admin-main">
      <!-- Header -->
      <header class="admin-header">
        <div class="header-left">
          <div class="breadcrumb">
            <span class="breadcrumb-icon">></span>
            <span class="breadcrumb-root">ADMIN</span>
            <span class="breadcrumb-separator">/</span>
            <span class="breadcrumb-current">{{ currentPageName }}</span>
          </div>
        </div>

        <div class="header-center">
          <div class="server-time">
            <el-icon><Clock /></el-icon>
            <span>{{ currentTime }}</span>
          </div>
        </div>

        <div class="header-right">
          <!-- Notifications -->
          <div class="header-notifications">
            <el-badge :value="3" :max="9" class="notification-badge">
              <button class="notification-btn">
                <el-icon><Bell /></el-icon>
              </button>
            </el-badge>
          </div>

          <!-- User Dropdown -->
          <el-dropdown @command="handleCommand" trigger="click">
            <div class="user-menu">
              <div class="user-avatar">
                <el-icon><UserFilled /></el-icon>
              </div>
              <div class="user-info">
                <span class="user-name">{{ userStore.user?.username }}</span>
                <span class="user-role">ADMINISTRATOR</span>
              </div>
              <el-icon class="dropdown-arrow"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu class="admin-dropdown">
                <el-dropdown-item command="front">
                  <el-icon><Monitor /></el-icon>
                  返回前台
                </el-dropdown-item>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人设置
                </el-dropdown-item>
                <el-dropdown-item command="logout" divided>
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>

      <!-- Content Area -->
      <main class="admin-content">
        <router-view />
      </main>

      <!-- Footer -->
      <footer class="admin-footer">
        <div class="footer-left">
          <span class="footer-version">v1.0.0</span>
          <span class="footer-separator">|</span>
          <span class="footer-text">REMOTE_DEBUG_ADMIN_PANEL</span>
        </div>
        <div class="footer-right">
          <span class="footer-text">SESSION: {{ sessionId }}</span>
        </div>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'
import {
  UserFilled, ArrowDown, Back, Clock, Bell, Monitor, User, SwitchButton,
  DataAnalysis, Connection, Iphone, Wallet, Tools, Warning
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const currentTime = ref('')
const systemLoad = ref(23)
const memoryUsage = ref(45)
const sessionId = ref('')

const navItems = [
  { path: '/admin', label: '数据概览', code: 'DASH', icon: DataAnalysis },
  { path: '/admin/users', label: '用户管理', code: 'USER', icon: User },
  { path: '/admin/ports', label: '端口管理', code: 'PORT', icon: Connection },
  { path: '/admin/devices', label: '设备管理', code: 'DEV', icon: Iphone },
  { path: '/admin/orders', label: '订单管理', code: 'ORD', icon: Wallet },
  { path: '/admin/config', label: '系统配置', code: 'CONF', icon: Tools }
]

const currentPageName = computed(() => {
  const current = navItems.find(item => item.path === route.path)
  return current ? current.code : 'DASHBOARD'
})

function isActive(path) {
  if (path === '/admin') {
    return route.path === '/admin'
  }
  return route.path.startsWith(path)
}

function handleCommand(command) {
  if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  } else if (command === 'front') {
    router.push('/')
  } else if (command === 'profile') {
    router.push('/profile')
  }
}

function updateTime() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour12: false })
}

function generateSessionId() {
  return Math.random().toString(36).substring(2, 10).toUpperCase()
}

function updateSystemStats() {
  // Simulate system stats
  systemLoad.value = Math.floor(15 + Math.random() * 30)
  memoryUsage.value = Math.floor(40 + Math.random() * 25)
}

let timeInterval
let statsInterval

onMounted(() => {
  updateTime()
  sessionId.value = generateSessionId()
  timeInterval = setInterval(updateTime, 1000)
  statsInterval = setInterval(updateSystemStats, 5000)
})

onUnmounted(() => {
  clearInterval(timeInterval)
  clearInterval(statsInterval)
})
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: var(--bg-void);
}

/* Sidebar */
.admin-sidebar {
  width: 240px;
  background: linear-gradient(180deg, var(--bg-card) 0%, var(--bg-surface) 100%);
  border-right: 1px solid var(--border-default);
  display: flex;
  flex-direction: column;
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 100;
}

/* Logo */
.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px;
  border-bottom: 1px solid var(--border-subtle);
}

.logo-icon {
  width: 40px;
  height: 40px;
  color: #f97316;
  animation: logo-pulse 3s ease-in-out infinite;
}

@keyframes logo-pulse {
  0%, 100% { filter: drop-shadow(0 0 8px rgba(249, 115, 22, 0.5)); }
  50% { filter: drop-shadow(0 0 16px rgba(249, 115, 22, 0.8)); }
}

.logo-text {
  display: flex;
  flex-direction: column;
}

.logo-title {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 800;
  color: var(--text-primary);
  letter-spacing: 2px;
}

.logo-subtitle {
  font-family: var(--font-mono);
  font-size: 10px;
  color: #f97316;
  letter-spacing: 1px;
}

/* Admin Badge */
.admin-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 16px 20px;
  padding: 8px 12px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: var(--radius-md);
  font-family: var(--font-mono);
  font-size: 10px;
  color: #ef4444;
  letter-spacing: 1px;
}

.badge-icon {
  display: flex;
  animation: badge-blink 2s ease-in-out infinite;
}

@keyframes badge-blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* Navigation */
.sidebar-nav {
  flex: 1;
  padding: 8px 12px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  margin-bottom: 4px;
  border-radius: var(--radius-md);
  text-decoration: none;
  color: var(--text-secondary);
  position: relative;
  transition: all 0.2s ease;
}

.nav-item:hover {
  background: var(--bg-elevated);
  color: var(--text-primary);
}

.nav-item.active {
  background: rgba(249, 115, 22, 0.1);
  color: #f97316;
}

.nav-indicator {
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 0;
  background: #f97316;
  border-radius: 0 2px 2px 0;
  transition: height 0.2s ease;
}

.nav-item.active .nav-indicator {
  height: 24px;
}

.nav-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  font-size: 16px;
}

.nav-label {
  flex: 1;
  font-size: 13px;
  font-weight: 500;
}

.nav-code {
  font-family: var(--font-mono);
  font-size: 9px;
  color: var(--text-muted);
  padding: 2px 6px;
  background: var(--bg-surface);
  border-radius: var(--radius-sm);
  letter-spacing: 0.5px;
}

.nav-item.active .nav-code {
  background: rgba(249, 115, 22, 0.2);
  color: #f97316;
}

/* System Status */
.system-status {
  margin: 16px 20px;
  padding: 12px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
}

.status-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--neon-green);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--neon-green);
  box-shadow: 0 0 8px var(--neon-green);
  animation: pulse-dot 2s ease-in-out infinite;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status-items {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.status-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-family: var(--font-mono);
  font-size: 11px;
}

.item-label {
  color: var(--text-muted);
}

.item-value {
  color: var(--neon-cyan);
}

/* Back Link */
.back-link {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 16px;
  margin: 16px 20px;
  background: var(--bg-elevated);
  border: 1px dashed var(--border-default);
  border-radius: var(--radius-md);
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
  text-decoration: none;
  transition: all 0.2s ease;
}

.back-link:hover {
  border-color: var(--neon-cyan);
  color: var(--neon-cyan);
  background: rgba(0, 240, 255, 0.05);
}

/* Main Content */
.admin-main {
  flex: 1;
  margin-left: 240px;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

/* Header */
.admin-header {
  height: 60px;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border-default);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 50;
}

.header-left {
  display: flex;
  align-items: center;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-mono);
  font-size: 13px;
}

.breadcrumb-icon {
  color: #f97316;
}

.breadcrumb-root {
  color: var(--text-muted);
}

.breadcrumb-separator {
  color: var(--text-disabled);
}

.breadcrumb-current {
  color: var(--text-primary);
  font-weight: 600;
}

.header-center {
  display: flex;
  align-items: center;
}

.server-time {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 16px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-full);
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-secondary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.notification-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.notification-btn:hover {
  border-color: #f97316;
  color: #f97316;
}

.notification-badge :deep(.el-badge__content) {
  background: #ef4444;
  border: none;
}

/* User Menu */
.user-menu {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 16px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
}

.user-menu:hover {
  border-color: #f97316;
}

.user-avatar {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(249, 115, 22, 0.15);
  border-radius: var(--radius-md);
  color: #f97316;
  font-size: 16px;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
}

.user-role {
  font-family: var(--font-mono);
  font-size: 9px;
  color: #f97316;
  letter-spacing: 0.5px;
}

.dropdown-arrow {
  color: var(--text-muted);
  font-size: 12px;
}

/* Content */
.admin-content {
  flex: 1;
  padding: 24px;
  background: var(--bg-void);
}

/* Footer */
.admin-footer {
  height: 40px;
  background: var(--bg-surface);
  border-top: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}

.footer-left,
.footer-right {
  display: flex;
  align-items: center;
  gap: 12px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
}

.footer-version {
  color: #f97316;
}

.footer-separator {
  color: var(--border-default);
}

/* Dropdown Menu Customization */
.admin-dropdown {
  background: var(--bg-card) !important;
  border: 1px solid var(--border-default) !important;
}

.admin-dropdown :deep(.el-dropdown-menu__item) {
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 8px;
}

.admin-dropdown :deep(.el-dropdown-menu__item:hover) {
  background: var(--bg-elevated);
  color: #f97316;
}

/* Responsive */
@media (max-width: 1024px) {
  .admin-sidebar {
    width: 200px;
  }

  .admin-main {
    margin-left: 200px;
  }

  .nav-code {
    display: none;
  }
}

@media (max-width: 768px) {
  .admin-sidebar {
    transform: translateX(-100%);
    transition: transform 0.3s ease;
  }

  .admin-sidebar.open {
    transform: translateX(0);
  }

  .admin-main {
    margin-left: 0;
  }
}
</style>
