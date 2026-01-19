<template>
  <div class="orders-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-icon">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <path d="M9 5a2 2 0 012-2h2a2 2 0 012 2v0a2 2 0 01-2 2h-2a2 2 0 01-2-2v0z" stroke="currentColor" stroke-width="2"/>
            <path d="M9 12h6M9 16h6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </div>
        <div class="header-text">
          <h1 class="page-title">TRANSACTION_LOG</h1>
          <p class="page-subtitle">// 账户交易记录</p>
        </div>
      </div>
      <div class="header-actions">
        <el-button class="refresh-btn" @click="fetchOrders" :loading="loading">
          <el-icon><Refresh /></el-icon>
          <span>刷新</span>
        </el-button>
      </div>
    </div>

    <!-- Stats Bar -->
    <div class="stats-bar">
      <div class="stat-item">
        <div class="stat-icon recharge">
          <el-icon><Top /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-label">TOTAL_RECHARGE</span>
          <span class="stat-value recharge">+¥{{ totalRecharge.toFixed(2) }}</span>
        </div>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <div class="stat-icon consume">
          <el-icon><Bottom /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-label">TOTAL_CONSUME</span>
          <span class="stat-value consume">-¥{{ totalConsume.toFixed(2) }}</span>
        </div>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <div class="stat-icon refund">
          <el-icon><RefreshLeft /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-label">TOTAL_REFUND</span>
          <span class="stat-value refund">+¥{{ totalRefund.toFixed(2) }}</span>
        </div>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <div class="stat-icon total">
          <el-icon><Wallet /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-label">RECORDS</span>
          <span class="stat-value">{{ total }} 条</span>
        </div>
      </div>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar">
      <div class="filter-group">
        <span class="filter-label">> TYPE_FILTER:</span>
        <div class="filter-chips">
          <button
            class="filter-chip"
            :class="{ active: typeFilter === '' }"
            @click="typeFilter = ''; fetchOrders()"
          >
            全部
          </button>
          <button
            class="filter-chip recharge"
            :class="{ active: typeFilter === 'recharge' }"
            @click="typeFilter = 'recharge'; fetchOrders()"
          >
            <span class="chip-dot"></span>
            充值
          </button>
          <button
            class="filter-chip consume"
            :class="{ active: typeFilter === 'consume' }"
            @click="typeFilter = 'consume'; fetchOrders()"
          >
            <span class="chip-dot"></span>
            消费
          </button>
          <button
            class="filter-chip refund"
            :class="{ active: typeFilter === 'refund' }"
            @click="typeFilter = 'refund'; fetchOrders()"
          >
            <span class="chip-dot"></span>
            退款
          </button>
        </div>
      </div>
    </div>

    <!-- Transaction List -->
    <div class="transaction-container" v-loading="loading">
      <div v-if="orders.length === 0 && !loading" class="empty-state">
        <div class="empty-icon">
          <svg viewBox="0 0 80 80" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="10" y="15" width="60" height="50" rx="4" stroke="currentColor" stroke-width="2" fill="none"/>
            <path d="M10 30h60" stroke="currentColor" stroke-width="2"/>
            <circle cx="20" cy="22.5" r="2.5" fill="currentColor" opacity="0.5"/>
            <circle cx="28" cy="22.5" r="2.5" fill="currentColor" opacity="0.5"/>
            <circle cx="36" cy="22.5" r="2.5" fill="currentColor" opacity="0.5"/>
            <path d="M25 45h30M25 52h20" stroke="currentColor" stroke-width="2" stroke-linecap="round" opacity="0.3"/>
          </svg>
        </div>
        <p class="empty-text">NO_TRANSACTIONS_FOUND</p>
        <p class="empty-hint">暂无交易记录</p>
      </div>

      <div v-else class="transaction-list">
        <div
          v-for="order in orders"
          :key="order.id"
          class="transaction-item"
          :class="order.type"
        >
          <div class="transaction-indicator">
            <div class="indicator-line"></div>
            <div class="indicator-dot" :class="order.type"></div>
          </div>

          <div class="transaction-content">
            <div class="transaction-header">
              <div class="transaction-type">
                <span class="type-badge" :class="order.type">
                  <span class="badge-icon">
                    <el-icon v-if="order.type === 'recharge'"><Top /></el-icon>
                    <el-icon v-else-if="order.type === 'consume'"><Bottom /></el-icon>
                    <el-icon v-else><RefreshLeft /></el-icon>
                  </span>
                  {{ typeNames[order.type] }}
                </span>
                <span class="transaction-id">#{{ order.id.toString().padStart(6, '0') }}</span>
              </div>
              <div class="transaction-amount" :class="order.type">
                <span class="amount-sign">{{ order.amount >= 0 ? '+' : '' }}</span>
                <span class="amount-value">¥{{ Math.abs(order.amount).toFixed(2) }}</span>
              </div>
            </div>

            <div class="transaction-body">
              <p class="transaction-desc">{{ order.description || '无描述' }}</p>
            </div>

            <div class="transaction-footer">
              <div class="transaction-time">
                <el-icon><Clock /></el-icon>
                <span>{{ formatTime(order.created_at) }}</span>
              </div>
              <div class="transaction-relative">
                {{ formatRelativeTime(order.created_at) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div class="pagination-container" v-if="total > 0">
        <div class="pagination-info">
          <span class="info-label">> SHOWING:</span>
          <span class="info-value">{{ (page - 1) * pageSize + 1 }}-{{ Math.min(page * pageSize, total) }}</span>
          <span class="info-label">OF</span>
          <span class="info-value">{{ total }}</span>
        </div>
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="sizes, prev, pager, next"
          @size-change="fetchOrders"
          @current-change="fetchOrders"
          class="cyber-pagination"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Refresh, Top, Bottom, RefreshLeft, Wallet, Clock } from '@element-plus/icons-vue'
import api from '../api'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const loading = ref(false)
const orders = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const typeFilter = ref('')

const typeNames = {
  recharge: '充值',
  consume: '消费',
  refund: '退款'
}

// Computed stats
const totalRecharge = computed(() => {
  return orders.value
    .filter(o => o.type === 'recharge')
    .reduce((sum, o) => sum + Math.abs(o.amount), 0)
})

const totalConsume = computed(() => {
  return orders.value
    .filter(o => o.type === 'consume')
    .reduce((sum, o) => sum + Math.abs(o.amount), 0)
})

const totalRefund = computed(() => {
  return orders.value
    .filter(o => o.type === 'refund')
    .reduce((sum, o) => sum + Math.abs(o.amount), 0)
})

function formatTime(time) {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm:ss') : '-'
}

function formatRelativeTime(time) {
  return time ? dayjs(time).fromNow() : '-'
}

async function fetchOrders() {
  loading.value = true
  try {
    const res = await api.user.getOrders(page.value, pageSize.value, typeFilter.value)
    if (res.code === 0) {
      orders.value = res.data.list || []
      total.value = res.data.total
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchOrders()
})
</script>

<style scoped>
.orders-page {
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
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.2), rgba(34, 197, 94, 0.05));
  border: 1px solid rgba(34, 197, 94, 0.3);
  border-radius: var(--radius-lg);
  color: var(--neon-green);
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

.refresh-btn {
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 6px;
}

.refresh-btn:hover {
  border-color: var(--neon-green);
  color: var(--neon-green);
}

/* Stats Bar */
.stats-bar {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 20px 24px;
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  margin-bottom: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.stat-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  font-size: 18px;
}

.stat-icon.recharge {
  background: rgba(34, 197, 94, 0.15);
  color: var(--neon-green);
}

.stat-icon.consume {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.stat-icon.refund {
  background: rgba(251, 191, 36, 0.15);
  color: #fbbf24;
}

.stat-icon.total {
  background: rgba(0, 240, 255, 0.15);
  color: var(--neon-cyan);
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 0.5px;
}

.stat-value {
  font-family: var(--font-mono);
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.stat-value.recharge {
  color: var(--neon-green);
}

.stat-value.consume {
  color: #ef4444;
}

.stat-value.refund {
  color: #fbbf24;
}

.stat-divider {
  width: 1px;
  height: 40px;
  background: var(--border-default);
}

/* Filter Bar */
.filter-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  margin-bottom: 20px;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-label {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--neon-green);
}

.filter-chips {
  display: flex;
  gap: 8px;
}

.filter-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-full);
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.filter-chip:hover {
  border-color: var(--text-muted);
}

.filter-chip.active {
  background: rgba(0, 240, 255, 0.1);
  border-color: var(--neon-cyan);
  color: var(--neon-cyan);
}

.filter-chip.recharge.active {
  background: rgba(34, 197, 94, 0.1);
  border-color: var(--neon-green);
  color: var(--neon-green);
}

.filter-chip.consume.active {
  background: rgba(239, 68, 68, 0.1);
  border-color: #ef4444;
  color: #ef4444;
}

.filter-chip.refund.active {
  background: rgba(251, 191, 36, 0.1);
  border-color: #fbbf24;
  color: #fbbf24;
}

.chip-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}

