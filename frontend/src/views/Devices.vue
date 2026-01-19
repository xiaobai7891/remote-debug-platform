<template>
  <div class="devices-page">
    <!-- Page header -->
    <div class="page-header">
      <div class="header-title">
        <div class="title-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="5" y="2" width="14" height="20" rx="2"/>
            <path d="M12 18h.01"/>
          </svg>
        </div>
        <div class="title-content">
          <h1>DEVICE_MONITOR</h1>
          <p class="subtitle">// 实时设备监控与远程控制</p>
        </div>
      </div>
      <div class="header-actions">
        <div class="live-indicator">
          <span class="live-dot"></span>
          <span class="live-text">LIVE</span>
        </div>
        <button class="refresh-btn" @click="fetchDevices" :disabled="loading">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: loading }">
            <path d="M23 4v6h-6M1 20v-6h6"/>
            <path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15"/>
          </svg>
          <span>刷新</span>
        </button>
      </div>
    </div>

    <!-- Stats overview -->
    <div class="stats-bar">
      <div class="stat-item">
        <span class="stat-label">TOTAL_DEVICES</span>
        <span class="stat-value">{{ devices.length }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">ONLINE</span>
        <span class="stat-value online">{{ onlineCount }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">OFFLINE</span>
        <span class="stat-value offline">{{ offlineCount }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">AVG_LATENCY</span>
        <span class="stat-value latency">{{ avgLatency }}ms</span>
      </div>
    </div>

    <!-- Device grid -->
    <div class="devices-container" v-loading="loading" element-loading-background="rgba(5, 5, 10, 0.8)">
      <!-- Empty state -->
      <div v-if="!loading && devices.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <rect x="5" y="2" width="14" height="20" rx="2"/>
            <path d="M12 18h.01"/>
            <path d="M9 6h6M9 9h6M9 12h3"/>
          </svg>
        </div>
        <p class="empty-text">暂无连接设备</p>
        <p class="empty-hint">等待设备通过端口连接...</p>
      </div>

      <!-- Device cards -->
      <div v-else class="devices-grid">
        <div
          v-for="device in devices"
          :key="device.id"
          class="device-card"
          :class="{ online: device.status === 'online', offline: device.status !== 'online' }"
        >
          <!-- Device header -->
          <div class="device-header">
            <div class="device-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="5" y="2" width="14" height="20" rx="2"/>
                <path d="M12 18h.01"/>
              </svg>
            </div>
            <div class="device-status">
              <span class="status-dot"></span>
              <span class="status-text">{{ device.status === 'online' ? 'ONLINE' : 'OFFLINE' }}</span>
            </div>
          </div>

          <!-- Device name -->
          <div class="device-name">
            <span class="name-label">DEVICE_NAME</span>
            <span class="name-value">{{ device.device_name || 'Unknown Device' }}</span>
          </div>

          <!-- Device info grid -->
          <div class="device-info">
            <div class="info-item">
              <span class="info-label">ID</span>
              <span class="info-value mono">{{ device.device_id?.substring(0, 12) }}...</span>
            </div>
            <div class="info-item">
              <span class="info-label">PORT</span>
              <span class="info-value highlight">{{ device.port }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">LATENCY</span>
              <span class="info-value" :class="getPingClass(device.last_ping)">
                {{ device.last_ping || 0 }}ms
              </span>
            </div>
            <div class="info-item">
              <span class="info-label">CONNECTED</span>
              <span class="info-value">{{ formatRelativeTime(device.connected_at) }}</span>
            </div>
          </div>

          <!-- Latency bar -->
          <div class="latency-bar">
            <div
              class="latency-fill"
              :class="getPingClass(device.last_ping)"
              :style="{ width: getLatencyWidth(device.last_ping) }"
            ></div>
          </div>

          <!-- Actions -->
          <div class="device-actions">
            <button class="action-btn info" @click="showDetailDialog(device)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <path d="M12 16v-4m0-4h.01"/>
              </svg>
              <span>详情</span>
            </button>
            <button
              class="action-btn exec"
              @click="showExecDialog(device)"
              :disabled="device.status !== 'online'"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="4 17 10 11 4 5"/>
                <line x1="12" y1="19" x2="20" y2="19"/>
              </svg>
              <span>执行代码</span>
            </button>
          </div>

          <!-- Corner decorations -->
          <div class="corner tl"></div>
          <div class="corner tr"></div>
          <div class="corner bl"></div>
          <div class="corner br"></div>

          <!-- Pulse animation for online devices -->
          <div v-if="device.status === 'online'" class="pulse-ring"></div>
        </div>
      </div>
    </div>

    <!-- Device Detail Dialog -->
    <el-dialog
      v-model="detailDialogVisible"
      title=""
      width="750"
      :show-close="false"
      class="cyber-dialog detail-dialog"
    >
      <template #header>
        <div class="dialog-header">
          <div class="dialog-title">
            <span class="title-bracket">[</span>
            DEVICE_{{ selectedDevice?.device_id?.substring(0, 8) }}_INFO
            <span class="title-bracket">]</span>
          </div>
          <div class="dialog-actions">
            <div class="device-status-badge" :class="selectedDevice?.status">
              <span class="status-dot"></span>
              {{ selectedDevice?.status === 'online' ? 'ONLINE' : 'OFFLINE' }}
            </div>
            <button class="close-btn" @click="detailDialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12"/>
              </svg>
            </button>
          </div>
        </div>
      </template>

      <div class="dialog-content">
        <!-- Tab navigation -->
        <div class="tab-nav">
          <button
            class="tab-btn"
            :class="{ active: activeTab === 'info' }"
            @click="activeTab = 'info'"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <path d="M12 16v-4m0-4h.01"/>
            </svg>
            基本信息
          </button>
          <button
            class="tab-btn"
            :class="{ active: activeTab === 'logs' }"
            @click="activeTab = 'logs'"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
              <path d="M14 2v6h6M16 13H8M16 17H8M10 9H8"/>
            </svg>
            设备日志
            <span v-if="deviceLogs.length" class="log-count">{{ deviceLogs.length }}</span>
          </button>
        </div>

        <!-- Info panel -->
        <div v-if="activeTab === 'info'" class="info-panel">
          <div class="info-grid">
            <div class="info-card">
              <span class="card-label">DEVICE_NAME</span>
              <span class="card-value">{{ selectedDevice?.device_name || '-' }}</span>
            </div>
            <div class="info-card">
              <span class="card-label">DEVICE_ID</span>
              <span class="card-value mono">{{ selectedDevice?.device_id }}</span>
            </div>
            <div class="info-card">
              <span class="card-label">PORT</span>
              <span class="card-value highlight">{{ selectedDevice?.port }}</span>
            </div>
            <div class="info-card">
              <span class="card-label">STATUS</span>
              <span class="card-value">
                <span class="status-indicator" :class="selectedDevice?.status">
                  {{ selectedDevice?.status === 'online' ? '在线' : '离线' }}
                </span>
              </span>
            </div>
            <div class="info-card">
              <span class="card-label">LATENCY</span>
              <span class="card-value" :class="getPingClass(selectedDevice?.last_ping)">
                {{ selectedDevice?.last_ping || 0 }}ms
              </span>
            </div>
            <div class="info-card">
              <span class="card-label">CONNECTED_AT</span>
              <span class="card-value">{{ formatTime(selectedDevice?.connected_at) }}</span>
            </div>
          </div>

          <!-- Device info JSON -->
          <div class="device-info-block">
            <div class="block-header">
              <span class="block-title">
                <span class="title-icon">#</span>
                DEVICE_INFO_JSON
              </span>
              <button class="copy-btn" @click="copyDeviceInfo">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="9" y="9" width="13" height="13" rx="2"/>
                  <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                </svg>
                复制
              </button>
            </div>
            <pre class="json-content"><code>{{ formatDeviceInfo(selectedDevice?.device_info) }}</code></pre>
          </div>
        </div>

        <!-- Logs panel -->
        <div v-else class="logs-panel">
          <div class="logs-header">
            <div class="logs-filter">
              <button
                v-for="level in ['all', 'info', 'warn', 'error']"
                :key="level"
                class="filter-btn"
                :class="{ active: logFilter === level }"
                @click="logFilter = level"
              >
                {{ level.toUpperCase() }}
              </button>
            </div>
            <button class="refresh-logs-btn" @click="fetchDeviceLogs">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M23 4v6h-6M1 20v-6h6"/>
                <path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15"/>
              </svg>
            </button>
          </div>

          <div class="logs-container">
            <div v-if="filteredLogs.length === 0" class="logs-empty">
              <span>暂无日志记录</span>
            </div>
            <div
              v-for="log in filteredLogs"
              :key="log.id"
              class="log-entry"
              :class="'level-' + log.level"
            >
              <span class="log-time">{{ formatTime(log.created_at) }}</span>
              <span class="log-level">[{{ log.level.toUpperCase() }}]</span>
              <span class="log-message">{{ log.message }}</span>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- Execute Code Dialog -->
    <el-dialog
      v-model="execDialogVisible"
      title=""
      width="650"
      :show-close="false"
      class="cyber-dialog exec-dialog"
    >
      <template #header>
        <div class="dialog-header">
          <div class="dialog-title">
            <span class="title-bracket">[</span>
            REMOTE_EXECUTE
            <span class="title-bracket">]</span>
          </div>
          <button class="close-btn" @click="execDialogVisible = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6L6 18M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </template>

      <div class="dialog-content">
        <!-- Target device info -->
        <div class="target-device">
          <span class="target-label">TARGET_DEVICE</span>
          <div class="target-info">
            <span class="target-name">{{ selectedDevice?.device_name || 'Unknown' }}</span>
            <span class="target-id">{{ selectedDevice?.device_id?.substring(0, 12) }}...</span>
            <span class="target-port">PORT: {{ selectedDevice?.port }}</span>
          </div>
        </div>

        <!-- Code editor -->
        <div class="code-editor">
          <div class="editor-header">
            <span class="editor-title">
              <span class="title-icon">></span>
              CODE_INPUT
            </span>
            <div class="editor-actions">
              <button class="editor-btn" @click="execCode = 'toast(&quot;Hello from remote!&quot;);'" title="插入示例">
                示例
              </button>
              <button class="editor-btn" @click="execCode = ''" title="清空">
                清空
              </button>
            </div>
          </div>
          <div class="editor-body">
            <div class="line-numbers">
              <span v-for="n in lineCount" :key="n">{{ n }}</span>
            </div>
            <textarea
              v-model="execCode"
              class="code-input"
              placeholder="// 输入要执行的 JavaScript 代码"
              spellcheck="false"
              @keydown.tab.prevent="insertTab"
            ></textarea>
          </div>
        </div>

        <!-- Quick actions -->
        <div class="quick-actions">
          <span class="quick-label">QUICK_COMMANDS</span>
          <div class="quick-btns">
            <button class="quick-btn" @click="execCode = 'toast(&quot;Hello!&quot;);'">Toast</button>
            <button class="quick-btn" @click="execCode = 'console.log(device.getAndroidId());'">Device ID</button>
            <button class="quick-btn" @click="execCode = 'home();'">Home</button>
            <button class="quick-btn" @click="execCode = 'back();'">Back</button>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <button class="cancel-btn" @click="execDialogVisible = false">取消</button>
          <button class="confirm-btn" @click="handleExec" :disabled="execLoading || !execCode.trim()">
            <span v-if="execLoading" class="loading-spinner"></span>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="5 3 19 12 5 21 5 3"/>
            </svg>
            <span>{{ execLoading ? '发送中...' : '执行代码' }}</span>
          </button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '../api'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const loading = ref(false)
const devices = ref([])

// Computed stats
const onlineCount = computed(() => devices.value.filter(d => d.status === 'online').length)
const offlineCount = computed(() => devices.value.filter(d => d.status !== 'online').length)
const avgLatency = computed(() => {
  const onlineDevices = devices.value.filter(d => d.status === 'online' && d.last_ping)
  if (onlineDevices.length === 0) return 0
  const total = onlineDevices.reduce((sum, d) => sum + (d.last_ping || 0), 0)
  return Math.round(total / onlineDevices.length)
})

// Detail dialog
const detailDialogVisible = ref(false)
const selectedDevice = ref(null)
const deviceLogs = ref([])
const activeTab = ref('info')
const logFilter = ref('all')

const filteredLogs = computed(() => {
  if (logFilter.value === 'all') return deviceLogs.value
  return deviceLogs.value.filter(log => log.level === logFilter.value)
})

// Exec dialog
const execDialogVisible = ref(false)
const execCode = ref('')
const execLoading = ref(false)

const lineCount = computed(() => {
  return Math.max(execCode.value.split('\n').length, 10)
})

function formatTime(time) {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm:ss') : '-'
}

function formatRelativeTime(time) {
  return time ? dayjs(time).fromNow() : '-'
}

function getPingClass(ping) {
  if (!ping) return ''
  if (ping < 100) return 'ping-good'
  if (ping < 300) return 'ping-normal'
  return 'ping-bad'
}

function getLatencyWidth(ping) {
  if (!ping) return '0%'
  const maxPing = 500
  const percentage = Math.min((ping / maxPing) * 100, 100)
  return percentage + '%'
}

function formatDeviceInfo(info) {
  if (!info) return '{\n  "info": "No device info available"\n}'
  try {
    return JSON.stringify(JSON.parse(info), null, 2)
  } catch {
    return info
  }
}

function copyDeviceInfo() {
  const info = formatDeviceInfo(selectedDevice.value?.device_info)
  navigator.clipboard.writeText(info)
  ElMessage.success('设备信息已复制')
}

function insertTab(e) {
  const start = e.target.selectionStart
  const end = e.target.selectionEnd
  execCode.value = execCode.value.substring(0, start) + '    ' + execCode.value.substring(end)
  setTimeout(() => {
    e.target.selectionStart = e.target.selectionEnd = start + 4
  }, 0)
}

async function fetchDevices() {
  loading.value = true
  try {
    const res = await api.devices.list()
    if (res.code === 0) {
      devices.value = res.data || []
    }
  } finally {
    loading.value = false
  }
}

async function fetchDeviceLogs() {
  if (!selectedDevice.value) return
  const res = await api.devices.getLogs(selectedDevice.value.id)
  if (res.code === 0) {
    deviceLogs.value = res.data || []
  }
}

async function showDetailDialog(device) {
  selectedDevice.value = device
  activeTab.value = 'info'
  logFilter.value = 'all'
  detailDialogVisible.value = true
  await fetchDeviceLogs()
}

function showExecDialog(device) {
  selectedDevice.value = device
  execCode.value = 'toast("Hello from remote!");'
  execDialogVisible.value = true
}

async function handleExec() {
  if (!execCode.value.trim()) {
    ElMessage.warning('请输入要执行的代码')
    return
  }
  execLoading.value = true
  try {
    const res = await api.devices.exec(selectedDevice.value.id, execCode.value)
    if (res.code === 0) {
      ElMessage.success('执行命令已发送')
      execDialogVisible.value = false
    } else {
      ElMessage.error(res.message)
    }
  } finally {
    execLoading.value = false
  }
}

// Auto refresh
let refreshInterval

onMounted(() => {
  fetchDevices()
  refreshInterval = setInterval(fetchDevices, 30000) // Refresh every 30s
})

onUnmounted(() => {
  clearInterval(refreshInterval)
})
</script>

<style scoped>
.devices-page {
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

.header-title {
  display: flex;
  align-items: center;
  gap: 16px;
}

.title-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(168, 85, 247, 0.1), rgba(0, 240, 255, 0.1));
  border: 1px solid var(--neon-purple);
  border-radius: var(--radius-lg);
  color: var(--neon-purple);
}

.title-icon svg {
  width: 24px;
  height: 24px;
}

.title-content h1 {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: 2px;
  margin: 0 0 4px 0;
}

.subtitle {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
  margin: 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.live-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.3);
  border-radius: var(--radius-full);
}

.live-dot {
  width: 8px;
  height: 8px;
  background: var(--neon-green);
  border-radius: 50%;
  box-shadow: 0 0 10px var(--neon-green);
  animation: pulse-live 1.5s ease-in-out infinite;
}

@keyframes pulse-live {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(0.8); }
}

.live-text {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  color: var(--neon-green);
  letter-spacing: 1px;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.refresh-btn:hover:not(:disabled) {
  border-color: var(--neon-cyan);
  color: var(--neon-cyan);
}

.refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.refresh-btn svg {
  width: 16px;
  height: 16px;
}

.refresh-btn svg.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Stats Bar */
.stats-bar {
  display: flex;
  gap: 24px;
  margin-bottom: 24px;
  padding: 16px 24px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-label {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.stat-value {
  font-family: var(--font-mono);
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-value.online {
  color: var(--neon-green);
}

.stat-value.offline {
  color: var(--text-muted);
}

.stat-value.latency {
  color: var(--neon-cyan);
}

/* Empty State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
}

.empty-icon {
  width: 80px;
  height: 80px;
  color: var(--text-muted);
  opacity: 0.5;
  margin-bottom: 20px;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.empty-icon svg {
  width: 100%;
  height: 100%;
}

.empty-text {
  font-size: 18px;
  color: var(--text-secondary);
  margin: 0 0 8px 0;
}

.empty-hint {
  font-family: var(--font-mono);
  font-size: 13px;
  color: var(--text-muted);
  margin: 0;
}

/* Devices Grid */
.devices-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

/* Device Card */
.device-card {
  position: relative;
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  padding: 20px;
  transition: all 0.3s ease;
  overflow: hidden;
}

.device-card:hover {
  border-color: var(--neon-purple-dim);
  transform: translateY(-2px);
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
}

.device-card.online {
  border-left: 3px solid var(--neon-green);
}

.device-card.offline {
  border-left: 3px solid var(--text-muted);
  opacity: 0.7;
}

.device-card .corner {
  position: absolute;
  width: 10px;
  height: 10px;
  border-color: var(--border-default);
  opacity: 0.5;
  transition: all 0.3s ease;
}

.device-card:hover .corner {
  border-color: var(--neon-purple);
  opacity: 1;
}

.corner.tl { top: 8px; left: 8px; border-top: 2px solid; border-left: 2px solid; }
.corner.tr { top: 8px; right: 8px; border-top: 2px solid; border-right: 2px solid; }
.corner.bl { bottom: 8px; left: 8px; border-bottom: 2px solid; border-left: 2px solid; }
.corner.br { bottom: 8px; right: 8px; border-bottom: 2px solid; border-right: 2px solid; }

.pulse-ring {
  position: absolute;
  top: 20px;
  right: 20px;
  width: 50px;
  height: 50px;
  border: 2px solid var(--neon-green);
  border-radius: 50%;
  opacity: 0;
  animation: pulse-ring 2s ease-out infinite;
}

@keyframes pulse-ring {
  0% { transform: scale(0.5); opacity: 0.8; }
  100% { transform: scale(1.5); opacity: 0; }
}

/* Device Header */
.device-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.device-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-elevated);
  border-radius: var(--radius-md);
  color: var(--neon-purple);
}

.device-icon svg {
  width: 22px;
  height: 22px;
}

.device-status {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  background: var(--bg-elevated);
  border-radius: var(--radius-full);
}

.device-card.online .device-status {
  background: rgba(34, 197, 94, 0.1);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--text-muted);
}

.device-card.online .status-dot {
  background: var(--neon-green);
  box-shadow: 0 0 8px var(--neon-green);
}

.status-text {
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 1px;
  color: var(--text-muted);
}

.device-card.online .status-text {
  color: var(--neon-green);
}

/* Device Name */
.device-name {
  margin-bottom: 16px;
}

.name-label {
  display: block;
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 1px;
  margin-bottom: 4px;
}

.name-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

/* Device Info */
.device-info {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border-subtle);
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-family: var(--font-mono);
  font-size: 9px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.info-value {
  font-family: var(--font-mono);
  font-size: 13px;
  color: var(--text-secondary);
}

.info-value.mono {
  font-size: 11px;
}

.info-value.highlight {
  color: var(--neon-cyan);
  font-weight: 600;
}

.info-value.ping-good {
  color: var(--neon-green);
}

.info-value.ping-normal {
  color: var(--status-warning);
}

.info-value.ping-bad {
  color: var(--status-error);
}

/* Latency Bar */
.latency-bar {
  height: 4px;
  background: var(--bg-elevated);
  border-radius: var(--radius-full);
  margin-bottom: 16px;
  overflow: hidden;
}

.latency-fill {
  height: 100%;
  border-radius: var(--radius-full);
  transition: width 0.5s ease;
}

.latency-fill.ping-good {
  background: var(--neon-green);
  box-shadow: 0 0 10px var(--neon-green);
}

.latency-fill.ping-normal {
  background: var(--status-warning);
  box-shadow: 0 0 10px var(--status-warning);
}

.latency-fill.ping-bad {
  background: var(--status-error);
  box-shadow: 0 0 10px var(--status-error);
}

/* Device Actions */
.device-actions {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn svg {
  width: 14px;
  height: 14px;
}

.action-btn.info:hover {
  border-color: var(--neon-cyan);
  color: var(--neon-cyan);
  background: rgba(0, 240, 255, 0.05);
}

.action-btn.exec:hover:not(:disabled) {
  border-color: var(--neon-purple);
  color: var(--neon-purple);
  background: rgba(168, 85, 247, 0.05);
}

.action-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* Dialog Styles */
.cyber-dialog :deep(.el-dialog) {
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-xl);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

.cyber-dialog :deep(.el-dialog__header),
.cyber-dialog :deep(.el-dialog__body),
.cyber-dialog :deep(.el-dialog__footer) {
  padding: 0;
  margin: 0;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-subtle);
}

.dialog-title {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: 1px;
}

.title-bracket {
  color: var(--neon-purple);
}

.dialog-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.device-status-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: var(--radius-full);
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 1px;
}

.device-status-badge.online {
  background: rgba(34, 197, 94, 0.15);
  color: var(--neon-green);
}

.device-status-badge.online .status-dot {
  background: var(--neon-green);
  box-shadow: 0 0 6px var(--neon-green);
}

.device-status-badge.offline {
  background: rgba(148, 163, 184, 0.15);
  color: var(--text-muted);
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: transparent;
  border: 1px solid var(--border-default);
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.2s ease;
}

.close-btn:hover {
  border-color: var(--status-error);
  color: var(--status-error);
}

.close-btn svg {
  width: 16px;
  height: 16px;
}

.dialog-content {
  padding: 24px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid var(--border-subtle);
}

.cancel-btn {
  padding: 10px 20px;
  background: transparent;
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.cancel-btn:hover {
  border-color: var(--text-secondary);
  color: var(--text-primary);
}

.confirm-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 24px;
  background: linear-gradient(135deg, var(--neon-purple), var(--neon-purple-dim));
  border: none;
  border-radius: var(--radius-md);
  color: white;
  font-family: var(--font-mono);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.confirm-btn:hover:not(:disabled) {
  box-shadow: 0 4px 20px rgba(168, 85, 247, 0.3);
}

.confirm-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.confirm-btn svg {
  width: 14px;
  height: 14px;
}

.loading-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid transparent;
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

/* Tab Navigation */
.tab-nav {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border-subtle);
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: transparent;
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab-btn svg {
  width: 16px;
  height: 16px;
}

.tab-btn:hover {
  border-color: var(--neon-purple-dim);
  color: var(--text-secondary);
}

.tab-btn.active {
  border-color: var(--neon-purple);
  background: rgba(168, 85, 247, 0.1);
  color: var(--neon-purple);
}

.log-count {
  padding: 2px 8px;
  background: var(--neon-purple);
  border-radius: var(--radius-full);
  font-size: 10px;
  color: white;
}

/* Info Panel */
.info-panel .info-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 24px;
}

.info-card {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 14px;
  background: var(--bg-elevated);
  border-radius: var(--radius-md);
}

.card-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.card-value {
  font-family: var(--font-mono);
  font-size: 14px;
  color: var(--text-primary);
}

.card-value.mono {
  font-size: 11px;
  word-break: break-all;
}

.card-value.highlight {
  color: var(--neon-cyan);
  font-weight: 600;
}

.status-indicator {
  display: inline-flex;
  padding: 3px 10px;
  border-radius: var(--radius-full);
  font-size: 11px;
}

.status-indicator.online {
  background: rgba(34, 197, 94, 0.15);
  color: var(--neon-green);
}

.status-indicator.offline {
  background: rgba(148, 163, 184, 0.15);
  color: var(--text-muted);
}

/* Device Info Block */
.device-info-block {
  margin-top: 20px;
}

.block-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.block-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.block-title .title-icon {
  color: var(--neon-purple);
}

.copy-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: transparent;
  border: 1px solid var(--border-default);
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.copy-btn svg {
  width: 12px;
  height: 12px;
}

.copy-btn:hover {
  border-color: var(--neon-cyan);
  color: var(--neon-cyan);
}

.json-content {
  padding: 16px;
  background: var(--bg-void);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  margin: 0;
  overflow-x: auto;
  max-height: 200px;
}

.json-content code {
  font-family: var(--font-mono);
  font-size: 12px;
  line-height: 1.6;
  color: var(--text-secondary);
  white-space: pre;
}

/* Logs Panel */
.logs-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.logs-filter {
  display: flex;
  gap: 4px;
}

.filter-btn {
  padding: 6px 12px;
  background: transparent;
  border: 1px solid var(--border-default);
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.filter-btn:hover {
  border-color: var(--neon-purple-dim);
  color: var(--text-secondary);
}

.filter-btn.active {
  border-color: var(--neon-purple);
  background: rgba(168, 85, 247, 0.1);
  color: var(--neon-purple);
}

.refresh-logs-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: transparent;
  border: 1px solid var(--border-default);
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.2s ease;
}

.refresh-logs-btn:hover {
  border-color: var(--neon-cyan);
  color: var(--neon-cyan);
}

.refresh-logs-btn svg {
  width: 14px;
  height: 14px;
}

.logs-container {
  max-height: 400px;
  overflow-y: auto;
  background: var(--bg-void);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
}

.logs-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px;
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 13px;
}

