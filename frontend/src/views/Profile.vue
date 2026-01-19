<template>
  <div class="profile-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-icon">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <circle cx="12" cy="8" r="4" stroke="currentColor" stroke-width="2"/>
            <path d="M4 20c0-4 4-6 8-6s8 2 8 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </div>
        <div class="header-text">
          <h1 class="page-title">USER_PROFILE</h1>
          <p class="page-subtitle">// 账户信息管理</p>
        </div>
      </div>
      <div class="header-badge">
        <span class="badge-label">ROLE:</span>
        <span class="badge-value" :class="userStore.user?.role">
          {{ userStore.user?.role === 'admin' ? 'ADMIN' : 'USER' }}
        </span>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card balance">
        <div class="stat-card-header">
          <div class="stat-icon">
            <el-icon><Wallet /></el-icon>
          </div>
          <span class="stat-label">BALANCE</span>
        </div>
        <div class="stat-card-value">
          <span class="currency">¥</span>
          <span class="amount">{{ userStore.user?.balance?.toFixed(2) || '0.00' }}</span>
        </div>
        <div class="stat-card-footer">
          <span class="footer-text">可用余额</span>
        </div>
        <div class="card-glow"></div>
      </div>

      <div class="stat-card ports">
        <div class="stat-card-header">
          <div class="stat-icon">
            <el-icon><Connection /></el-icon>
          </div>
          <span class="stat-label">ACTIVE_PORTS</span>
        </div>
        <div class="stat-card-value">
          <span class="amount">{{ stats.active_port_count || 0 }}</span>
          <span class="unit">/ {{ userStore.user?.max_ports || 0 }}</span>
        </div>
        <div class="stat-card-footer">
          <div class="progress-bar">
            <div
              class="progress-fill"
              :style="{ width: portUsagePercent + '%' }"
            ></div>
          </div>
        </div>
        <div class="card-glow"></div>
      </div>

      <div class="stat-card devices">
        <div class="stat-card-header">
          <div class="stat-icon">
            <el-icon><Monitor /></el-icon>
          </div>
          <span class="stat-label">ONLINE_DEVICES</span>
        </div>
        <div class="stat-card-value">
          <span class="amount">{{ stats.online_device_count || 0 }}</span>
          <span class="unit">/ {{ userStore.user?.max_devices || 0 }}</span>
        </div>
        <div class="stat-card-footer">
          <div class="progress-bar">
            <div
              class="progress-fill"
              :style="{ width: deviceUsagePercent + '%' }"
            ></div>
          </div>
        </div>
        <div class="card-glow"></div>
      </div>

      <div class="stat-card consume">
        <div class="stat-card-header">
          <div class="stat-icon">
            <el-icon><TrendCharts /></el-icon>
          </div>
          <span class="stat-label">TOTAL_CONSUME</span>
        </div>
        <div class="stat-card-value">
          <span class="currency">¥</span>
          <span class="amount">{{ stats.total_consume?.toFixed(2) || '0.00' }}</span>
        </div>
        <div class="stat-card-footer">
          <span class="footer-text">累计消费</span>
        </div>
        <div class="card-glow"></div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="content-grid">
      <!-- Profile Info Card -->
      <div class="info-card">
        <div class="card-header">
          <div class="header-icon small">
            <el-icon><User /></el-icon>
          </div>
          <h2 class="card-title">PROFILE_INFO</h2>
          <span class="card-subtitle">// 个人信息</span>
        </div>

        <div class="card-body">
          <el-form :model="profileForm" label-position="top" class="cyber-form">
            <div class="form-row">
              <div class="form-field">
                <label class="field-label">
                  <span class="label-icon">></span>
                  USER_ID
                </label>
                <div class="field-value readonly">
                  <el-icon><User /></el-icon>
                  <span>{{ profileForm.username }}</span>
                  <span class="readonly-badge">LOCKED</span>
                </div>
              </div>
            </div>

            <div class="form-row">
              <div class="form-field">
                <label class="field-label">
                  <span class="label-icon">></span>
                  EMAIL_ADDRESS
                </label>
                <el-input
                  v-model="profileForm.email"
                  placeholder="请输入邮箱地址"
                  size="large"
                  :prefix-icon="Message"
                />
              </div>
            </div>

            <div class="form-row info-row">
              <div class="info-item">
                <span class="info-label">MAX_PORTS:</span>
                <span class="info-value">{{ userStore.user?.max_ports || 0 }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">MAX_DEVICES:</span>
                <span class="info-value">{{ userStore.user?.max_devices || 0 }}</span>
              </div>
            </div>

            <div class="form-row info-row">
              <div class="info-item full">
                <span class="info-label">CREATED_AT:</span>
                <span class="info-value">{{ formatTime(userStore.user?.created_at) }}</span>
              </div>
            </div>

            <div class="form-actions">
              <el-button
                type="primary"
                @click="handleUpdateProfile"
                :loading="profileLoading"
                class="action-btn"
              >
                <el-icon><Check /></el-icon>
                <span>保存修改</span>
              </el-button>
            </div>
          </el-form>
        </div>

        <!-- Corner decorations -->
        <div class="corner-decor top-left"></div>
        <div class="corner-decor top-right"></div>
        <div class="corner-decor bottom-left"></div>
        <div class="corner-decor bottom-right"></div>
      </div>

      <!-- Password Change Card -->
      <div class="info-card password-card">
        <div class="card-header">
          <div class="header-icon small">
            <el-icon><Lock /></el-icon>
          </div>
          <h2 class="card-title">SECURITY_KEY</h2>
          <span class="card-subtitle">// 修改密码</span>
        </div>

        <div class="card-body">
          <el-form
            ref="passwordFormRef"
            :model="passwordForm"
            :rules="passwordRules"
            label-position="top"
            class="cyber-form"
          >
            <div class="form-row">
              <div class="form-field">
                <label class="field-label">
                  <span class="label-icon">></span>
                  CURRENT_KEY
                </label>
                <el-form-item prop="oldPassword">
                  <el-input
                    v-model="passwordForm.oldPassword"
                    type="password"
                    placeholder="输入当前密码"
                    size="large"
                    show-password
                    :prefix-icon="Key"
                  />
                </el-form-item>
              </div>
            </div>

            <div class="form-row">
              <div class="form-field">
                <label class="field-label">
                  <span class="label-icon">></span>
                  NEW_KEY
                </label>
                <el-form-item prop="newPassword">
                  <el-input
                    v-model="passwordForm.newPassword"
                    type="password"
                    placeholder="设置新密码 (至少6位)"
                    size="large"
                    show-password
                    :prefix-icon="Key"
                  />
                </el-form-item>
              </div>
            </div>

            <div class="form-row">
              <div class="form-field">
                <label class="field-label">
                  <span class="label-icon">></span>
                  CONFIRM_KEY
                </label>
                <el-form-item prop="confirmPassword">
                  <el-input
                    v-model="passwordForm.confirmPassword"
                    type="password"
                    placeholder="确认新密码"
                    size="large"
                    show-password
                    :prefix-icon="Key"
                  />
                </el-form-item>
              </div>
            </div>

            <div class="security-tips">
              <div class="tip-header">
                <el-icon><Warning /></el-icon>
                <span>SECURITY_NOTICE</span>
              </div>
              <ul class="tip-list">
                <li>密码长度至少 6 个字符</li>
                <li>建议使用字母、数字和符号的组合</li>
                <li>修改密码后需要重新登录</li>
              </ul>
            </div>

            <div class="form-actions">
              <el-button
                type="primary"
                @click="handleUpdatePassword"
                :loading="passwordLoading"
                class="action-btn warning"
              >
                <el-icon><Unlock /></el-icon>
                <span>更新密码</span>
              </el-button>
            </div>
          </el-form>
        </div>

        <!-- Corner decorations -->
        <div class="corner-decor top-left"></div>
        <div class="corner-decor top-right"></div>
        <div class="corner-decor bottom-left"></div>
        <div class="corner-decor bottom-right"></div>
      </div>
    </div>

    <!-- Account Status -->
    <div class="status-card">
      <div class="status-header">
        <el-icon><CircleCheck /></el-icon>
        <span>ACCOUNT_STATUS</span>
      </div>
      <div class="status-content">
        <div class="status-item">
          <span class="status-label">STATUS:</span>
          <span class="status-value active">
            <span class="status-dot"></span>
            {{ userStore.user?.status === 'active' ? 'ACTIVE' : 'DISABLED' }}
          </span>
        </div>
        <div class="status-item">
          <span class="status-label">LAST_LOGIN:</span>
          <span class="status-value">{{ formatTime(userStore.user?.updated_at) || 'N/A' }}</span>
        </div>
        <div class="status-item">
          <span class="status-label">SESSION:</span>
          <span class="status-value online">CONNECTED</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  User, Lock, Message, Wallet, Connection, Monitor, TrendCharts,
  Check, Key, Warning, Unlock, CircleCheck
} from '@element-plus/icons-vue'
import { useUserStore } from '../stores/user'
import api from '../api'
import dayjs from 'dayjs'

