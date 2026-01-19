<template>
  <div class="admin-devices-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <rect x="5" y="2" width="14" height="20" rx="2" ry="2"/>
            <path d="M12 18h.01"/>
          </svg>
        </div>
        <div class="header-text">
          <h1 class="page-title">
            <span class="title-prefix">//</span>
            DEVICE_MONITOR
          </h1>
          <p class="page-subtitle">系统设备状态监控</p>
        </div>
      </div>
      <div class="header-actions">
        <div class="filter-group">
          <button
            v-for="filter in statusFilters"
            :key="filter.value"
            :class="['filter-btn', { active: statusFilter === filter.value }]"
            @click="statusFilter = filter.value; fetchDevices()"
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
        <span class="stat-label">TOTAL_DEVICES</span>
        <span class="stat-value">{{ total }}</span>
      </div>
      <div class="stat-chip online">
        <span class="stat-label">ONLINE</span>
        <span class="stat-value">{{ onlineCount }}</span>
      </div>
      <div class="stat-chip offline">
        <span class="stat-label">OFFLINE</span>
        <span class="stat-value">{{ offlineCount }}</span>
      </div>
      <div class="stat-chip">
        <span class="stat-label">AVG_LATENCY</span>
        <span class="stat-value latency">{{ avgLatency }}ms</span>
      </div>
    </div>

    <!-- Devices Table -->
    <div class="table-container">
      <div class="table-header">
        <div class="table-title">
          <span class="title-dot"></span>
          DEVICE_REGISTRY
        </div>
        <div class="table-actions">
          <button class="action-btn refresh" @click="fetchDevices">
            <el-icon><Refresh /></el-icon>
            REFRESH
          </button>
        </div>
      </div>

      <el-table
        :data="devices"
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
        <el-table-column prop="device_name" label="DEVICE_NAME" width="180">
          <template #default="{ row }">
            <div class="device-cell">
              <div class="device-icon">
                <el-icon><Iphone /></el-icon>
              </div>
              <span class="device-name">{{ row.device_name || '未知设备' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="device_id" label="DEVICE_ID" width="200">
          <template #default="{ row }">
            <span class="device-id-text">{{ row.device_id }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="OWNER" width="120">
          <template #default="{ row }">
            <div class="user-cell">
              <div class="user-avatar-small">
                {{ row.username?.charAt(0).toUpperCase() || '?' }}
              </div>
              <span class="username">{{ row.username || '-' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="port_number" label="PORT" width="100">
          <template #default="{ row }">
            <span class="port-badge">:{{ row.port_number }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="last_ping" label="LATENCY" width="100">
          <template #default="{ row }">
            <span :class="['latency-value', getLatencyClass(row.last_ping)]">
              {{ row.last_ping || 0 }}ms
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="STATUS" width="110">
          <template #default="{ row }">
            <span :class="['status-indicator', row.status]">
              <span class="status-dot"></span>
              {{ row.status === 'online' ? 'ONLINE' : 'OFFLINE' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="connected_at" label="CONNECTED" width="160">
          <template #default="{ row }">
            <span class="time-text">{{ formatTime(row.connected_at) }}</span>
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
          @size-change="fetchDevices"
          @current-change="fetchDevices"
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
import { Refresh, Iphone } from '@element-plus/icons-vue'
import api from '../../api'
import dayjs from 'dayjs'

const loading = ref(false)
const devices = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const statusFilter = ref('')

const statusFilters = [
  { value: '', label: 'ALL', color: 'var(--text-muted)' },
  { value: 'online', label: 'ONLINE', color: 'var(--neon-green)' },
  { value: 'offline', label: 'OFFLINE', color: 'var(--text-muted)' }
]

// Computed stats
const onlineCount = computed(() => devices.value.filter(d => d.status === 'online').length)
const offlineCount = computed(() => devices.value.filter(d => d.status === 'offline').length)
const avgLatency = computed(() => {
  const onlineDevices = devices.value.filter(d => d.status === 'online' && d.last_ping)
  if (onlineDevices.length === 0) return 0
  const sum = onlineDevices.reduce((acc, d) => acc + (d.last_ping || 0), 0)
  return Math.round(sum / onlineDevices.length)
})

function formatTime(time) {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-'
}

function getLatencyClass(ping) {
  if (!ping) return 'none'
  if (ping < 50) return 'good'
  if (ping < 150) return 'medium'
  return 'high'
}

async function fetchDevices() {
  loading.value = true
  try {
    const res = await api.admin.getDevices(page.value, pageSize.value, statusFilter.value)
    if (res.code === 0) {
      devices.value = res.data.list || []
      total.value = res.data.total
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDevices()
})
</script>

<style scoped>
.admin-devices-page {
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

.stat-chip.online .stat-value { color: var(--neon-green); }
.stat-chip.offline .stat-value { color: var(--text-muted); }
.stat-value.latency { color: var(--neon-cyan); }

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

.device-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.device-icon {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, var(--neon-cyan), rgba(0, 240, 255, 0.5));
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.device-name {
  font-weight: 500;
  color: var(--text-primary);
}

.device-id-text {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
  background: var(--bg-elevated);
  padding: 4px 8px;
  border-radius: var(--radius-sm);
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

.port-badge {
  font-family: var(--font-mono);
  font-size: 13px;
  font-weight: 600;
  color: var(--neon-cyan);
}

.latency-value {
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 600;
}

.latency-value.none {
  color: var(--text-muted);
}

.latency-value.good {
  color: var(--neon-green);
}

.latency-value.medium {
  color: #eab308;
}

.latency-value.high {
  color: #ef4444;
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

.status-indicator.online {
  color: var(--neon-green);
}

.status-indicator.online .status-dot {
  background: var(--neon-green);
  box-shadow: 0 0 8px var(--neon-green);
  animation: pulse 2s infinite;
}

.status-indicator.offline {
  color: var(--text-muted);
}

.status-indicator.offline .status-dot {
  background: var(--text-muted);
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.time-text {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
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
