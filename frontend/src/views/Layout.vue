<template>
  <div class="app-layout">
    <!-- Sidebar -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <!-- Logo -->
      <div class="sidebar-header">
        <div class="logo">
          <div class="logo-icon">
            <svg viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M20 4L36 12V28L20 36L4 28V12L20 4Z" stroke="currentColor" stroke-width="2" fill="none"/>
              <circle cx="20" cy="20" r="6" stroke="currentColor" stroke-width="2"/>
            </svg>
          </div>
          <transition name="fade">
            <span v-if="!sidebarCollapsed" class="logo-text">
              <span class="text-primary">REMOTE</span>
              <span class="text-accent">DEBUG</span>
            </span>
          </transition>
        </div>
        <button class="collapse-btn" @click="sidebarCollapsed = !sidebarCollapsed">
          <el-icon><Fold v-if="!sidebarCollapsed" /><Expand v-else /></el-icon>
        </button>
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
          <span class="nav-icon">
            <el-icon><component :is="item.icon" /></el-icon>
          </span>
          <transition name="fade">
            <span v-if="!sidebarCollapsed" class="nav-label">{{ item.label }}</span>
          </transition>
          <span v-if="!sidebarCollapsed && item.badge" class="nav-badge">{{ item.badge }}</span>
        </router-link>
      </nav>

      <!-- Sidebar footer -->
      <div class="sidebar-footer">
        <router-link
          v-if="userStore.user?.role === 'admin'"
          to="/admin"
          class="nav-item admin-link"
        >
          <span class="nav-icon">
            <el-icon><Setting /></el-icon>
          </span>
          <transition name="fade">
            <span v-if="!sidebarCollapsed" class="nav-label">管理后台</span>
          </transition>
        </router-link>
      </div>
    </aside>

    <!-- Main content area -->
    <div class="main-wrapper">
      <!-- Top header -->
      <header class="top-header">
        <div class="header-left">
          <div class="page-title">
            <span class="title-icon">&#62;</span>
            <span class="title-text">{{ currentPageTitle }}</span>
          </div>
        </div>

        <div class="header-right">
          <!-- Balance display -->
          <div class="balance-display">
            <div class="balance-icon">
              <el-icon><Wallet /></el-icon>
            </div>
            <div class="balance-info">
              <span class="balance-label">BALANCE</span>
              <span class="balance-value">¥{{ userStore.user?.balance?.toFixed(2) || '0.00' }}</span>
            </div>
          </div>

          <!-- User dropdown -->
          <el-dropdown @command="handleCommand" trigger="click">
            <div class="user-avatar">
              <div class="avatar-ring">
                <div class="avatar-inner">
                  {{ userStore.user?.username?.charAt(0).toUpperCase() }}
                </div>
              </div>
              <transition name="fade">
                <span class="user-name">{{ userStore.user?.username }}</span>
              </transition>
              <el-icon class="dropdown-arrow"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu class="user-dropdown">
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  <span>个人中心</span>
                </el-dropdown-item>
                <el-dropdown-item command="orders">
                  <el-icon><Wallet /></el-icon>
                  <span>消费记录</span>
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  <span>退出登录</span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>

      <!-- Page content -->
      <main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>

      <!-- Status bar -->
      <footer class="status-bar">
        <div class="status-left">
          <span class="status-item">
            <span class="status-dot" :class="connectionStatus"></span>
            <span>{{ connectionStatus === 'online' ? 'CONNECTED' : 'DISCONNECTED' }}</span>
          </span>
        </div>
        <div class="status-center">
          <span class="status-item">
            <el-icon><Connection /></el-icon>
            <span>{{ stats.ports || 0 }} PORTS</span>
          </span>
          <span class="status-divider">|</span>
          <span class="status-item">
            <el-icon><Iphone /></el-icon>
            <span>{{ stats.devices || 0 }} DEVICES</span>
          </span>
        </div>
        <div class="status-right">
          <span class="status-item">{{ currentTime }}</span>
        </div>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import {
  HomeFilled,
  Connection,
  Iphone,
  Wallet,
  User,
  Setting,
  Fold,
  Expand,
  ArrowDown,
  SwitchButton
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const sidebarCollapsed = ref(false)
const currentTime = ref('')
const connectionStatus = ref('online')
const stats = ref({ ports: 0, devices: 0 })

const navItems = [
  { path: '/', label: '控制台', icon: HomeFilled },
  { path: '/ports', label: '端口管理', icon: Connection },
  { path: '/devices', label: '设备管理', icon: Iphone },
  { path: '/orders', label: '消费记录', icon: Wallet },
  { path: '/profile', label: '个人中心', icon: User }
]

const pageTitles = {
  '/': 'DASHBOARD',
  '/ports': 'PORT_MANAGER',
  '/devices': 'DEVICE_MANAGER',
  '/orders': 'TRANSACTIONS',
  '/profile': 'PROFILE'
}

const currentPageTitle = computed(() => {
  return pageTitles[route.path] || 'DASHBOARD'
})

function isActive(path) {
  if (path === '/') {
    return route.path === '/'
  }
  return route.path.startsWith(path)
}

function handleCommand(command) {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'orders':
      router.push('/orders')
      break
    case 'logout':
      userStore.logout()
      router.push('/login')
      break
  }
}

function updateTime() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour12: false })
}

async function fetchStats() {
  try {
    const res = await userStore.refreshProfile()
    if (res?.data) {
      stats.value = {
        ports: res.data.active_port_count || 0,
        devices: res.data.online_device_count || 0
      }
    }
  } catch (e) {
    connectionStatus.value = 'offline'
  }
}