.log-entry {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 10px 16px;
  border-bottom: 1px solid var(--border-subtle);
  font-family: var(--font-mono);
  font-size: 12px;
}

.log-entry:last-child {
  border-bottom: none;
}

.log-time {
  color: var(--text-muted);
  white-space: nowrap;
}

.log-level {
  font-weight: 600;
  white-space: nowrap;
}

.log-entry.level-info .log-level {
  color: var(--neon-cyan);
}

.log-entry.level-warn .log-level {
  color: var(--status-warning);
}

.log-entry.level-error .log-level {
  color: var(--status-error);
}

.log-message {
  flex: 1;
  color: var(--text-secondary);
  word-break: break-word;
}

/* Exec Dialog */
.target-device {
  padding: 16px;
  background: var(--bg-elevated);
  border-radius: var(--radius-md);
  margin-bottom: 20px;
}

.target-label {
  display: block;
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 1px;
  margin-bottom: 10px;
}

.target-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.target-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.target-id {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
}

.target-port {
  padding: 4px 10px;
  background: rgba(0, 240, 255, 0.1);
  border-radius: var(--radius-sm);
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--neon-cyan);
}

/* Code Editor */
.code-editor {
  margin-bottom: 20px;
}

.editor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.editor-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.editor-title .title-icon {
  color: var(--neon-purple);
}