const userStore = useUserStore()

const profileForm = reactive({
  username: userStore.user?.username || '',
  email: userStore.user?.email || ''
})

const profileLoading = ref(false)

const passwordFormRef = ref()
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})
const passwordLoading = ref(false)

const stats = ref({})

// Computed
const portUsagePercent = computed(() => {
  const max = userStore.user?.max_ports || 1
  const active = stats.value.active_port_count || 0
  return Math.min((active / max) * 100, 100)
})

const deviceUsagePercent = computed(() => {
  const max = userStore.user?.max_devices || 1
  const online = stats.value.online_device_count || 0
  return Math.min((online / max) * 100, 100)
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = {
  oldPassword: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

function formatTime(time) {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm:ss') : '-'
}

async function handleUpdateProfile() {
  profileLoading.value = true
  try {
    const res = await api.user.updateProfile({ email: profileForm.email })
    if (res.code === 0) {
      ElMessage.success('保存成功')
      userStore.refreshProfile()
    } else {
      ElMessage.error(res.message)
    }
  } finally {
    profileLoading.value = false
  }
}

async function handleUpdatePassword() {
  await passwordFormRef.value.validate()
  passwordLoading.value = true
  try {
    const res = await api.user.updatePassword(passwordForm.oldPassword, passwordForm.newPassword)
    if (res.code === 0) {
      ElMessage.success('密码修改成功')
      passwordForm.oldPassword = ''
      passwordForm.newPassword = ''
      passwordForm.confirmPassword = ''
    } else {
      ElMessage.error(res.message)
    }
  } finally {
    passwordLoading.value = false
  }
}

async function fetchStats() {
  const res = await api.user.getProfile()
  if (res.code === 0) {
    stats.value = res.data.stats || {}
  }
}

onMounted(() => {
  fetchStats()
  // Sync profile form with user store
  profileForm.username = userStore.user?.username || ''
  profileForm.email = userStore.user?.email || ''
})
</script>

<style scoped>
.profile-page {
  padding: 0;
  animation: page-enter 0.4s ease-out;
}

@keyframes page-enter {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Page Header */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--border-subtle);
}

.header-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(0, 240, 255, 0.2), rgba(0, 240, 255, 0.05));
  border: 1px solid rgba(0, 240, 255, 0.3);
  border-radius: var(--radius-lg);
  color: var(--neon-cyan);
}

