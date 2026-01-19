<template>
  <div class="dashboard">
    <!-- Stats Grid -->
    <div class="stats-grid">
      <div class="stat-card" v-for="(stat, index) in statsCards" :key="stat.label">
        <div class="stat-glow" :style="{ background: stat.glowColor }"></div>
        <div class="stat-icon" :style="{ background: stat.iconBg, color: stat.iconColor }">
          <el-icon :size="24"><component :is="stat.icon" /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-value" :style="{ color: stat.valueColor }">{{ stat.value }}</div>
        </div>
        <div class="stat-trend" v-if="stat.trend">
          <span :class="['trend-badge', stat.trend > 0 ? 'up' : 'down']">
            {{ stat.trend > 0 ? '+' : '' }}{{ stat.trend }}%
          </span>
        </div>
      </div>
    </div>

    <!-- Main content grid -->
    <div class="content-grid">
      <!-- Ports panel -->
      <div class="panel ports-panel">
        <div class="panel-header">
          <div class="panel-title">
            <el-icon><Connection /></el-icon>
            <span>PORT_STATUS</span>
          </div>
          <el-button size="small" @click="$router.push('/ports')">
            管理端口
            <el-icon class="btn-arrow"><ArrowRight /></el-icon>
          </el-button>
        </div>
        <div class="panel-content">
          <div v-if="ports.length > 0" class="port-list">
            <div
              v-for="port in ports.slice(0, 5)"
              :key="port.id"
              class="port-item"
              :class="{ active: port.status === 'active' }"
            >
              <div class="port-info">
                <span class="port-number">:{{ port.port }}</span>
                <span class="port-name">{{ port.name || 'Unnamed' }}</span>
              </div>
              <div class="port-meta">
                <span class="port-protocol">{{ port.protocol?.toUpperCase() }}</span>
                <span class="port-devices">
                  <el-icon><Iphone /></el-icon>
                  {{ port.device_count || 0 }}
                </span>
              </div>
              <div class="port-status">
                <span class="status-indicator" :class="port.status"></span>
              </div>
            </div>
          </div>
          <div v-else class="empty-state">
            <div class="empty-icon">
              <el-icon :size="48"><Connection /></el-icon>
            </div>
            <p>暂无活跃端口</p>
            <el-button type="primary" size="small" @click="$router.push('/ports')">
              申请端口
            </el-button>
          </div>
        </div>
      </div>

      <!-- Devices panel -->
      <div class="panel devices-panel">
        <div class="panel-header">
          <div class="panel-title">
            <el-icon><Iphone /></el-icon>
            <span>ONLINE_DEVICES</span>
          </div>
          <el-button size="small" @click="$router.push('/devices')">
            查看全部
            <el-icon class="btn-arrow"><ArrowRight /></el-icon>
          </el-button>
        </div>
        <div class="panel-content">
          <div v-if="onlineDevices.length > 0" class="device-list">
            <div
              v-for="device in onlineDevices.slice(0, 5)"
              :key="device.id"
              class="device-item"
            >
              <div class="device-avatar">
                <el-icon><Iphone /></el-icon>
              </div>
              <div class="device-info">
                <span class="device-name">{{ device.device_name }}</span>
                <span class="device-port">Port :{{ device.port }}</span>
              </div>
              <div class="device-stats">
                <span class="device-ping" :class="getPingClass(device.last_ping)">
                  {{ device.last_ping || 0 }}ms
                </span>
              </div>
            </div>
          </div>
          <div v-else class="empty-state">
            <div class="empty-icon">
              <el-icon :size="48"><Iphone /></el-icon>
            </div>
            <p>暂无在线设备</p>
            <span class="empty-hint">使用连接密钥接入设备</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick start guide -->
    <div class="quick-start-panel">
      <div class="panel-header">
        <div class="panel-title">
          <el-icon><Guide /></el-icon>
          <span>QUICK_START</span>
        </div>
      </div>
      <div class="panel-content">
        <div class="steps-container">
          <div
            v-for="(step, index) in quickStartSteps"
            :key="index"
            class="step-item"
            :class="{ completed: index < currentStep, active: index === currentStep }"
          >
            <div class="step-number">
              <span v-if="index < currentStep" class="check-icon">✓</span>
              <span v-else>{{ index + 1 }}</span>
            </div>
            <div class="step-content">
              <div class="step-title">{{ step.title }}</div>
              <div class="step-desc">{{ step.desc }}</div>
            </div>
            <div class="step-connector" v-if="index < quickStartSteps.length - 1"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../api'
