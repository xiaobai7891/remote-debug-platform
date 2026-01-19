<template>
  <div class="admin-ports-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M12 2L2 7l10 5 10-5-10-5z"/>
            <path d="M2 17l10 5 10-5"/>
            <path d="M2 12l10 5 10-5"/>
          </svg>
        </div>
        <div class="header-text">
          <h1 class="page-title">
            <span class="title-prefix">//</span>
            PORT_MANAGEMENT
          </h1>
          <p class="page-subtitle">系统端口分配与监控</p>
        </div>
      </div>
      <div class="header-actions">
        <div class="filter-group">
          <button
            v-for="filter in statusFilters"
            :key="filter.value"
            :class="['filter-btn', { active: statusFilter === filter.value }]"
            @click="statusFilter = filter.value; fetchPorts()"
          >
            <span class="filter-dot" :style="{ background: filter.color }"></span>
            {{ filter.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- Stats Bar -->
    <div class="stats-bar">
      <div class="stat-chip">
        <span class="stat-label">TOTAL_PORTS</span>
        <span class="stat-value">{{ total }}</span>
      </div>
      <div class="stat-chip active">
        <span class="stat-label">RUNNING</span>
        <span class="stat-value">{{ runningCount }}</span>
      </div>
      <div class="stat-chip expired">
        <span class="stat-label">EXPIRED</span>
        <span class="stat-value">{{ expiredCount }}</span>
      </div>
      <div class="stat-chip disabled">
        <span class="stat-label">DISABLED</span>
        <span class="stat-value">{{ disabledCount }}</span>
      </div>
    </div>

    <!-- Ports Table -->
    <div class="table-container">
      <div class="table-header">
        <div class="table-title">
          <span class="title-dot"></span>
          PORT_REGISTRY
        </div>
        <div class="table-actions">
          <button class="action-btn refresh" @click="fetchPorts">
            <el-icon><Refresh /></el-icon>
            REFRESH
          </button>
        </div>
      </div>

      <el-table
        :data="ports"
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
        <el-table-column prop="port" label="PORT" width="120">
          <template #default="{ row }">
            <div class="port-cell">
              <span class="port-number">:{{ row.port }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="OWNER" width="130">
          <template #default="{ row }">
            <div class="user-cell">
              <div class="user-avatar-small">
                {{ row.username?.charAt(0).toUpperCase() || '?' }}
              </div>
              <span class="username">{{ row.username || '-' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="NAME">
          <template #default="{ row }">
            <span class="port-name">{{ row.name || '未命名' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="protocol" label="PROTOCOL" width="120">
          <template #default="{ row }">
            <span :class="['protocol-badge', row.protocol]">
              {{ protocolNames[row.protocol] || row.protocol }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="device_count" label="DEVICES" width="100">
          <template #default="{ row }">
            <div class="device-count">
              <el-icon><Iphone /></el-icon>
              <span>{{ row.device_count || 0 }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="STATUS" width="110">
          <template #default="{ row }">
            <span :class="['status-indicator', row.status]">
              <span class="status-dot"></span>
              {{ statusNames[row.status] }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="expire_at" label="EXPIRES" width="150">
          <template #default="{ row }">
            <div class="expire-cell">
              <span class="time-text">{{ formatTime(row.expire_at) }}</span>
              <span v-if="isExpiringSoon(row.expire_at)" class="expire-warning">
                <el-icon><Warning /></el-icon>
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="ACTIONS" width="140" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <button
                v-if="row.status === 'active'"
                class="btn-action stop"
                @click="handleStop(row)"
              >
                <el-icon><VideoPause /></el-icon>
                停止
              </button>
              <button
                v-else
                class="btn-action start"
                @click="handleStart(row)"
              >
                <el-icon><VideoPlay /></el-icon>
                启动
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
          @size-change="fetchPorts"
          @current-change="fetchPorts"
          class="cyber-pagination"
        />
      </div>

      <!-- Corner Decorations -->
      <div class="corner-decor top-left"></div>
      <div class="corner-decor top-right"></div>
      <div class="corner-decor bottom-left"></div>
      <div class="corner-decor bottom-right"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Iphone, Warning, VideoPause, VideoPlay } from '@element-plus/icons-vue'
import api from '../../api'
import dayjs from 'dayjs'

const loading = ref(false)
const ports = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const statusFilter = ref('')

const statusFilters = [
  { value: '', label: 'ALL', color: 'var(--text-muted)' },
  { value: 'active', label: 'ACTIVE', color: 'var(--neon-green)' },
  { value: 'expired', label: 'EXPIRED', color: '#ef4444' },
  { value: 'disabled', label: 'DISABLED', color: 'var(--text-muted)' }
]

const statusNames = {
  active: 'RUNNING',
  expired: 'EXPIRED',
  disabled: 'DISABLED'
}

const protocolNames = {
  autojs: 'Auto.js',
  lazyman: '懒人精灵',
  easyclick: 'EasyClick',
  websocket: 'WebSocket'
}

// Computed stats
const runningCount = computed(() => ports.value.filter(p => p.status === 'active').length)
const expiredCount = computed(() => ports.value.filter(p => p.status === 'expired').length)
const disabledCount = computed(() => ports.value.filter(p => p.status === 'disabled').length)

function formatTime(time) {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-'
}

function isExpiringSoon(time) {
  if (!time) return false
  const expireDate = dayjs(time)
  const now = dayjs()
  return expireDate.diff(now, 'day') <= 3 && expireDate.isAfter(now)
}

async function fetchPorts() {
  loading.value = true
  try {
    const res = await api.admin.getPorts(page.value, pageSize.value, statusFilter.value)
    if (res.code === 0) {
      ports.value = res.data.list || []
      total.value = res.data.total
    }
  } finally {
    loading.value = false
  }
}

async function handleStop(port) {
  const res = await api.admin.stopPort(port.id)
  if (res.code === 0) {
    ElMessage.success('端口已停止')
    fetchPorts()
  } else {
    ElMessage.error(res.message)
  }
}

async function handleStart(port) {
  const res = await api.admin.startPort(port.id)
  if (res.code === 0) {
    ElMessage.success('端口已启动')
    fetchPorts()
  } else {
    ElMessage.error(res.message)
  }
}

onMounted(() => {
  fetchPorts()
})
</script>

<style scoped>
.admin-ports-page {
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

/* Filter Group */
.filter-group {
  display: flex;
  gap: 8px;
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.filter-btn:hover {
  border-color: var(--neon-orange);
  color: var(--text-secondary);
}

.filter-btn.active {
  background: rgba(249, 115, 22, 0.1);
  border-color: var(--neon-orange);
  color: var(--neon-orange);
}

.filter-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
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
.stat-chip.expired .stat-value { color: #ef4444; }
.stat-chip.disabled .stat-value { color: var(--text-muted); }

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

.port-cell {
  display: flex;
  align-items: center;
}

.port-number {
  font-family: var(--font-mono);
  font-size: 16px;
  font-weight: 700;
  color: var(--neon-cyan);
  text-shadow: 0 0 10px rgba(0, 240, 255, 0.3);
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-avatar-small {
  width: 26px;
  height: 26px;
  background: linear-gradient(135deg, var(--neon-purple), var(--neon-purple-dim));
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-display);
  font-size: 12px;
  font-weight: 700;
  color: white;
}

.username {
  font-size: 13px;
  color: var(--text-secondary);
}

.port-name {
  color: var(--text-secondary);
}

.protocol-badge {
  display: inline-flex;
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.5px;
  background: rgba(168, 85, 247, 0.2);
  color: var(--neon-purple);
  border: 1px solid rgba(168, 85, 247, 0.3);
}

.protocol-badge.autojs {
  background: rgba(34, 197, 94, 0.2);
  color: var(--neon-green);
  border-color: rgba(34, 197, 94, 0.3);
}

.protocol-badge.lazyman {
  background: rgba(0, 240, 255, 0.2);
  color: var(--neon-cyan);
  border-color: rgba(0, 240, 255, 0.3);
}

.protocol-badge.easyclick {
  background: rgba(249, 115, 22, 0.2);
  color: var(--neon-orange);
  border-color: rgba(249, 115, 22, 0.3);
}

.device-count {
  display: flex;
  align-items: center;
  gap: 6px;
  font-family: var(--font-mono);
  font-size: 13px;
  color: var(--text-secondary);
}

.device-count .el-icon {
  color: var(--neon-cyan);
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
  animation: pulse 2s infinite;
}

.status-indicator.expired {
  color: #ef4444;
}

.status-indicator.expired .status-dot {
  background: #ef4444;
}

.status-indicator.disabled {
  color: var(--text-muted);
}

.status-indicator.disabled .status-dot {
  background: var(--text-muted);
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.expire-cell {
  display: flex;
  align-items: center;
  gap: 6px;
}

.time-text {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
}

.expire-warning {
  color: #eab308;
  animation: blink 1s infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
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
  padding: 6px 12px;
  border: 1px solid;
  border-radius: var(--radius-sm);
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  background: transparent;
}

.btn-action.stop {
  border-color: #ef4444;
  color: #ef4444;
}

.btn-action.stop:hover {
  background: rgba(239, 68, 68, 0.1);
}

.btn-action.start {
  border-color: var(--neon-green);
  color: var(--neon-green);
}

.btn-action.start:hover {
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
</style>