.header-icon svg {
  width: 24px;
  height: 24px;
}

.header-icon.small {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-md);
}

.header-icon.small svg,
.header-icon.small .el-icon {
  font-size: 16px;
}

.page-title {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: 1px;
}

.page-subtitle {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
  margin: 4px 0 0 0;
}

.header-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-full);
  font-family: var(--font-mono);
  font-size: 12px;
}

.badge-label {
  color: var(--text-muted);
}

.badge-value {
  font-weight: 600;
}

.badge-value.admin {
  color: var(--neon-purple);
}

.badge-value.user {
  color: var(--neon-cyan);
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  position: relative;
  padding: 20px;
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  overflow: hidden;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  border-color: var(--border-active);
}

.stat-card-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.stat-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  font-size: 16px;
}

.stat-card.balance .stat-icon {
  background: rgba(34, 197, 94, 0.15);
  color: var(--neon-green);
}

.stat-card.ports .stat-icon {
  background: rgba(0, 240, 255, 0.15);
  color: var(--neon-cyan);
}

.stat-card.devices .stat-icon {
  background: rgba(168, 85, 247, 0.15);
  color: var(--neon-purple);
}

.stat-card.consume .stat-icon {
  background: rgba(251, 191, 36, 0.15);
  color: #fbbf24;
}

.stat-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 0.5px;
}

.stat-card-value {
  display: flex;
  align-items: baseline;
  gap: 4px;
  margin-bottom: 12px;
}

.currency {
  font-family: var(--font-mono);
  font-size: 16px;
  color: var(--text-muted);
}

.amount {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-card.balance .amount {
  color: var(--neon-green);
}

.stat-card.consume .amount {
  color: #fbbf24;
}

.unit {
  font-family: var(--font-mono);
  font-size: 14px;
  color: var(--text-muted);
}

.stat-card-footer {
  min-height: 20px;
}

.footer-text {
  font-size: 12px;
  color: var(--text-muted);
}

.progress-bar {
  height: 4px;
  background: var(--bg-elevated);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: var(--radius-full);
  transition: width 0.5s ease;
}

.stat-card.ports .progress-fill {
  background: linear-gradient(90deg, var(--neon-cyan), var(--neon-cyan-dim));
}

.stat-card.devices .progress-fill {
  background: linear-gradient(90deg, var(--neon-purple), var(--neon-purple-dim));
}

.card-glow {
  position: absolute;
  top: 0;
  right: 0;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  filter: blur(40px);
  opacity: 0.15;
  pointer-events: none;
}

.stat-card.balance .card-glow {
  background: var(--neon-green);
}

.stat-card.ports .card-glow {
  background: var(--neon-cyan);
}

.stat-card.devices .card-glow {
  background: var(--neon-purple);
}

.stat-card.consume .card-glow {
  background: #fbbf24;
}

/* Content Grid */
.content-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

/* Info Card */
.info-card {
  position: relative;
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-subtle);
  background: var(--bg-surface);
}