import { Connection, Iphone, Wallet, TrendCharts, ArrowRight, Guide } from '@element-plus/icons-vue'

const user = ref(null)
const stats = ref({})
const ports = ref([])
const onlineDevices = ref([])

const statsCards = computed(() => [
  {
    label: 'ACTIVE_PORTS',
    value: stats.value.active_port_count || 0,
    icon: Connection,
    iconBg: 'rgba(0, 240, 255, 0.15)',
    iconColor: 'var(--neon-cyan)',
    valueColor: 'var(--neon-cyan)',
    glowColor: 'var(--neon-cyan)'
  },
  {
    label: 'ONLINE_DEVICES',
    value: stats.value.online_device_count || 0,
    icon: Iphone,
    iconBg: 'rgba(74, 222, 128, 0.15)',
    iconColor: 'var(--neon-green)',
    valueColor: 'var(--neon-green)',
    glowColor: 'var(--neon-green)'
  },
  {
    label: 'BALANCE',
    value: `¥${user.value?.balance?.toFixed(2) || '0.00'}`,
    icon: Wallet,
    iconBg: 'rgba(251, 191, 36, 0.15)',
    iconColor: 'var(--neon-yellow)',
    valueColor: 'var(--neon-yellow)',
    glowColor: 'var(--neon-yellow)'
  },
  {
    label: 'TOTAL_SPENT',
    value: `¥${stats.value.total_consume?.toFixed(2) || '0.00'}`,
    icon: TrendCharts,
    iconBg: 'rgba(168, 85, 247, 0.15)',
    iconColor: 'var(--neon-purple)',
    valueColor: 'var(--neon-purple)',
    glowColor: 'var(--neon-purple)'
  }
])

const quickStartSteps = [
  { title: '创建账号', desc: '注册并登录平台' },
  { title: '充值余额', desc: '联系管理员充值余额' },
  { title: '申请端口', desc: '在端口管理中申请调试端口' },
  { title: '连接设备', desc: '使用连接密钥接入设备' }
]

const currentStep = computed(() => {
  if (onlineDevices.value.length > 0) return 4
  if (ports.value.length > 0) return 3
  if (user.value?.balance > 0) return 2
  return 1
})

function getPingClass(ping) {
  if (!ping) return 'unknown'
  if (ping < 100) return 'good'
  if (ping < 300) return 'normal'
  return 'bad'
}

async function fetchData() {
  const res = await api.user.getDashboard()
  if (res.code === 0) {
    user.value = res.data.user
    stats.value = res.data.stats || {}
    ports.value = res.data.ports || []
    onlineDevices.value = res.data.online_devices || []
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.stat-card {
  position: relative;
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.stat-card:hover {
  border-color: var(--border-glow);
  transform: translateY(-2px);
}

.stat-glow {
  position: absolute;
  top: 0;
  right: 0;
  width: 100px;
  height: 100px;
  opacity: 0.1;
  filter: blur(40px);
  border-radius: 50%;
}

.stat-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-label {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 500;
  color: var(--text-muted);
  letter-spacing: 1px;
  margin-bottom: 4px;
}

.stat-value {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
}

.stat-trend {
  position: absolute;
  top: 12px;
  right: 12px;
}

.trend-badge {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
}

.trend-badge.up {
  background: rgba(74, 222, 128, 0.15);
  color: var(--neon-green);
}

.trend-badge.down {
  background: rgba(244, 63, 94, 0.15);
  color: var(--neon-red);
}

/* Content Grid */
.content-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

/* Panel styles */
.panel {
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-subtle);
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-family: var(--font-mono);
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: 1px;
}

.panel-title .el-icon {
  color: var(--neon-cyan);
}

.panel-content {
  padding: 16px;
}

.btn-arrow {
  margin-left: 4px;
  transition: transform 0.2s ease;
}

.panel-header .el-button:hover .btn-arrow {
  transform: translateX(4px);
}

/* Port list */
.port-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.port-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 16px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  transition: all 0.2s ease;
}

.port-item:hover {
  border-color: var(--border-default);
  background: var(--bg-hover);
}

.port-item.active {
  border-left: 3px solid var(--neon-green);
}

.port-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.port-number {
  font-family: var(--font-mono);
  font-size: 16px;
  font-weight: 600;
  color: var(--neon-cyan);
}

.port-name {
  font-size: 12px;
  color: var(--text-muted);
}

.port-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.port-protocol {
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 600;
  padding: 2px 8px;
  background: rgba(0, 240, 255, 0.1);
  color: var(--neon-cyan);
  border-radius: 4px;
  letter-spacing: 0.5px;
}

.port-devices {
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-secondary);
}

