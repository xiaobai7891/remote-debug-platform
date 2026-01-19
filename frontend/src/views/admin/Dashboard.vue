<template>
  <div class="admin-dashboard">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-icon">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="3" y="3" width="7" height="7" rx="1" stroke="currentColor" stroke-width="2"/>
            <rect x="14" y="3" width="7" height="7" rx="1" stroke="currentColor" stroke-width="2"/>
            <rect x="3" y="14" width="7" height="7" rx="1" stroke="currentColor" stroke-width="2"/>
            <rect x="14" y="14" width="7" height="7" rx="1" stroke="currentColor" stroke-width="2"/>
          </svg>
        </div>
        <div class="header-text">
          <h1 class="page-title">SYSTEM_OVERVIEW</h1>
          <p class="page-subtitle">// 系统数据概览</p>
        </div>
      </div>
      <div class="header-actions">
        <div class="last-updated">
          <span class="update-label">LAST_UPDATE:</span>
          <span class="update-time">{{ lastUpdated }}</span>
        </div>
        <el-button class="refresh-btn" @click="fetchStats" :loading="loading">
          <el-icon><Refresh /></el-icon>
          <span>刷新</span>
        </el-button>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="stats-grid" v-loading="loading">
      <div class="stat-card users">
        <div class="card-background">
          <div class="hex-pattern"></div>
        </div>
        <div class="card-content">
          <div class="stat-header">
            <div class="stat-icon">
              <el-icon><User /></el-icon>
            </div>
            <span class="stat-code">USR</span>
          </div>
          <div class="stat-value">{{ stats.user_count || 0 }}</div>
          <div class="stat-label">REGISTERED_USERS</div>
          <div class="stat-footer">
            <span class="footer-text">注册用户总数</span>
          </div>
        </div>
        <div class="card-glow"></div>
      </div>

      <div class="stat-card ports">
        <div class="card-background">
          <div class="hex-pattern"></div>
        </div>
        <div class="card-content">
          <div class="stat-header">
            <div class="stat-icon">
              <el-icon><Connection /></el-icon>
            </div>
            <span class="stat-code">PRT</span>
          </div>
          <div class="stat-value">
            {{ stats.active_port_count || 0 }}
            <span class="stat-total">/ {{ stats.port_count || 0 }}</span>
          </div>
          <div class="stat-label">ACTIVE_PORTS</div>
          <div class="stat-footer">
            <div class="progress-mini">
              <div class="progress-fill" :style="{ width: portUsagePercent + '%' }"></div>
            </div>
            <span class="footer-percent">{{ portUsagePercent.toFixed(1) }}%</span>
          </div>
        </div>
        <div class="card-glow"></div>
      </div>

      <div class="stat-card devices">
        <div class="card-background">
          <div class="hex-pattern"></div>
        </div>
        <div class="card-content">
          <div class="stat-header">
            <div class="stat-icon">
              <el-icon><Iphone /></el-icon>
            </div>
            <span class="stat-code">DEV</span>
          </div>
          <div class="stat-value">{{ stats.online_device_count || 0 }}</div>
          <div class="stat-label">ONLINE_DEVICES</div>
          <div class="stat-footer">
            <span class="status-indicator online"></span>
            <span class="footer-text">实时在线</span>
          </div>
        </div>
        <div class="card-glow"></div>
      </div>

      <div class="stat-card revenue">
        <div class="card-background">
          <div class="hex-pattern"></div>
        </div>
        <div class="card-content">
          <div class="stat-header">
            <div class="stat-icon">
              <el-icon><Wallet /></el-icon>
            </div>
            <span class="stat-code">REV</span>
          </div>
          <div class="stat-value">
            <span class="currency">¥</span>
            {{ formatNumber(stats.total_recharge || 0) }}
          </div>
          <div class="stat-label">TOTAL_REVENUE</div>
          <div class="stat-footer">
            <span class="footer-text">累计充值</span>
          </div>
        </div>
        <div class="card-glow"></div>
      </div>
    </div>

    <!-- Dashboard Panels -->
    <div class="dashboard-panels">
      <!-- Revenue Panel -->
      <div class="panel revenue-panel">
        <div class="panel-header">
          <div class="header-left">
            <div class="panel-icon">
              <el-icon><TrendCharts /></el-icon>
            </div>
            <h2 class="panel-title">REVENUE_ANALYTICS</h2>
          </div>
          <span class="panel-subtitle">// 收入统计</span>
        </div>
        <div class="panel-content">
          <div class="revenue-grid">
            <div class="revenue-item total-recharge">
              <div class="item-header">
                <span class="item-label">TOTAL_RECHARGE</span>
                <el-icon class="item-icon"><Top /></el-icon>
              </div>
              <div class="item-value positive">
                +¥{{ formatNumber(stats.total_recharge || 0) }}
              </div>
              <div class="item-desc">累计充值金额</div>
            </div>

            <div class="revenue-item total-consume">
              <div class="item-header">
                <span class="item-label">TOTAL_CONSUME</span>
                <el-icon class="item-icon"><Bottom /></el-icon>
              </div>
              <div class="item-value negative">
                -¥{{ formatNumber(stats.total_consume || 0) }}
              </div>
              <div class="item-desc">累计消费金额</div>
            </div>

            <div class="revenue-item today-recharge">
              <div class="item-header">
                <span class="item-label">TODAY_RECHARGE</span>
                <el-icon class="item-icon"><Calendar /></el-icon>
              </div>
              <div class="item-value positive">
                +¥{{ formatNumber(stats.today_recharge || 0) }}
              </div>
              <div class="item-desc">今日充值</div>
            </div>

            <div class="revenue-item today-consume">
              <div class="item-header">
                <span class="item-label">TODAY_CONSUME</span>
                <el-icon class="item-icon"><Calendar /></el-icon>
              </div>
              <div class="item-value negative">
                -¥{{ formatNumber(stats.today_consume || 0) }}
              </div>
              <div class="item-desc">今日消费</div>
            </div>
          </div>

          <div class="revenue-summary">
            <div class="summary-item">
              <span class="summary-label">NET_REVENUE:</span>
              <span class="summary-value" :class="netRevenue >= 0 ? 'positive' : 'negative'">
                {{ netRevenue >= 0 ? '+' : '' }}¥{{ formatNumber(netRevenue) }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- System Status Panel -->
      <div class="panel status-panel">
        <div class="panel-header">
          <div class="header-left">
            <div class="panel-icon">
              <el-icon><Monitor /></el-icon>
            </div>
            <h2 class="panel-title">SYSTEM_STATUS</h2>
          </div>
          <span class="panel-subtitle">// 系统状态</span>
        </div>
        <div class="panel-content">
          <div class="status-grid">
            <div class="status-row">
              <div class="row-label">
                <span class="label-icon">></span>
                REGISTERED_USERS
              </div>
              <div class="row-value">{{ stats.user_count || 0 }}</div>
            </div>

            <div class="status-row">
              <div class="row-label">
                <span class="label-icon">></span>
                TOTAL_PORTS
              </div>
              <div class="row-value">{{ stats.port_count || 0 }}</div>
            </div>

            <div class="status-row highlight">
              <div class="row-label">
                <span class="label-icon">></span>
                ACTIVE_PORTS
              </div>
              <div class="row-value active">{{ stats.active_port_count || 0 }}</div>
            </div>

            <div class="status-row highlight">
              <div class="row-label">
                <span class="label-icon">></span>
                ONLINE_DEVICES
              </div>
              <div class="row-value online">{{ stats.online_device_count || 0 }}</div>
            </div>

            <div class="status-row">
              <div class="row-label">
                <span class="label-icon">></span>
                PORT_RANGE
              </div>
              <div class="row-value">10000 - 60000</div>
            </div>

            <div class="status-row">
              <div class="row-label">
                <span class="label-icon">></span>
                SYSTEM_HEALTH
              </div>
              <div class="row-value health">
                <span class="health-indicator"></span>
                OPERATIONAL
              </div>
            </div>
          </div>

          <div class="status-actions">
            <router-link to="/admin/users" class="action-link">
              <span>查看用户列表</span>
              <el-icon><ArrowRight /></el-icon>
            </router-link>
            <router-link to="/admin/ports" class="action-link">
              <span>管理端口</span>
              <el-icon><ArrowRight /></el-icon>
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="quick-actions">
      <div class="actions-header">
        <el-icon><Operation /></el-icon>
        <span>QUICK_ACTIONS</span>
      </div>
      <div class="actions-grid">
        <router-link to="/admin/users" class="action-card">
          <div class="action-icon">
            <el-icon><User /></el-icon>
          </div>
          <span class="action-label">用户管理</span>
        </router-link>
        <router-link to="/admin/ports" class="action-card">
          <div class="action-icon">
            <el-icon><Connection /></el-icon>
          </div>
          <span class="action-label">端口管理</span>
        </router-link>
        <router-link to="/admin/devices" class="action-card">
          <div class="action-icon">
            <el-icon><Iphone /></el-icon>
          </div>
          <span class="action-label">设备管理</span>
        </router-link>
        <router-link to="/admin/orders" class="action-card">
          <div class="action-icon">
            <el-icon><Wallet /></el-icon>
          </div>
          <span class="action-label">订单管理</span>
        </router-link>
        <router-link to="/admin/config" class="action-card">
          <div class="action-icon">
            <el-icon><Tools /></el-icon>
          </div>
          <span class="action-label">系统配置</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import {
  Refresh, User, Connection, Iphone, Wallet, TrendCharts, Monitor,
  Top, Bottom, Calendar, ArrowRight, Operation, Tools
} from '@element-plus/icons-vue'
import api from '../../api'
import dayjs from 'dayjs'

const loading = ref(false)
const stats = ref({})
const lastUpdated = ref('')

const portUsagePercent = computed(() => {
  const total = stats.value.port_count || 1
  const active = stats.value.active_port_count || 0
  return Math.min((active / total) * 100, 100)
})

const netRevenue = computed(() => {
  const recharge = stats.value.total_recharge || 0
  const consume = stats.value.total_consume || 0
  return recharge - consume
})

function formatNumber(num) {
  if (num >= 10000) {
    return (num / 10000).toFixed(2) + 'w'
  }
  return num.toFixed(2)
}

async function fetchStats() {
  loading.value = true
  try {
    const res = await api.admin.getStats()
    if (res.code === 0) {
      stats.value = res.data
      lastUpdated.value = dayjs().format('HH:mm:ss')
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.admin-dashboard {
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
  background: linear-gradient(135deg, rgba(249, 115, 22, 0.2), rgba(249, 115, 22, 0.05));
  border: 1px solid rgba(249, 115, 22, 0.3);
  border-radius: var(--radius-lg);
  color: #f97316;
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

.page-subtitle {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
  margin: 4px 0 0 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.last-updated {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-mono);
  font-size: 11px;
}

.update-label {
  color: var(--text-muted);
}

.update-time {
  color: var(--neon-green);
}

.refresh-btn {
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 6px;
}

.refresh-btn:hover {
  border-color: #f97316;
  color: #f97316;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  position: relative;
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  overflow: hidden;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  border-color: var(--border-active);
}

.card-background {
  position: absolute;
  inset: 0;
  opacity: 0.03;
  pointer-events: none;
}

.hex-pattern {
  width: 100%;
  height: 100%;
  background-image: url("data:image/svg+xml,%3Csvg width='60' height='60' xmlns='http://www.w3.org/2000/svg'%3E%3Cpolygon points='30,0 60,15 60,45 30,60 0,45 0,15' fill='none' stroke='%23fff' stroke-width='1'/%3E%3C/svg%3E");
  background-size: 30px 30px;
}

.card-content {
  position: relative;
  padding: 24px;
}

.stat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.stat-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  font-size: 20px;
}

.stat-card.users .stat-icon {
  background: rgba(59, 130, 246, 0.15);
  color: #3b82f6;
}

.stat-card.ports .stat-icon {
  background: rgba(34, 197, 94, 0.15);
  color: var(--neon-green);
}

.stat-card.devices .stat-icon {
  background: rgba(168, 85, 247, 0.15);
  color: var(--neon-purple);
}

.stat-card.revenue .stat-icon {
  background: rgba(249, 115, 22, 0.15);
  color: #f97316;
}

.stat-code {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  padding: 2px 8px;
  background: var(--bg-elevated);
  border-radius: var(--radius-sm);
  letter-spacing: 1px;
}

.stat-value {
  font-family: var(--font-display);
  font-size: 36px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.stat-total {
  font-size: 16px;
  color: var(--text-muted);
}

.currency {
  font-size: 20px;
  color: var(--text-muted);
}

.stat-label {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 0.5px;
  margin-bottom: 16px;
}

.stat-footer {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-top: 16px;
  border-top: 1px solid var(--border-subtle);
}

.footer-text {
  font-size: 12px;
  color: var(--text-disabled);
}

.footer-percent {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--neon-green);
}

.progress-mini {
  flex: 1;
  height: 4px;
  background: var(--bg-elevated);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.progress-mini .progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--neon-green), var(--neon-cyan));
  border-radius: var(--radius-full);
  transition: width 0.5s ease;
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-indicator.online {
  background: var(--neon-green);
  box-shadow: 0 0 8px var(--neon-green);
  animation: pulse-status 2s ease-in-out infinite;
}

@keyframes pulse-status {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.card-glow {
  position: absolute;
  top: -50px;
  right: -50px;
  width: 150px;
  height: 150px;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.1;
  pointer-events: none;
}

.stat-card.users .card-glow {
  background: #3b82f6;
}

.stat-card.ports .card-glow {
  background: var(--neon-green);
}

.stat-card.devices .card-glow {
  background: var(--neon-purple);
}

.stat-card.revenue .card-glow {
  background: #f97316;
}

/* Dashboard Panels */
.dashboard-panels {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 24px;
}

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
  padding: 20px 24px;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border-subtle);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.panel-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  font-size: 16px;
}

.revenue-panel .panel-icon {
  background: rgba(34, 197, 94, 0.15);
  color: var(--neon-green);
}

.status-panel .panel-icon {
  background: rgba(0, 240, 255, 0.15);
  color: var(--neon-cyan);
}

.panel-title {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: 0.5px;
}

.panel-subtitle {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
}

.panel-content {
  padding: 24px;
}

/* Revenue Grid */
.revenue-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 20px;
}

.revenue-item {
  padding: 16px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
}

.item-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.item-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 0.5px;
}