.card-title {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: 0.5px;
}

.card-subtitle {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
  margin-left: auto;
}

.card-body {
  padding: 24px;
}

/* Corner decorations */
.corner-decor {
  position: absolute;
  width: 12px;
  height: 12px;
  border-color: var(--neon-cyan);
  opacity: 0.3;
}

.corner-decor.top-left {
  top: 8px;
  left: 8px;
  border-top: 2px solid;
  border-left: 2px solid;
}

.corner-decor.top-right {
  top: 8px;
  right: 8px;
  border-top: 2px solid;
  border-right: 2px solid;
}

.corner-decor.bottom-left {
  bottom: 8px;
  left: 8px;
  border-bottom: 2px solid;
  border-left: 2px solid;
}

.corner-decor.bottom-right {
  bottom: 8px;
  right: 8px;
  border-bottom: 2px solid;
  border-right: 2px solid;
}

.password-card .corner-decor {
  border-color: var(--neon-purple);
}

.password-card .header-icon {
  background: linear-gradient(135deg, rgba(168, 85, 247, 0.2), rgba(168, 85, 247, 0.05));
  border-color: rgba(168, 85, 247, 0.3);
  color: var(--neon-purple);
}

/* Form Styles */
.cyber-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-row {
  display: flex;
  gap: 16px;
}

.form-field {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.field-label {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 500;
  color: var(--text-muted);
  letter-spacing: 0.5px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.label-icon {
  color: var(--neon-cyan);
}

.password-card .label-icon {
  color: var(--neon-purple);
}

.field-value {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  font-family: var(--font-mono);
  font-size: 14px;
  color: var(--text-primary);
}

.field-value.readonly {
  background: var(--bg-surface);
  border-style: dashed;
}

.field-value .el-icon {
  color: var(--text-muted);
}

.readonly-badge {
  margin-left: auto;
  padding: 2px 8px;
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  font-size: 9px;
  color: var(--text-disabled);
  letter-spacing: 0.5px;
}

.cyber-form :deep(.el-form-item) {
  margin-bottom: 0;
}

.cyber-form :deep(.el-input__wrapper) {
  padding: 10px 14px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
}

.cyber-form :deep(.el-input__wrapper:hover) {
  border-color: var(--neon-cyan-dim);
}

.password-card .cyber-form :deep(.el-input__wrapper:hover) {
  border-color: var(--neon-purple-dim);
}

.cyber-form :deep(.el-input__wrapper.is-focus) {
  border-color: var(--neon-cyan);
  box-shadow: 0 0 0 3px rgba(0, 240, 255, 0.1);
}

.password-card .cyber-form :deep(.el-input__wrapper.is-focus) {
  border-color: var(--neon-purple);
  box-shadow: 0 0 0 3px rgba(168, 85, 247, 0.1);
}

.cyber-form :deep(.el-input__inner) {
  font-family: var(--font-mono);
  font-size: 14px;
}

/* Info Row */
.info-row {
  padding: 12px 16px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.info-item.full {
  flex: none;
  width: 100%;
}

.info-label {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
}

.info-value {
  font-family: var(--font-mono);
  font-size: 13px;
  color: var(--neon-cyan);
}

/* Security Tips */
.security-tips {
  padding: 16px;
  background: rgba(251, 191, 36, 0.05);
  border: 1px solid rgba(251, 191, 36, 0.2);
  border-radius: var(--radius-md);
}

.tip-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: #fbbf24;
  letter-spacing: 0.5px;
}

.tip-list {
  margin: 0;
  padding-left: 20px;
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.8;
}

.tip-list li::marker {
  color: #fbbf24;
}

/* Form Actions */
.form-actions {
  display: flex;
  justify-content: flex-end;
  padding-top: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 24px;
  font-weight: 500;
  background: linear-gradient(135deg, var(--neon-cyan), var(--neon-cyan-dim));
  border: none;
}

.action-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(0, 240, 255, 0.3);
}

.action-btn.warning {
  background: linear-gradient(135deg, var(--neon-purple), var(--neon-purple-dim));
}

.action-btn.warning:hover {
  box-shadow: 0 6px 20px rgba(168, 85, 247, 0.3);
}

/* Status Card */
.status-card {
  padding: 16px 24px;
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
}

.status-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--neon-green);
}

.status-content {
  display: flex;
  gap: 32px;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-label {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
}

.status-value {
  display: flex;
  align-items: center;
  gap: 6px;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-secondary);
}

.status-value.active,
.status-value.online {
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

/* Responsive */
@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .content-grid {
    grid-template-columns: 1fr;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .status-content {
    flex-direction: column;
    gap: 12px;
  }

  .info-row {
    flex-direction: column;
    gap: 12px;
  }
}
</style>
