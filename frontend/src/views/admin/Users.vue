<template>
  <div class="admin-users-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
        </div>
        <div class="header-text">
          <h1 class="page-title">
            <span class="title-prefix">//</span>
            USER_MANAGEMENT
          </h1>
          <p class="page-subtitle">系统用户管理与权限控制</p>
        </div>
      </div>
      <div class="header-actions">
        <div class="search-box">
          <el-input
            v-model="keyword"
            placeholder="搜索用户名或邮箱..."
            @keyup.enter="fetchUsers"
            clearable
            class="search-input"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <button class="search-btn" @click="fetchUsers">
            <el-icon><Search /></el-icon>
            SCAN
          </button>
        </div>
      </div>
    </div>

    <!-- Stats Bar -->
    <div class="stats-bar">
      <div class="stat-chip">
        <span class="stat-label">TOTAL_USERS</span>
        <span class="stat-value">{{ total }}</span>
      </div>
      <div class="stat-chip active">
        <span class="stat-label">ACTIVE</span>
        <span class="stat-value">{{ activeCount }}</span>
      </div>
      <div class="stat-chip disabled">
        <span class="stat-label">DISABLED</span>
        <span class="stat-value">{{ disabledCount }}</span>
      </div>
      <div class="stat-chip admin">
        <span class="stat-label">ADMINS</span>
        <span class="stat-value">{{ adminCount }}</span>
      </div>
    </div>

    <!-- Users Table -->
    <div class="table-container">
      <div class="table-header">
        <div class="table-title">
          <span class="title-dot"></span>
          USER_DATABASE
        </div>
        <div class="table-actions">
          <button class="action-btn refresh" @click="fetchUsers">
            <el-icon><Refresh /></el-icon>
            REFRESH
          </button>
        </div>
      </div>

      <el-table
        :data="users"
        v-loading="loading"
        class="cyber-table"
        :header-cell-style="{ background: 'var(--bg-elevated)', color: 'var(--text-muted)', fontFamily: 'var(--font-mono)', fontSize: '11px', fontWeight: '600', letterSpacing: '1px', borderBottom: '1px solid var(--border-default)' }"
        :row-style="{ background: 'transparent' }"
      >
        <el-table-column prop="id" label="ID" width="70">
          <template #default="{ row }">
            <span class="id-badge">#{{ row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="USER_NAME" width="140">
          <template #default="{ row }">
            <div class="user-cell">
              <div class="user-avatar">
                {{ row.username.charAt(0).toUpperCase() }}
              </div>
              <span class="username">{{ row.username }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="EMAIL" width="200">
          <template #default="{ row }">
            <span class="email-text">{{ row.email || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="balance" label="BALANCE" width="120">
          <template #default="{ row }">
            <span class="balance-value">¥{{ row.balance?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="max_ports" label="MAX_PORTS" width="100">
          <template #default="{ row }">
            <span class="mono-value">{{ row.max_ports }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="role" label="ROLE" width="110">
          <template #default="{ row }">
            <span :class="['role-badge', row.role]">
              {{ row.role === 'admin' ? 'ADMIN' : 'USER' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="STATUS" width="100">
          <template #default="{ row }">
            <span :class="['status-indicator', row.status]">
              <span class="status-dot"></span>
              {{ row.status === 'active' ? 'ACTIVE' : 'DISABLED' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="CREATED" width="150">
          <template #default="{ row }">
            <span class="time-text">{{ formatTime(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="ACTIONS" width="220" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <button class="btn-action recharge" @click="showRechargeDialog(row)">
                <el-icon><Wallet /></el-icon>
                充值
              </button>
              <button class="btn-action edit" @click="showEditDialog(row)">
                <el-icon><Edit /></el-icon>
                编辑
              </button>
              <button
                :class="['btn-action', row.status === 'active' ? 'disable' : 'enable']"
                @click="toggleStatus(row)"
              >
                <el-icon><component :is="row.status === 'active' ? 'Lock' : 'Unlock'" /></el-icon>
                {{ row.status === 'active' ? '禁用' : '启用' }}
              </button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination -->
      <div class="pagination-container">
        <div class="pagination-info">
          显示 {{ (page - 1) * pageSize + 1 }}-{{ Math.min(page * pageSize, total) }} / {{ total }} 条记录
        </div>
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="sizes, prev, pager, next"
          @size-change="fetchUsers"
          @current-change="fetchUsers"
          class="cyber-pagination"
        />
      </div>

      <!-- Corner Decorations -->
      <div class="corner-decor top-left"></div>
      <div class="corner-decor top-right"></div>
      <div class="corner-decor bottom-left"></div>
      <div class="corner-decor bottom-right"></div>
    </div>

    <!-- Recharge Dialog -->
    <el-dialog
      v-model="rechargeDialogVisible"
      title=""
      width="460"
      class="cyber-dialog"
      :show-close="false"
    >
      <template #header>
        <div class="dialog-header">
          <div class="dialog-icon recharge">
            <el-icon><Wallet /></el-icon>
          </div>
          <div class="dialog-title">
            <h3>BALANCE_RECHARGE</h3>
            <p>用户余额充值</p>
          </div>
          <button class="dialog-close" @click="rechargeDialogVisible = false">
            <el-icon><Close /></el-icon>
          </button>
        </div>
      </template>

      <div class="dialog-content">
        <div class="info-row">
          <span class="info-label">TARGET_USER</span>
          <span class="info-value user">
            <span class="user-badge">{{ selectedUser?.username }}</span>
          </span>
        </div>
        <div class="info-row">
          <span class="info-label">CURRENT_BALANCE</span>
          <span class="info-value balance">¥{{ selectedUser?.balance?.toFixed(2) }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">
            <span class="label-icon">></span>
            AMOUNT
          </label>
          <div class="amount-input-wrapper">
            <span class="currency-prefix">¥</span>
            <el-input-number
              v-model="rechargeAmount"
              :min="0"
              :precision="2"
              :controls="false"
              class="amount-input"
            />
          </div>
          <div class="quick-amounts">
            <button
              v-for="amt in [50, 100, 200, 500]"
              :key="amt"
              class="quick-btn"
              @click="rechargeAmount = amt"
            >
              ¥{{ amt }}
            </button>
          </div>
        </div>

        <div class="form-group">
          <label class="form-label">
            <span class="label-icon">></span>
            DESCRIPTION
          </label>
          <el-input
            v-model="rechargeDescription"
            placeholder="备注信息（可选）"
            class="cyber-input"
          />
        </div>

        <div class="preview-row">
          <span class="preview-label">充值后余额</span>
          <span class="preview-value">¥{{ ((selectedUser?.balance || 0) + rechargeAmount).toFixed(2) }}</span>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <button class="btn-cancel" @click="rechargeDialogVisible = false">取消</button>
          <button class="btn-confirm" @click="handleRecharge" :disabled="rechargeLoading">
            <el-icon v-if="rechargeLoading" class="is-loading"><Loading /></el-icon>
            <span>{{ rechargeLoading ? '处理中...' : '确认充值' }}</span>
          </button>
        </div>
      </template>
    </el-dialog>

    <!-- Edit Dialog -->
    <el-dialog
      v-model="editDialogVisible"
      title=""
      width="520"
      class="cyber-dialog"
      :show-close="false"
    >
      <template #header>
        <div class="dialog-header">
          <div class="dialog-icon edit">
            <el-icon><Edit /></el-icon>
          </div>
          <div class="dialog-title">
            <h3>USER_CONFIG</h3>
            <p>编辑用户信息</p>
          </div>
          <button class="dialog-close" @click="editDialogVisible = false">
            <el-icon><Close /></el-icon>
          </button>
        </div>
      </template>

      <div class="dialog-content">
        <div class="form-group">
          <label class="form-label">
            <span class="label-icon">></span>
            USERNAME
          </label>
          <el-input v-model="editForm.username" disabled class="cyber-input" />
        </div>

        <div class="form-group">
          <label class="form-label">
            <span class="label-icon">></span>
            EMAIL
          </label>
          <el-input v-model="editForm.email" placeholder="用户邮箱" class="cyber-input" />
        </div>

        <div class="form-row">
          <div class="form-group half">
            <label class="form-label">
              <span class="label-icon">></span>
              MAX_PORTS
            </label>
            <el-input-number v-model="editForm.max_ports" :min="1" class="cyber-number" />
          </div>
          <div class="form-group half">
            <label class="form-label">
              <span class="label-icon">></span>
              MAX_DEVICES
            </label>
            <el-input-number v-model="editForm.max_devices" :min="1" class="cyber-number" />
          </div>
        </div>

        <div class="form-group">
          <label class="form-label">
            <span class="label-icon">></span>
            ROLE
          </label>
          <div class="role-selector">
            <button
              :class="['role-option', { active: editForm.role === 'user' }]"
              @click="editForm.role = 'user'"
            >
              <el-icon><User /></el-icon>
              普通用户
            </button>
            <button
              :class="['role-option admin', { active: editForm.role === 'admin' }]"
              @click="editForm.role = 'admin'"
            >
              <el-icon><UserFilled /></el-icon>
              管理员
            </button>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <button class="btn-cancel" @click="editDialogVisible = false">取消</button>
          <button class="btn-confirm" @click="handleEdit" :disabled="editLoading">
            <el-icon v-if="editLoading" class="is-loading"><Loading /></el-icon>
            <span>{{ editLoading ? '保存中...' : '保存修改' }}</span>
          </button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Search, Refresh, Wallet, Edit, Lock, Unlock, Close,
  Loading, User, UserFilled
} from '@element-plus/icons-vue'
import api from '../../api'
import dayjs from 'dayjs'

const loading = ref(false)
const users = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const keyword = ref('')

// Computed stats
const activeCount = computed(() => users.value.filter(u => u.status === 'active').length)
const disabledCount = computed(() => users.value.filter(u => u.status === 'disabled').length)
const adminCount = computed(() => users.value.filter(u => u.role === 'admin').length)

// Recharge
const rechargeDialogVisible = ref(false)
const rechargeLoading = ref(false)
const selectedUser = ref(null)
const rechargeAmount = ref(100)
const rechargeDescription = ref('')

// Edit
const editDialogVisible = ref(false)
const editLoading = ref(false)
const editForm = reactive({
  username: '',
  email: '',
  max_ports: 3,
  max_devices: 10,
  role: 'user'
})

function formatTime(time) {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-'
}

async function fetchUsers() {
  loading.value = true
  try {
    const res = await api.admin.getUsers(page.value, pageSize.value, keyword.value)
    if (res.code === 0) {
      users.value = res.data.list || []
      total.value = res.data.total
    }
  } finally {
    loading.value = false
  }
}

function showRechargeDialog(user) {
  selectedUser.value = user
  rechargeAmount.value = 100
  rechargeDescription.value = ''
  rechargeDialogVisible.value = true
}

async function handleRecharge() {
  rechargeLoading.value = true
  try {
    const res = await api.admin.rechargeUser(selectedUser.value.id, rechargeAmount.value, rechargeDescription.value)
    if (res.code === 0) {
      ElMessage.success('充值成功')
      rechargeDialogVisible.value = false
      fetchUsers()
    } else {
      ElMessage.error(res.message)
    }
  } finally {
    rechargeLoading.value = false
  }
}

function showEditDialog(user) {
  selectedUser.value = user
  editForm.username = user.username
  editForm.email = user.email
  editForm.max_ports = user.max_ports
  editForm.max_devices = user.max_devices
  editForm.role = user.role
  editDialogVisible.value = true
}

async function handleEdit() {
  editLoading.value = true
  try {
    const res = await api.admin.updateUser(selectedUser.value.id, {
      email: editForm.email,
      max_ports: editForm.max_ports,
      max_devices: editForm.max_devices,
      role: editForm.role
    })
    if (res.code === 0) {
      ElMessage.success('保存成功')
      editDialogVisible.value = false
      fetchUsers()
    } else {
      ElMessage.error(res.message)
    }
  } finally {
    editLoading.value = false
  }
}

async function toggleStatus(user) {
  const newStatus = user.status === 'active' ? 'disabled' : 'active'
  const res = await api.admin.updateUserStatus(user.id, newStatus)
  if (res.code === 0) {
    ElMessage.success('操作成功')
    fetchUsers()
  } else {
    ElMessage.error(res.message)
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.admin-users-page {
  padding: 0;
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
  background: linear-gradient(135deg, rgba(249, 115, 22, 0.2), rgba(249, 115, 22, 0.05));
  border: 1px solid rgba(249, 115, 22, 0.3);
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--neon-orange);
}

.header-icon svg {
  width: 24px;
  height: 24px;
}

.page-title {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: 1px;
}

.title-prefix {
  color: var(--neon-orange);
  margin-right: 8px;
}

.page-subtitle {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
  margin: 4px 0 0;
}

/* Search Box */
.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
}

.search-input {
  width: 260px;
}

.search-input :deep(.el-input__wrapper) {
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  padding: 8px 12px;
}

.search-input :deep(.el-input__wrapper:focus-within) {
  border-color: var(--neon-orange);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
}

.search-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  background: linear-gradient(135deg, var(--neon-orange), #ea580c);
  border: none;
  border-radius: var(--radius-md);
  color: white;
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.search-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

/* Stats Bar */
.stats-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.stat-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  font-family: var(--font-mono);
}

.stat-label {
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.stat-value {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.stat-chip.active .stat-value { color: var(--neon-green); }
.stat-chip.disabled .stat-value { color: var(--text-muted); }
.stat-chip.admin .stat-value { color: var(--neon-orange); }

/* Table Container */
.table-container {
  position: relative;
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-xl);
  padding: 20px;
}

.table-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-subtle);
}

.table-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
  letter-spacing: 1px;
}

.title-dot {
  width: 8px;
  height: 8px;
  background: var(--neon-orange);
  border-radius: 50%;
  box-shadow: 0 0 8px var(--neon-orange);
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn:hover {
  border-color: var(--neon-orange);
  color: var(--neon-orange);
}

/* Table Styles */
.cyber-table {
  --el-table-bg-color: transparent;
  --el-table-tr-bg-color: transparent;
  --el-table-border-color: var(--border-subtle);
  --el-table-text-color: var(--text-primary);
}

.cyber-table :deep(.el-table__body-wrapper) {
  background: transparent;
}

.cyber-table :deep(.el-table__row) {
  transition: all 0.2s ease;
}

.cyber-table :deep(.el-table__row:hover > td) {
  background: rgba(249, 115, 22, 0.05) !important;
}

.cyber-table :deep(.el-table__cell) {
  border-bottom: 1px solid var(--border-subtle);
  padding: 12px 8px;
}

/* Table Cell Styles */
.id-badge {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, var(--neon-orange), #ea580c);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-display);
  font-size: 14px;
  font-weight: 700;
  color: white;
}

.username {
  font-weight: 500;
  color: var(--text-primary);
}

.email-text {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-secondary);
}

.balance-value {
  font-family: var(--font-mono);
  font-weight: 600;
  color: var(--neon-green);
}

.mono-value {
  font-family: var(--font-mono);
  color: var(--text-secondary);
}

.role-badge {
  display: inline-flex;
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 1px;
}

.role-badge.user {
  background: rgba(100, 116, 139, 0.2);
  color: var(--text-secondary);
  border: 1px solid rgba(100, 116, 139, 0.3);
}

.role-badge.admin {
  background: rgba(249, 115, 22, 0.2);
  color: var(--neon-orange);
  border: 1px solid rgba(249, 115, 22, 0.3);
}

.status-indicator {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 500;
}

.status-indicator .status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.status-indicator.active {
  color: var(--neon-green);
}

.status-indicator.active .status-dot {
  background: var(--neon-green);
  box-shadow: 0 0 8px var(--neon-green);
}

.status-indicator.disabled {
  color: var(--text-muted);
}

.status-indicator.disabled .status-dot {
  background: var(--text-muted);
}

.time-text {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
}

/* Action Buttons */
.action-buttons {
  display: flex;
  gap: 6px;
}

.btn-action {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 5px 10px;
  border: 1px solid;
  border-radius: var(--radius-sm);
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  background: transparent;
}

.btn-action.recharge {
  border-color: var(--neon-green);
  color: var(--neon-green);
}

.btn-action.recharge:hover {
  background: rgba(34, 197, 94, 0.1);
}

.btn-action.edit {
  border-color: var(--neon-cyan);
  color: var(--neon-cyan);
}

.btn-action.edit:hover {
  background: rgba(0, 240, 255, 0.1);
}

.btn-action.disable {
  border-color: #ef4444;
  color: #ef4444;
}

.btn-action.disable:hover {
  background: rgba(239, 68, 68, 0.1);
}

.btn-action.enable {
  border-color: var(--neon-green);
  color: var(--neon-green);
}

.btn-action.enable:hover {
  background: rgba(34, 197, 94, 0.1);
}

/* Pagination */
.pagination-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--border-subtle);
}

.pagination-info {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
}

.cyber-pagination :deep(.el-pagination__sizes),
.cyber-pagination :deep(.el-pager li),
.cyber-pagination :deep(.btn-prev),
.cyber-pagination :deep(.btn-next) {
  background: var(--bg-elevated) !important;
  color: var(--text-secondary) !important;
  border: 1px solid var(--border-default) !important;
}

.cyber-pagination :deep(.el-pager li.is-active) {
  background: var(--neon-orange) !important;
  border-color: var(--neon-orange) !important;
  color: white !important;
}

/* Corner Decorations */
.corner-decor {
  position: absolute;
  width: 16px;
  height: 16px;
  border-color: var(--neon-orange);
  opacity: 0.4;
}

.corner-decor.top-left { top: 8px; left: 8px; border-top: 2px solid; border-left: 2px solid; }
.corner-decor.top-right { top: 8px; right: 8px; border-top: 2px solid; border-right: 2px solid; }
.corner-decor.bottom-left { bottom: 8px; left: 8px; border-bottom: 2px solid; border-left: 2px solid; }
.corner-decor.bottom-right { bottom: 8px; right: 8px; border-bottom: 2px solid; border-right: 2px solid; }

/* Dialog Styles */
.cyber-dialog :deep(.el-dialog) {
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-xl);
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.5);
}

.cyber-dialog :deep(.el-dialog__header) {
  padding: 0;
  margin: 0;
}

.cyber-dialog :deep(.el-dialog__body) {
  padding: 0;
}

.cyber-dialog :deep(.el-dialog__footer) {
  padding: 0;
}

.dialog-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-subtle);
}

.dialog-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.dialog-icon.recharge {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.2), rgba(34, 197, 94, 0.05));
  border: 1px solid rgba(34, 197, 94, 0.3);
  color: var(--neon-green);
}

.dialog-icon.edit {
  background: linear-gradient(135deg, rgba(0, 240, 255, 0.2), rgba(0, 240, 255, 0.05));
  border: 1px solid rgba(0, 240, 255, 0.3);
  color: var(--neon-cyan);
}

.dialog-title {
  flex: 1;
}

.dialog-title h3 {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: 1px;
}

.dialog-title p {
  font-size: 12px;
  color: var(--text-muted);
  margin: 4px 0 0;
}

.dialog-close {
  width: 32px;
  height: 32px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  color: var(--text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.dialog-close:hover {
  border-color: #ef4444;
  color: #ef4444;
}

.dialog-content {
  padding: 24px;
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: var(--bg-elevated);
  border-radius: var(--radius-md);
  margin-bottom: 12px;
}

.info-label {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.info-value {
  font-weight: 600;
}

.info-value.user {
  color: var(--text-primary);
}

.info-value.balance {
  color: var(--neon-green);
  font-family: var(--font-mono);
}

.user-badge {
  padding: 4px 12px;
  background: rgba(249, 115, 22, 0.1);
  border: 1px solid rgba(249, 115, 22, 0.3);
  border-radius: var(--radius-sm);
  font-family: var(--font-mono);
}

.form-group {
  margin-bottom: 20px;
}

.form-group.half {
  flex: 1;
}

.form-row {
  display: flex;
  gap: 16px;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 500;
  color: var(--text-muted);
  letter-spacing: 1px;
  margin-bottom: 8px;
}

.label-icon {
  color: var(--neon-orange);
}

.amount-input-wrapper {
  display: flex;
  align-items: center;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  padding: 0 12px;
  transition: all 0.2s ease;
}

.amount-input-wrapper:focus-within {
  border-color: var(--neon-orange);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
}

.currency-prefix {
  font-family: var(--font-mono);
  font-size: 16px;
  font-weight: 600;
  color: var(--neon-green);
  margin-right: 8px;
}

.amount-input {
  flex: 1;
}

.amount-input :deep(.el-input__wrapper) {
  background: transparent !important;
  box-shadow: none !important;
  padding: 8px 0;
}

.amount-input :deep(.el-input__inner) {
  font-family: var(--font-mono);
  font-size: 20px;
  font-weight: 600;
}

.quick-amounts {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.quick-btn {
  flex: 1;
  padding: 8px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.quick-btn:hover {
  border-color: var(--neon-green);
  color: var(--neon-green);
}

.cyber-input :deep(.el-input__wrapper) {
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  padding: 10px 14px;
}

.cyber-input :deep(.el-input__wrapper:focus-within) {
  border-color: var(--neon-orange);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
}

.cyber-number {
  width: 100%;
}

.cyber-number :deep(.el-input__wrapper) {
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
}

.preview-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.2);
  border-radius: var(--radius-md);
  margin-top: 20px;
}

.preview-label {
  font-size: 13px;
  color: var(--text-secondary);
}

.preview-value {
  font-family: var(--font-mono);
  font-size: 20px;
  font-weight: 700;
  color: var(--neon-green);
}

.role-selector {
  display: flex;
  gap: 12px;
}

.role-option {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px;
  background: var(--bg-elevated);
  border: 2px solid var(--border-default);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.role-option:hover {
  border-color: var(--neon-cyan);
}

.role-option.active {
  border-color: var(--neon-cyan);
  background: rgba(0, 240, 255, 0.1);
  color: var(--neon-cyan);
}

.role-option.admin:hover {
  border-color: var(--neon-orange);
}

.role-option.admin.active {
  border-color: var(--neon-orange);
  background: rgba(249, 115, 22, 0.1);
  color: var(--neon-orange);
}

.dialog-footer {
  display: flex;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid var(--border-subtle);
}

.btn-cancel {
  flex: 1;
  padding: 12px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-cancel:hover {
  border-color: var(--text-muted);
  color: var(--text-primary);
}

.btn-confirm {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  background: linear-gradient(135deg, var(--neon-orange), #ea580c);
  border: none;
  border-radius: var(--radius-md);
  color: white;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-confirm:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.btn-confirm:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.is-loading {
  animation: spin 1s linear infinite;
}
</style>