.item-icon {
  font-size: 14px;
}

.total-recharge .item-icon,
.today-recharge .item-icon {
  color: var(--neon-green);
}

.total-consume .item-icon,
.today-consume .item-icon {
  color: #ef4444;
}

.item-value {
  font-family: var(--font-mono);
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 6px;
}

.item-value.positive {
  color: var(--neon-green);
}

.item-value.negative {
  color: #ef4444;
}

.item-desc {
  font-size: 11px;
  color: var(--text-disabled);
}

.revenue-summary {
  padding: 16px;
  background: rgba(249, 115, 22, 0.05);
  border: 1px solid rgba(249, 115, 22, 0.2);
  border-radius: var(--radius-md);
}

.summary-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.summary-label {
  font-family: var(--font-mono);
  font-size: 12px;
  color: #f97316;
}

.summary-value {
  font-family: var(--font-mono);
  font-size: 18px;
  font-weight: 700;
}

.summary-value.positive {
  color: var(--neon-green);
}

.summary-value.negative {
  color: #ef4444;
}

/* Status Grid */
.status-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 20px;
}

.status-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
}

.status-row.highlight {
  background: var(--bg-elevated);
  border-color: var(--border-default);
}

.row-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-secondary);
}

.label-icon {
  color: var(--neon-cyan);
}

