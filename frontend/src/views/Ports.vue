<template>
  <div class="ports-page">
    <!-- Page header -->
    <div class="page-header">
      <div class="header-title">
        <div class="title-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
          </svg>
        </div>
        <div class="title-content">
          <h1>PORT_MANAGER</h1>
          <p class="subtitle">// 端口分配与管理控制台</p>
        </div>
      </div>
      <button class="create-btn" @click="showCreateDialog">
        <span class="btn-icon">+</span>
        <span class="btn-text">申请新端口</span>
        <span class="btn-glow"></span>
      </button>
    </div>

    <!-- Stats bar -->
    <div class="stats-bar">
      <div class="stat-item">
        <span class="stat-label">TOTAL_PORTS</span>
        <span class="stat-value">{{ ports.length }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">ACTIVE</span>
        <span class="stat-value active">{{ activePorts }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">EXPIRED</span>
        <span class="stat-value expired">{{ expiredPorts }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">TOTAL_DEVICES</span>
        <span class="stat-value">{{ totalDevices }}</span>
      </div>
    </div>

    <!-- Ports grid -->
    <div class="ports-container" v-loading="loading" element-loading-background="rgba(5, 5, 10, 0.8)">
      <div v-if="ports.length === 0 && !loading" class="empty-state">
        <div class="empty-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M9 3H5a2 2 0 00-2 2v4m6-6h10a2 2 0 012 2v4M9 3v18m0 0h10a2 2 0 002-2v-4M9 21H5a2 2 0 01-2-2v-4m0-6v6"/>
          </svg>
        </div>
        <p class="empty-text">暂无端口</p>
        <p class="empty-hint">点击"申请新端口"开始使用</p>
      </div>

      <div v-else class="ports-grid">
        <div
          v-for="port in ports"
          :key="port.id"
          class="port-card"
          :class="{
            'status-active': port.status === 'active',
            'status-expired': port.status === 'expired',
            'status-disabled': port.status === 'disabled',
            'expiring-soon': isExpiringSoon(port.expire_at)
          }"
        >
          <!-- Port header -->
          <div class="port-header">
            <div class="port-number">
              <span class="port-label">PORT</span>
              <span class="port-value">{{ port.port }}</span>
            </div>
            <div class="port-status">
              <span class="status-dot"></span>
              <span class="status-text">{{ statusNames[port.status] }}</span>
            </div>
          </div>

          <!-- Port info -->
          <div class="port-info">
            <div class="info-row">
              <span class="info-label">NAME</span>
              <span class="info-value">{{ port.name || '未命名' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">PROTOCOL</span>
              <span class="protocol-badge" :class="port.protocol">
                {{ protocolNames[port.protocol] }}
              </span>
            </div>
            <div class="info-row">
              <span class="info-label">DEVICES</span>
              <span class="info-value devices">
                <span class="device-count">{{ port.device_count || 0 }}</span>
                <span class="device-label">online</span>
              </span>
            </div>
            <div class="info-row">
              <span class="info-label">EXPIRE</span>
              <span class="info-value expire" :class="{ warning: isExpiringSoon(port.expire_at) }">
                {{ formatTime(port.expire_at) }}
              </span>
            </div>
          </div>

          <!-- Token section -->
          <div class="token-section">
            <div class="token-label">
              <span class="label-icon">#</span>
              CONNECTION_TOKEN
            </div>
            <div class="token-value">
              <code>{{ port.token.substring(0, 20) }}...</code>
              <button class="copy-btn" @click="copyToken(port.token)" title="复制密钥">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="9" y="9" width="13" height="13" rx="2"/>
                  <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                </svg>
              </button>
            </div>
          </div>

          <!-- Actions -->
          <div class="port-actions">
            <button class="action-btn detail" @click="showDetailDialog(port)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <path d="M12 16v-4m0-4h.01"/>
              </svg>
              <span>详情</span>
            </button>
            <button class="action-btn renew" @click="showRenewDialog(port)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M23 4v6h-6M1 20v-6h6"/>
                <path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15"/>
              </svg>
              <span>续费</span>
            </button>
            <button class="action-btn change" @click="handleChangePort(port)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M16 3h5v5M8 3H3v5M3 16v5h5M21 16v5h-5"/>
              </svg>
              <span>换端口</span>
            </button>
            <button class="action-btn delete" @click="confirmDelete(port)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
              </svg>
              <span>释放</span>
            </button>
          </div>

          <!-- Decorative corners -->
          <div class="corner tl"></div>
          <div class="corner tr"></div>
          <div class="corner bl"></div>
          <div class="corner br"></div>
        </div>
      </div>
    </div>

    <!-- Create Port Dialog -->
    <el-dialog
      v-model="createDialogVisible"
      title=""
      width="520"
      :show-close="false"
      class="cyber-dialog"
    >
      <template #header>
        <div class="dialog-header">
          <div class="dialog-title">
            <span class="title-bracket">[</span>
            NEW_PORT_REQUEST
            <span class="title-bracket">]</span>
          </div>
          <button class="close-btn" @click="createDialogVisible = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6L6 18M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </template>

      <div class="dialog-content">
        <div class="form-group">
          <label class="form-label">
            <span class="label-icon">></span>
            PORT_NAME
          </label>
          <el-input
            v-model="createForm.name"
            placeholder="可选，用于备注识别"
            class="cyber-input"
          />
        </div>

        <div class="form-group">
          <label class="form-label">
            <span class="label-icon">></span>
            PROTOCOL_TYPE
          </label>
          <div class="protocol-select">
            <div
              v-for="proto in protocols"
              :key="proto.value"
              class="protocol-option"
              :class="{ active: createForm.protocol === proto.value }"
              @click="createForm.protocol = proto.value"
            >
              <div class="option-icon">{{ proto.icon }}</div>
              <div class="option-info">
                <span class="option-name">{{ proto.label }}</span>
                <span class="option-desc">{{ proto.desc }}</span>
              </div>
              <div class="option-check" v-if="createForm.protocol === proto.value">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                  <path d="M20 6L9 17l-5-5"/>
                </svg>
              </div>
            </div>
          </div>
        </div>

        <div class="form-group">
          <label class="form-label">
            <span class="label-icon">></span>
            SUBSCRIPTION_DAYS
          </label>
          <div class="days-selector">
            <button
              v-for="day in [7, 30, 90, 180, 365]"
              :key="day"
              class="day-btn"
              :class="{ active: createForm.days === day }"
              @click="createForm.days = day"
            >
              {{ day }}天
            </button>
            <div class="custom-days">
              <el-input-number
                v-model="createForm.days"
                :min="1"
                :max="365"
                size="small"
                controls-position="right"
              />
            </div>
          </div>
        </div>

        <div class="cost-display">
          <div class="cost-label">ESTIMATED_COST</div>
          <div class="cost-value">
            <span class="currency">¥</span>
            <span class="amount">{{ (createForm.days * 1).toFixed(2) }}</span>
          </div>
          <div class="cost-rate">@ ¥1.00/天</div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <button class="cancel-btn" @click="createDialogVisible = false">
            取消
          </button>
          <button class="confirm-btn" @click="handleCreate" :disabled="createLoading">
            <span v-if="createLoading" class="loading-spinner"></span>
            <span v-else>确认申请</span>
          </button>
        </div>
      </template>
    </el-dialog>

    <!-- Renew Dialog -->
    <el-dialog
      v-model="renewDialogVisible"
      title=""
      width="450"
      :show-close="false"
      class="cyber-dialog"
    >
      <template #header>
        <div class="dialog-header">
          <div class="dialog-title">
            <span class="title-bracket">[</span>
            RENEW_SUBSCRIPTION
            <span class="title-bracket">]</span>
          </div>
          <button class="close-btn" @click="renewDialogVisible = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6L6 18M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </template>

      <div class="dialog-content">
        <div class="renew-info">
          <div class="renew-row">
            <span class="renew-label">PORT_NUMBER</span>
            <span class="renew-value highlight">{{ selectedPort?.port }}</span>
          </div>
          <div class="renew-row">
            <span class="renew-label">CURRENT_EXPIRY</span>
            <span class="renew-value">{{ formatTime(selectedPort?.expire_at) }}</span>
          </div>
        </div>

        <div class="form-group">
          <label class="form-label">
            <span class="label-icon">></span>
            EXTEND_DAYS
          </label>
          <div class="days-selector">
            <button
              v-for="day in [7, 30, 90, 180]"
              :key="day"
              class="day-btn"
              :class="{ active: renewDays === day }"
              @click="renewDays = day"
            >
              {{ day }}天
            </button>
          </div>
        </div>

        <div class="cost-display">
          <div class="cost-label">RENEWAL_COST</div>
          <div class="cost-value">
            <span class="currency">¥</span>
            <span class="amount">{{ (renewDays * 1).toFixed(2) }}</span>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <button class="cancel-btn" @click="renewDialogVisible = false">取消</button>
          <button class="confirm-btn" @click="handleRenew" :disabled="renewLoading">
            <span v-if="renewLoading" class="loading-spinner"></span>
            <span v-else>确认续费</span>
          </button>
        </div>
      </template>
    </el-dialog>

    <!-- Detail Dialog -->
    <el-dialog
      v-model="detailDialogVisible"
      title=""
      width="680"
      :show-close="false"
      class="cyber-dialog detail-dialog"
    >
      <template #header>
        <div class="dialog-header">
          <div class="dialog-title">
            <span class="title-bracket">[</span>
            PORT_{{ selectedPort?.port }}_DETAILS
            <span class="title-bracket">]</span>
          </div>
          <button class="close-btn" @click="detailDialogVisible = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6L6 18M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </template>

      <div class="dialog-content">
        <div class="detail-grid">
          <div class="detail-item">
            <span class="detail-label">PORT</span>
            <span class="detail-value large">{{ selectedPort?.port }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">NAME</span>
            <span class="detail-value">{{ selectedPort?.name || '-' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">PROTOCOL</span>
            <span class="protocol-badge" :class="selectedPort?.protocol">
              {{ protocolNames[selectedPort?.protocol] }}
            </span>
          </div>
          <div class="detail-item">
            <span class="detail-label">STATUS</span>
            <span class="status-badge" :class="selectedPort?.status">
              {{ statusNames[selectedPort?.status] }}
            </span>
          </div>
          <div class="detail-item">
            <span class="detail-label">EXPIRE_AT</span>
            <span class="detail-value">{{ formatTime(selectedPort?.expire_at) }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">CREATED_AT</span>
            <span class="detail-value">{{ formatTime(selectedPort?.created_at) }}</span>
          </div>
        </div>

        <div class="token-block">
          <div class="token-header">
            <span class="token-title">
              <span class="title-icon">#</span>
              CONNECTION_TOKEN
            </span>
            <div class="token-actions">
              <button class="token-btn" @click="copyToken(selectedPort?.token)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="9" y="9" width="13" height="13" rx="2"/>
                  <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                </svg>
                复制
              </button>
              <button class="token-btn warning" @click="handleRegenerateToken">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M23 4v6h-6M1 20v-6h6"/>
                  <path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15"/>
                </svg>
                重新生成
              </button>
            </div>
          </div>
          <div class="token-content">
            <code>{{ selectedPort?.token }}</code>
          </div>
        </div>

        <div class="code-block">
          <div class="code-header">
            <span class="code-title">CONNECTION_EXAMPLE</span>
            <div class="code-tabs">
              <button
                v-for="tab in ['Auto.js', 'EasyClick', '懒人精灵']"
                :key="tab"
                class="code-tab"
                :class="{ active: codeTab === tab }"
                @click="codeTab = tab"
              >
                {{ tab }}
              </button>
            </div>
          </div>
          <pre class="code-content"><code>{{ getConnectCode() }}</code></pre>
          <button class="copy-code-btn" @click="copyCode">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="9" y="9" width="13" height="13" rx="2"/>
              <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
            </svg>
            复制代码
          </button>
        </div>
      </div>
    </el-dialog>

    <!-- Delete Confirmation -->
    <el-dialog
      v-model="deleteDialogVisible"
      title=""
      width="400"
      :show-close="false"
      class="cyber-dialog delete-dialog"
    >
      <template #header>
        <div class="dialog-header warning">
          <div class="dialog-title">
            <span class="title-bracket">[</span>
            CONFIRM_DELETE
            <span class="title-bracket">]</span>
          </div>
        </div>
      </template>

      <div class="dialog-content">
        <div class="warning-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 9v4m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
          </svg>
        </div>
        <p class="warning-text">确定要释放端口 <strong>{{ portToDelete?.port }}</strong> 吗？</p>
        <p class="warning-hint">此操作不可撤销，端口将被回收至端口池</p>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <button class="cancel-btn" @click="deleteDialogVisible = false">取消</button>
          <button class="confirm-btn danger" @click="handleDelete">
            确认释放
          </button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'
import dayjs from 'dayjs'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const loading = ref(false)
const ports = ref([])
const codeTab = ref('Auto.js')

const protocols = [
  { value: 'autojs', label: 'Auto.js', icon: 'AJ', desc: 'Auto.js Pro 协议' },
  { value: 'lazyman', label: '懒人精灵', icon: 'LR', desc: '懒人精灵协议' },
  { value: 'easyclick', label: 'EasyClick', icon: 'EC', desc: 'EasyClick 协议' },
  { value: 'generic', label: '通用协议', icon: 'WS', desc: '通用 WebSocket' }
]

const protocolNames = {
  autojs: 'Auto.js',
  lazyman: '懒人精灵',
  easyclick: 'EasyClick',
  generic: '通用协议'
}

const statusNames = {
  active: '运行中',
  expired: '已过期',
  disabled: '已禁用'
}

// Computed stats
const activePorts = computed(() => ports.value.filter(p => p.status === 'active').length)
const expiredPorts = computed(() => ports.value.filter(p => p.status === 'expired').length)
const totalDevices = computed(() => ports.value.reduce((sum, p) => sum + (p.device_count || 0), 0))

// Create dialog
const createDialogVisible = ref(false)
const createLoading = ref(false)
const createForm = ref({
  name: '',
  protocol: 'autojs',
  days: 30
})

// Renew dialog
const renewDialogVisible = ref(false)
const renewLoading = ref(false)
const renewDays = ref(30)
const selectedPort = ref(null)

// Detail dialog
const detailDialogVisible = ref(false)

// Delete dialog
const deleteDialogVisible = ref(false)
const portToDelete = ref(null)

function formatTime(time) {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-'
}

function isExpiringSoon(time) {
  return time && dayjs(time).diff(dayjs(), 'day') < 7
}

function copyToken(token) {
  navigator.clipboard.writeText(token)
  ElMessage.success('密钥已复制到剪贴板')
}

function copyCode() {
  navigator.clipboard.writeText(getConnectCode())
  ElMessage.success('代码已复制到剪贴板')
}

function getConnectCode() {
  const host = window.location.hostname
  const port = selectedPort.value?.port
  const token = selectedPort.value?.token

  if (codeTab.value === 'Auto.js') {
    return `// Auto.js Pro 连接示例
var ws = new WebSocket("ws://${host}:${port}");

ws.on("open", function() {
    ws.send(JSON.stringify({
        type: "register",
        token: "${token}",
        device: {
            id: device.getAndroidId(),
            name: device.brand + " " + device.model,
            os: "Android " + device.release
        }
    }));
});

ws.on("message", function(msg) {
    var data = JSON.parse(msg);
    if (data.type === "exec") {
        eval(data.code);
    }
});`
  } else if (codeTab.value === 'EasyClick') {
    return `// EasyClick 连接示例
var ws = new WebSocket("ws://${host}:${port}");

ws.onOpen(function() {
    ws.sendText(JSON.stringify({
        type: "register",
        token: "${token}",
        device: {
            id: device.getDeviceId(),
            name: device.getBrand() + " " + device.getModel(),
            os: "Android " + device.getOsVersion()
        }
    }));
});

ws.onMessage(function(msg) {
    var data = JSON.parse(msg);
    if (data.type === "exec") {
        eval(data.code);
    }
});`
  } else {
    return `-- 懒人精灵连接示例
local ws = require("websocket")
local json = require("json")

local client = ws.connect("ws://${host}:${port}")

client:send(json.encode({
    type = "register",
    token = "${token}",
    device = {
        id = getDeviceID(),
        name = getDeviceName(),
        os = "Android"
    }
}))

while true do
    local msg = client:receive()
    if msg then
        local data = json.decode(msg)
        if data.type == "exec" then
            load(data.code)()
        end
    end
end`
  }
}

async function fetchPorts() {
  loading.value = true
  try {
    const res = await api.ports.list()
    if (res.code === 0) {
      ports.value = res.data || []
    }
  } finally {
    loading.value = false
  }
}

function showCreateDialog() {
  createForm.value = { name: '', protocol: 'autojs', days: 30 }
  createDialogVisible.value = true
}

async function handleCreate() {
  createLoading.value = true
  try {
    const res = await api.ports.create(createForm.value)
    if (res.code === 0) {
      ElMessage.success('端口申请成功')
      createDialogVisible.value = false
      fetchPorts()
      userStore.refreshProfile()
    } else {
      ElMessage.error(res.message)
    }
  } finally {
    createLoading.value = false
  }
}

function showRenewDialog(port) {
  selectedPort.value = port
  renewDays.value = 30
  renewDialogVisible.value = true
}

async function handleRenew() {
  renewLoading.value = true
  try {
    const res = await api.ports.renew(selectedPort.value.id, renewDays.value)
    if (res.code === 0) {
      ElMessage.success('续费成功')
      renewDialogVisible.value = false
      fetchPorts()
      userStore.refreshProfile()
    } else {
      ElMessage.error(res.message)
    }
  } finally {
    renewLoading.value = false
  }
}

function showDetailDialog(port) {
  selectedPort.value = port
  codeTab.value = 'Auto.js'
  detailDialogVisible.value = true
}

async function handleChangePort(port) {
  try {
    await ElMessageBox.confirm(
      '更换端口后，原端口将失效，连接密钥也会更新，确定要更换吗？',
      '确认更换',
      { type: 'warning' }
    )
    const res = await api.ports.change(port.id)
    if (res.code === 0) {
      ElMessage.success(`端口已更换为 ${res.data.port}`)
      fetchPorts()
    } else {
      ElMessage.error(res.message)
    }
  } catch {}
}

function confirmDelete(port) {
  portToDelete.value = port
  deleteDialogVisible.value = true
}

async function handleDelete() {
  const res = await api.ports.delete(portToDelete.value.id)
  if (res.code === 0) {
    ElMessage.success('端口已释放')
    deleteDialogVisible.value = false
    fetchPorts()
  } else {
    ElMessage.error(res.message)
  }
}

async function handleRegenerateToken() {
  try {
    await ElMessageBox.confirm('重新生成密钥后，所有设备需要重新连接，确定吗？', '确认', { type: 'warning' })
    const res = await api.ports.regenerateToken(selectedPort.value.id)
    if (res.code === 0) {
      ElMessage.success('密钥已更新')
      selectedPort.value = res.data
      fetchPorts()
    } else {
      ElMessage.error(res.message)
    }
  } catch {}
}

onMounted(() => {
  fetchPorts()
})
</script>

<style scoped>
.ports-page {
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
  background: linear-gradient(135deg, rgba(0, 240, 255, 0.1), rgba(168, 85, 247, 0.1));
  border: 1px solid var(--neon-cyan);
  border-radius: var(--radius-lg);
  color: var(--neon-cyan);
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

.create-btn {
  position: relative;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 24px;
  background: linear-gradient(135deg, var(--neon-cyan), var(--neon-cyan-dim));
  border: none;
  border-radius: var(--radius-lg);
  color: var(--bg-void);
  font-family: var(--font-mono);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.3s ease;
}

.create-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(0, 240, 255, 0.3);
}

.btn-icon {
  font-size: 18px;
  font-weight: 700;
}

.btn-glow {
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  transform: translateX(-100%);
}

.create-btn:hover .btn-glow {
  animation: glow-sweep 0.6s ease;
}

@keyframes glow-sweep {
  to { transform: translateX(100%); }
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

.stat-value.active {
  color: var(--neon-green);
}

.stat-value.expired {
  color: var(--status-error);
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

/* Ports Grid */
.ports-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 20px;
}

/* Port Card */
.port-card {
  position: relative;
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  padding: 20px;
  transition: all 0.3s ease;
}

.port-card:hover {
  border-color: var(--neon-cyan-dim);
  transform: translateY(-2px);
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
}

.port-card.status-active {
  border-left: 3px solid var(--neon-green);
}

.port-card.status-expired {
  border-left: 3px solid var(--status-error);
}

.port-card.status-disabled {
  border-left: 3px solid var(--text-muted);
  opacity: 0.7;
}

.port-card.expiring-soon {
  animation: warning-pulse 2s ease-in-out infinite;
}

@keyframes warning-pulse {
  0%, 100% { box-shadow: 0 0 0 0 rgba(245, 158, 11, 0); }
  50% { box-shadow: 0 0 0 4px rgba(245, 158, 11, 0.2); }
}

/* Card corners */
.port-card .corner {
  position: absolute;
  width: 12px;
  height: 12px;
  border-color: var(--border-default);
  opacity: 0.5;
  transition: all 0.3s ease;
}

.port-card:hover .corner {
  border-color: var(--neon-cyan);
  opacity: 1;
}

.corner.tl { top: 8px; left: 8px; border-top: 2px solid; border-left: 2px solid; }
.corner.tr { top: 8px; right: 8px; border-top: 2px solid; border-right: 2px solid; }
.corner.bl { bottom: 8px; left: 8px; border-bottom: 2px solid; border-left: 2px solid; }
.corner.br { bottom: 8px; right: 8px; border-bottom: 2px solid; border-right: 2px solid; }

/* Port Header */
.port-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 16px;
}

.port-number {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.port-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.port-value {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 800;
  color: var(--neon-cyan);
  text-shadow: 0 0 20px rgba(0, 240, 255, 0.3);
}

.port-status {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: var(--bg-elevated);
  border-radius: var(--radius-full);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--neon-green);
  box-shadow: 0 0 8px var(--neon-green);
}

.port-card.status-expired .status-dot {
  background: var(--status-error);
  box-shadow: 0 0 8px var(--status-error);
}

.port-card.status-disabled .status-dot {
  background: var(--text-muted);
  box-shadow: none;
}

.status-text {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-secondary);
}

/* Port Info */
.port-info {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border-subtle);
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.info-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.info-value {
  font-family: var(--font-mono);
  font-size: 13px;
  color: var(--text-secondary);
}

.info-value.devices {
  display: flex;
  align-items: center;
  gap: 6px;
}

.device-count {
  font-weight: 700;
  color: var(--neon-purple);
}

.device-label {
  font-size: 10px;
  color: var(--text-muted);
}

.info-value.expire.warning {
  color: var(--status-warning);
}

/* Protocol Badge */
.protocol-badge {
  padding: 3px 10px;
  border-radius: var(--radius-full);
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 600;
}

.protocol-badge.autojs {
  background: rgba(0, 240, 255, 0.15);
  color: var(--neon-cyan);
  border: 1px solid rgba(0, 240, 255, 0.3);
}

.protocol-badge.lazyman {
  background: rgba(168, 85, 247, 0.15);
  color: var(--neon-purple);
  border: 1px solid rgba(168, 85, 247, 0.3);
}

.protocol-badge.easyclick {
  background: rgba(34, 197, 94, 0.15);
  color: var(--neon-green);
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.protocol-badge.generic {
  background: rgba(245, 158, 11, 0.15);
  color: var(--status-warning);
  border: 1px solid rgba(245, 158, 11, 0.3);
}

/* Token Section */
.token-section {
  margin-bottom: 16px;
}

.token-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 1px;
  margin-bottom: 8px;
}

.label-icon {
  color: var(--neon-purple);
}

.token-value {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
}

.token-value code {
  flex: 1;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-secondary);
  word-break: break-all;
}

.copy-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: transparent;
  border: 1px solid var(--border-default);
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.2s ease;
}

.copy-btn:hover {
  border-color: var(--neon-cyan);
  color: var(--neon-cyan);
}

.copy-btn svg {
  width: 14px;
  height: 14px;
}

/* Port Actions */
.port-actions {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 10px 8px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn svg {
  width: 16px;
  height: 16px;
}

.action-btn.detail:hover {
  border-color: var(--neon-cyan);
  color: var(--neon-cyan);
  background: rgba(0, 240, 255, 0.05);
}

.action-btn.renew:hover {
  border-color: var(--neon-green);
  color: var(--neon-green);
  background: rgba(34, 197, 94, 0.05);
}

.action-btn.change:hover {
  border-color: var(--status-warning);
  color: var(--status-warning);
  background: rgba(245, 158, 11, 0.05);
}

.action-btn.delete:hover {
  border-color: var(--status-error);
  color: var(--status-error);
  background: rgba(239, 68, 68, 0.05);
}

/* Dialog Styles */
.cyber-dialog :deep(.el-dialog) {
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-xl);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
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
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-subtle);
}

.dialog-header.warning {
  background: rgba(239, 68, 68, 0.1);
  border-bottom-color: rgba(239, 68, 68, 0.3);
}

.dialog-title {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: 1px;
}

.title-bracket {
  color: var(--neon-cyan);
}

.dialog-header.warning .title-bracket {
  color: var(--status-error);
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
  background: linear-gradient(135deg, var(--neon-cyan), var(--neon-cyan-dim));
  border: none;
  border-radius: var(--radius-md);
  color: var(--bg-void);
  font-family: var(--font-mono);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.confirm-btn:hover {
  box-shadow: 0 4px 20px rgba(0, 240, 255, 0.3);
}

.confirm-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.confirm-btn.danger {
  background: linear-gradient(135deg, var(--status-error), #dc2626);
}

.confirm-btn.danger:hover {
  box-shadow: 0 4px 20px rgba(239, 68, 68, 0.3);
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Form Styles */
.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 1px;
  margin-bottom: 10px;
}

.form-label .label-icon {
  color: var(--neon-cyan);
}

.cyber-input :deep(.el-input__wrapper) {
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  padding: 12px 14px;
}

.cyber-input :deep(.el-input__wrapper:hover) {
  border-color: var(--neon-cyan-dim);
}

.cyber-input :deep(.el-input__wrapper.is-focus) {
  border-color: var(--neon-cyan);
  box-shadow: 0 0 0 3px rgba(0, 240, 255, 0.1);
}

/* Protocol Select */
.protocol-select {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.protocol-option {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
}

.protocol-option:hover {
  border-color: var(--neon-cyan-dim);
}

.protocol-option.active {
  border-color: var(--neon-cyan);
  background: rgba(0, 240, 255, 0.05);
}

.option-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-surface);
  border-radius: var(--radius-md);
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 700;
  color: var(--neon-cyan);
}

.option-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.option-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.option-desc {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
}

.option-check {
  width: 20px;
  height: 20px;
  color: var(--neon-cyan);
}

.option-check svg {
  width: 100%;
  height: 100%;
}

/* Days Selector */
.days-selector {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.day-btn {
  padding: 10px 18px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.day-btn:hover {
  border-color: var(--neon-cyan-dim);
}

.day-btn.active {
  border-color: var(--neon-cyan);
  background: rgba(0, 240, 255, 0.1);
  color: var(--neon-cyan);
}

.custom-days {
  display: flex;
  align-items: center;
}

.custom-days :deep(.el-input-number) {
  width: 100px;
}

/* Cost Display */
.cost-display {
  display: flex;
  align-items: baseline;
  gap: 12px;
  padding: 20px;
  background: linear-gradient(135deg, rgba(0, 240, 255, 0.05), rgba(168, 85, 247, 0.05));
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  margin-top: 20px;
}

.cost-label {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.cost-value {
  display: flex;
  align-items: baseline;
}

.cost-value .currency {
  font-size: 16px;
  color: var(--text-secondary);
  margin-right: 2px;
}

.cost-value .amount {
  font-family: var(--font-display);
  font-size: 32px;
  font-weight: 800;
  color: var(--neon-cyan);
}

.cost-rate {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
  margin-left: auto;
}

/* Renew Info */
.renew-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  background: var(--bg-elevated);
  border-radius: var(--radius-lg);
  margin-bottom: 20px;
}

.renew-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.renew-label {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.renew-value {
  font-family: var(--font-mono);
  font-size: 14px;
  color: var(--text-secondary);
}

.renew-value.highlight {
  font-size: 20px;
  font-weight: 700;
  color: var(--neon-cyan);
}

/* Detail Grid */
.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 14px;
  background: var(--bg-elevated);
  border-radius: var(--radius-md);
}

.detail-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.detail-value {
  font-family: var(--font-mono);
  font-size: 14px;
  color: var(--text-primary);
}

.detail-value.large {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 800;
  color: var(--neon-cyan);
}

.status-badge {
  display: inline-flex;
  padding: 4px 12px;
  border-radius: var(--radius-full);
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 600;
}

.status-badge.active {
  background: rgba(34, 197, 94, 0.15);
  color: var(--neon-green);
}

.status-badge.expired {
  background: rgba(239, 68, 68, 0.15);
  color: var(--status-error);
}

.status-badge.disabled {
  background: rgba(148, 163, 184, 0.15);
  color: var(--text-muted);
}

/* Token Block */
.token-block {
  margin-bottom: 24px;
}

.token-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.token-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.title-icon {
  color: var(--neon-purple);
}

.token-actions {
  display: flex;
  gap: 8px;
}

.token-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: transparent;
  border: 1px solid var(--border-default);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.token-btn svg {
  width: 14px;
  height: 14px;
}

.token-btn:hover {
  border-color: var(--neon-cyan);
  color: var(--neon-cyan);
}

.token-btn.warning:hover {
  border-color: var(--status-warning);
  color: var(--status-warning);
}

.token-content {
  padding: 14px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
}

.token-content code {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--neon-cyan);
  word-break: break-all;
}

/* Code Block */
.code-block {
  position: relative;
}

.code-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.code-title {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.code-tabs {
  display: flex;
  gap: 4px;
}

.code-tab {
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

.code-tab:hover {
  border-color: var(--neon-cyan-dim);
  color: var(--text-secondary);
}

.code-tab.active {
  border-color: var(--neon-cyan);
  background: rgba(0, 240, 255, 0.1);
  color: var(--neon-cyan);
}

.code-content {
  padding: 16px;
  background: var(--bg-void);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  margin: 0;
  overflow-x: auto;
}

.code-content code {
  font-family: var(--font-mono);
  font-size: 12px;
  line-height: 1.6;
  color: var(--text-secondary);
  white-space: pre;
}

.copy-code-btn {
  position: absolute;
  top: 44px;
  right: 10px;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--bg-surface);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.copy-code-btn svg {
  width: 12px;
  height: 12px;
}

.copy-code-btn:hover {
  border-color: var(--neon-cyan);
  color: var(--neon-cyan);
}

/* Delete Dialog */
.delete-dialog .warning-icon {
  width: 60px;
  height: 60px;
  margin: 0 auto 20px;
  color: var(--status-error);
  animation: pulse-warning 1.5s ease-in-out infinite;
}

@keyframes pulse-warning {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.7; transform: scale(1.05); }
}

.delete-dialog .warning-icon svg {
  width: 100%;
  height: 100%;
}

.delete-dialog .warning-text {
  text-align: center;
  font-size: 16px;
  color: var(--text-primary);
  margin: 0 0 10px 0;
}

.delete-dialog .warning-text strong {
  color: var(--neon-cyan);
}

.delete-dialog .warning-hint {
  text-align: center;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
  margin: 0;
}

/* Responsive */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .stats-bar {
    flex-wrap: wrap;
  }

  .ports-grid {
    grid-template-columns: 1fr;
  }

  .port-actions {
    grid-template-columns: repeat(2, 1fr);
  }

  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>