let timeInterval
let statsInterval

onMounted(() => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
  fetchStats()
  statsInterval = setInterval(fetchStats, 30000)
})

onUnmounted(() => {
  clearInterval(timeInterval)
  clearInterval(statsInterval)
})
</script>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
  background: var(--bg-void);
}

/* Sidebar */
.sidebar {
  width: 240px;
  background: var(--bg-surface);
  border-right: 1px solid var(--border-subtle);
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
  position: relative;
  z-index: 100;
}

.sidebar.collapsed {
  width: 72px;
}

/* Sidebar header */
.sidebar-header {
  height: 64px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--border-subtle);
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  overflow: hidden;
}

.logo-icon {
  width: 36px;
  height: 36px;
  color: var(--neon-cyan);
  flex-shrink: 0;
}

.logo-text {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 700;
  letter-spacing: 1px;
  white-space: nowrap;
}

.logo-text .text-primary {
  color: var(--text-primary);
}

.logo-text .text-accent {
  color: var(--neon-cyan);
}

.collapse-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.2s ease;
}

.collapse-btn:hover {
  color: var(--neon-cyan);
  border-color: var(--neon-cyan);
}

.sidebar.collapsed .collapse-btn {
  margin-left: auto;
  margin-right: auto;
}

/* Navigation */
.sidebar-nav {
  flex: 1;
  padding: 16px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  text-decoration: none;
  transition: all 0.2s ease;
  position: relative;
}

.nav-item:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.nav-item.active {
  background: rgba(0, 240, 255, 0.1);
  color: var(--neon-cyan);
}

.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 24px;
  background: var(--neon-cyan);
  border-radius: 0 2px 2px 0;
  box-shadow: 0 0 10px var(--neon-cyan);
}

.nav-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.nav-label {
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
}

.nav-badge {
  margin-left: auto;
  padding: 2px 8px;
  background: var(--neon-cyan);
  color: var(--bg-void);
  font-size: 11px;
  font-weight: 600;
  border-radius: 10px;
}

.sidebar.collapsed .nav-item {
  justify-content: center;
  padding: 12px;
}

/* Sidebar footer */
.sidebar-footer {
  padding: 12px;
  border-top: 1px solid var(--border-subtle);
}

.admin-link {
  background: rgba(168, 85, 247, 0.1);
  border: 1px solid rgba(168, 85, 247, 0.2);
}

.admin-link:hover {
  background: rgba(168, 85, 247, 0.2);
  border-color: var(--neon-purple);
}

.admin-link .nav-icon {
  color: var(--neon-purple);
}

/* Main wrapper */
.main-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

/* Top header */
.top-header {
  height: 64px;
  padding: 0 24px;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-mono);
  font-size: 14px;
  font-weight: 500;
}

.title-icon {
  color: var(--neon-cyan);
}

.title-text {
  color: var(--text-primary);
  letter-spacing: 1px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

/* Balance display */
.balance-display {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 16px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
}

.balance-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(251, 191, 36, 0.15);
  border-radius: var(--radius-sm);
  color: var(--neon-yellow);
}

.balance-info {
  display: flex;
  flex-direction: column;
}

.balance-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.balance-value {
  font-family: var(--font-mono);
  font-size: 16px;
  font-weight: 600;
  color: var(--neon-yellow);
}

/* User avatar */
.user-avatar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 12px 6px 6px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: 100px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.user-avatar:hover {
  border-color: var(--neon-cyan-dim);
}

.avatar-ring {
  width: 36px;
  height: 36px;
  padding: 2px;
  background: linear-gradient(135deg, var(--neon-cyan), var(--neon-purple));
  border-radius: 50%;
}

.avatar-inner {
  width: 100%;
  height: 100%;
  background: var(--bg-card);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-display);
  font-size: 14px;
  font-weight: 700;
  color: var(--neon-cyan);
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.dropdown-arrow {
  color: var(--text-muted);
  transition: transform 0.2s ease;
}

.user-avatar:hover .dropdown-arrow {
  color: var(--neon-cyan);
}

/* Main content */
.main-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  background: var(--bg-deep);
}

/* Status bar */
.status-bar {
  height: 32px;
  padding: 0 24px;
  background: var(--bg-surface);
  border-top: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
}

.status-left,
.status-center,
.status-right {
  display: flex;
  align-items: center;
  gap: 12px;
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
}

.status-dot.online {
  background: var(--neon-green);
  box-shadow: 0 0 8px var(--neon-green);
}

.status-dot.offline {
  background: var(--neon-red);
  box-shadow: 0 0 8px var(--neon-red);
}

.status-divider {
  color: var(--border-default);
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.page-enter-active {
  animation: page-in 0.3s ease-out;
}

.page-leave-active {
  animation: page-out 0.2s ease-in;
}

@keyframes page-in {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes page-out {
  from {
    opacity: 1;
    transform: translateY(0);
  }
  to {
    opacity: 0;
    transform: translateY(-10px);
  }
}

/* Dropdown menu override */
:deep(.user-dropdown) {
  min-width: 160px;
}

:deep(.user-dropdown .el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
}

/* Responsive */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    z-index: 1000;
    transform: translateX(-100%);
  }

  .sidebar:not(.collapsed) {
    transform: translateX(0);
  }

  .balance-display {
    display: none;
  }

  .user-name {
    display: none;
  }
}
</style>