.row-value {
  font-family: var(--font-mono);
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

.row-value.active {
  color: var(--neon-green);
}

.row-value.online {
  color: var(--neon-purple);
}

.row-value.health {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--neon-green);
}

.health-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--neon-green);
  box-shadow: 0 0 8px var(--neon-green);
}

.status-actions {
  display: flex;
  gap: 12px;
}

.action-link {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px;
  background: var(--bg-elevated);
  border: 1px dashed var(--border-default);
  border-radius: var(--radius-md);
  font-size: 12px;
  color: var(--text-secondary);
  text-decoration: none;
  transition: all 0.2s ease;
}

.action-link:hover {
  border-color: var(--neon-cyan);
  color: var(--neon-cyan);
  background: rgba(0, 240, 255, 0.05);
}

/* Quick Actions */
.quick-actions {
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  padding: 20px 24px;
}

.actions-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
  font-family: var(--font-mono);
  font-size: 12px;
  color: #f97316;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 16px;
}

.action-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 20px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  text-decoration: none;
  transition: all 0.2s ease;
}

.action-card:hover {
  background: var(--bg-elevated);
  border-color: #f97316;
  transform: translateY(-2px);
}

.action-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(249, 115, 22, 0.1);
  border-radius: var(--radius-md);
  color: #f97316;
  font-size: 22px;
}

.action-card:hover .action-icon {
  background: rgba(249, 115, 22, 0.2);
}

.action-label {
  font-size: 13px;
  color: var(--text-secondary);
}

.action-card:hover .action-label {
  color: var(--text-primary);
}

/* Responsive */
@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .actions-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .dashboard-panels {
    grid-template-columns: 1fr;
  }

  .actions-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .revenue-grid {
    grid-template-columns: 1fr;
  }
}
</style>