.port-status {
  display: flex;
  align-items: center;
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-indicator.active {
  background: var(--neon-green);
  box-shadow: 0 0 8px var(--neon-green);
}

.status-indicator.expired,
.status-indicator.disabled {
  background: var(--text-muted);
}

/* Device list */
.device-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.device-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  transition: all 0.2s ease;
}

.device-item:hover {
  border-color: var(--border-default);
  background: var(--bg-hover);
}

.device-avatar {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(74, 222, 128, 0.15);
  color: var(--neon-green);
  border-radius: var(--radius-md);
}

.device-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.device-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.device-port {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
}

.device-stats {
  display: flex;
  align-items: center;
}

.device-ping {
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: var(--radius-sm);
}

.device-ping.good {
  background: rgba(74, 222, 128, 0.15);
  color: var(--neon-green);
}

.device-ping.normal {
  background: rgba(251, 191, 36, 0.15);
  color: var(--neon-yellow);
}

.device-ping.bad {
  background: rgba(244, 63, 94, 0.15);
  color: var(--neon-red);
}

.device-ping.unknown {
  background: var(--bg-hover);
  color: var(--text-muted);
}

/* Empty state */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.empty-icon {
  color: var(--text-muted);
  opacity: 0.5;
  margin-bottom: 16px;
}

.empty-state p {
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.empty-hint {
  font-size: 12px;
  color: var(--text-muted);
}

.empty-state .el-button {
  margin-top: 16px;
}

/* Quick start panel */
.quick-start-panel {
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
}

.steps-container {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
}

.step-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  position: relative;
  padding: 0 20px;
}

.step-number {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-elevated);
  border: 2px solid var(--border-default);
  border-radius: 50%;
  font-family: var(--font-mono);
  font-size: 14px;
  font-weight: 600;
  color: var(--text-muted);
  margin-bottom: 12px;
  transition: all 0.3s ease;
}

.step-item.completed .step-number {
  background: var(--neon-green);
  border-color: var(--neon-green);
  color: var(--bg-void);
}

.step-item.active .step-number {
  background: var(--neon-cyan);
  border-color: var(--neon-cyan);
  color: var(--bg-void);
  box-shadow: 0 0 20px rgba(0, 240, 255, 0.3);
}

.check-icon {
  font-size: 18px;
}

.step-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.step-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.step-item.completed .step-title {
  color: var(--neon-green);
}

.step-item.active .step-title {
  color: var(--neon-cyan);
}

.step-desc {
  font-size: 12px;
  color: var(--text-muted);
}

.step-connector {
  position: absolute;
  top: 20px;
  left: calc(50% + 30px);
  width: calc(100% - 60px);
  height: 2px;
  background: var(--border-default);
}

.step-item.completed .step-connector {
  background: var(--neon-green);
}

/* Responsive */
@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 900px) {
  .content-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .steps-container {
    flex-direction: column;
    gap: 16px;
  }

  .step-item {
    flex-direction: row;
    text-align: left;
    padding: 12px 0;
  }

  .step-number {
    margin-bottom: 0;
    margin-right: 16px;
  }

  .step-connector {
    display: none;
  }
}
</style>