/* Transaction Container */
.transaction-container {
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  min-height: 400px;
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
  margin-bottom: 20px;
  opacity: 0.5;
}

.empty-text {
  font-family: var(--font-mono);
  font-size: 14px;
  color: var(--text-muted);
  margin: 0 0 8px 0;
  letter-spacing: 1px;
}

.empty-hint {
  font-size: 13px;
  color: var(--text-disabled);
  margin: 0;
}

/* Transaction List */
.transaction-list {
  padding: 20px;
}

.transaction-item {
  display: flex;
  gap: 16px;
  padding: 16px 0;
  border-bottom: 1px solid var(--border-subtle);
  transition: background 0.2s ease;
}

.transaction-item:last-child {
  border-bottom: none;
}

.transaction-item:hover {
  background: rgba(255, 255, 255, 0.02);
}

/* Transaction Indicator */
.transaction-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 20px;
  padding-top: 4px;
}

.indicator-line {
  width: 2px;
  flex: 1;
  background: var(--border-subtle);
  margin-top: 8px;
}

.transaction-item:last-child .indicator-line {
  display: none;
}

.indicator-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 2px solid;
  background: var(--bg-card);
}

.indicator-dot.recharge {
  border-color: var(--neon-green);
  box-shadow: 0 0 8px rgba(34, 197, 94, 0.4);
}