.editor-actions {
  display: flex;
  gap: 8px;
}

.editor-btn {
  padding: 4px 10px;
  background: transparent;
  border: 1px solid var(--border-default);
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.editor-btn:hover {
  border-color: var(--neon-purple-dim);
  color: var(--text-secondary);
}

.editor-body {
  display: flex;
  background: var(--bg-void);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.line-numbers {
  display: flex;
  flex-direction: column;
  padding: 16px 12px;
  background: var(--bg-surface);
  border-right: 1px solid var(--border-subtle);
  font-family: var(--font-mono);
  font-size: 12px;
  line-height: 1.6;
  color: var(--text-muted);
  text-align: right;
  user-select: none;
}

.code-input {
  flex: 1;
  min-height: 250px;
  padding: 16px;
  background: transparent;
  border: none;
  outline: none;
  resize: vertical;
  font-family: var(--font-mono);
  font-size: 12px;
  line-height: 1.6;
  color: var(--text-primary);
}

.code-input::placeholder {
  color: var(--text-muted);
}

/* Quick Actions */
.quick-actions {
  padding: 16px;
  background: var(--bg-elevated);
  border-radius: var(--radius-md);
}

.quick-label {
  display: block;
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 1px;
  margin-bottom: 12px;
}

.quick-btns {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.quick-btn {
  padding: 8px 14px;
  background: var(--bg-surface);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.quick-btn:hover {
  border-color: var(--neon-purple);
  color: var(--neon-purple);
}

/* Responsive */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .header-actions {
    justify-content: space-between;
  }

  .stats-bar {
    flex-wrap: wrap;
  }

  .devices-grid {
    grid-template-columns: 1fr;
  }

  .info-panel .info-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .target-info {
    flex-wrap: wrap;
  }
}
</style>