.indicator-dot.consume {
  border-color: #ef4444;
  box-shadow: 0 0 8px rgba(239, 68, 68, 0.4);
}

.indicator-dot.refund {
  border-color: #fbbf24;
  box-shadow: 0 0 8px rgba(251, 191, 36, 0.4);
}

/* Transaction Content */
.transaction-content {
  flex: 1;
  min-width: 0;
}

.transaction-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.transaction-type {
  display: flex;
  align-items: center;
  gap: 12px;
}

.type-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: var(--radius-full);
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 500;
}

.type-badge.recharge {
  background: rgba(34, 197, 94, 0.15);
  color: var(--neon-green);
}

.type-badge.consume {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.type-badge.refund {
  background: rgba(251, 191, 36, 0.15);
  color: #fbbf24;
}

.badge-icon {
  display: flex;
  font-size: 12px;
}

.transaction-id {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
}

.transaction-amount {
  font-family: var(--font-mono);
  font-size: 18px;
  font-weight: 700;
}

.transaction-amount.recharge {
  color: var(--neon-green);
}

.transaction-amount.consume {
  color: #ef4444;
}

.transaction-amount.refund {
  color: #fbbf24;
}

.amount-sign {
  font-size: 14px;
  margin-right: 2px;
}

.transaction-body {
  margin-bottom: 10px;
}

.transaction-desc {
  font-size: 13px;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.5;
}

.transaction-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.transaction-time {
  display: flex;
  align-items: center;
  gap: 6px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
}

.transaction-relative {
  font-size: 11px;
  color: var(--text-disabled);
}

/* Pagination */
.pagination-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-top: 1px solid var(--border-subtle);
}

.pagination-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-family: var(--font-mono);
  font-size: 12px;
}

.info-label {
  color: var(--text-muted);
}

.info-value {
  color: var(--neon-cyan);
}

.cyber-pagination :deep(.el-pagination) {
  --el-pagination-bg-color: transparent;
  --el-pagination-button-bg-color: var(--bg-elevated);
  --el-pagination-hover-color: var(--neon-cyan);
}

.cyber-pagination :deep(.el-pager li) {
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  margin: 0 2px;
}

.cyber-pagination :deep(.el-pager li:hover) {
  color: var(--neon-cyan);
  border-color: var(--neon-cyan);
}

.cyber-pagination :deep(.el-pager li.is-active) {
  background: rgba(0, 240, 255, 0.15);
  border-color: var(--neon-cyan);
  color: var(--neon-cyan);
}

.cyber-pagination :deep(.btn-prev),
.cyber-pagination :deep(.btn-next) {
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
}

.cyber-pagination :deep(.btn-prev:hover),
.cyber-pagination :deep(.btn-next:hover) {
  color: var(--neon-cyan);
  border-color: var(--neon-cyan);
}

/* Responsive */
@media (max-width: 768px) {
  .stats-bar {
    flex-wrap: wrap;
    gap: 16px;
  }

  .stat-item {
    flex: 1 1 45%;
  }

  .stat-divider {
    display: none;
  }

  .filter-bar {
    flex-direction: column;
    align-items: flex-start;
  }

  .filter-chips {
    flex-wrap: wrap;
  }

  .transaction-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .pagination-container {
    flex-direction: column;
    gap: 16px;
  }
}
</style>
