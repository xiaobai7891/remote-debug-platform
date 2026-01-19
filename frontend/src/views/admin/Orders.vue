<template>
  <div class="admin-orders-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M21 4H3c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h18c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2z"/>
            <path d="M1 10h22"/>
          </svg>
        </div>
        <div class="header-text">
          <h1 class="page-title">
            <span class="title-prefix">//</span>
            ORDER_LEDGER
          </h1>
          <p class="page-subtitle">系统交易记录管理</p>
        </div>
      </div>
      <div class="header-actions">
        <div class="filter-group">
          <button
            v-for="filter in typeFilters"
            :key="filter.value"
            :class="['filter-btn', { active: typeFilter === filter.value }]"
            @click="typeFilter = filter.value; fetchOrders()"
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
        <span class="stat-label">TOTAL_ORDERS</span>
        <span class="stat-value">{{ total }}</span>
      </div>
      <div class="stat-chip recharge">
        <span class="stat-label">RECHARGE</span>
        <span class="stat-value">¥{{ totalRecharge.toFixed(2) }}</span>
      </div>
      <div class="stat-chip consume">
        <span class="stat-label">CONSUME</span>
        <span class="stat-value">¥{{ totalConsume.toFixed(2) }}</span>
      </div>
      <div class="stat-chip refund">
        <span class="stat-label">REFUND</span>
        <span class="stat-value">¥{{ totalRefund.toFixed(2) }}</span>
      </div>
    </div>

    <!-- Orders Table -->
    <div class="table-container">
      <div class="table-header">
        <div class="table-title">
          <span class="title-dot"></span>
          TRANSACTION_LOG
        </div>
        <div class="table-actions">
          <button class="action-btn refresh" @click="fetchOrders">
            <el-icon><Refresh /></el-icon>
            REFRESH
          </button>
        </div>
      </div>

      <el-table
        :data="orders"
        v-loading="loading"
        class="cyber-table"
        :header-cell-style="{ background: 'var(--bg-elevated)', color: 'var(--text-muted)', fontFamily: 'var(--font-mono)', fontSize: '11px', fontWeight: '600', letterSpacing: '1px', borderBottom: '1px solid var(--border-default)' }"
        :row-style="{ background: 'transparent' }"
      >
        <el-table-column prop="id" label="ORDER_ID" width="100">
          <template #default="{ row }">
            <span class="order-id">#{{ String(row.id).padStart(6, '0') }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="USER" width="140">
          <template #default="{ row }">
            <div class="user-cell">
              <div class="user-avatar-small">
                {{ row.username?.charAt(0).toUpperCase() || '?' }}
              </div>
              <span class="username">{{ row.username || '-' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="TYPE" width="120">
          <template #default="{ row }">
            <span :class="['type-badge', row.type]">
              <span class="type-icon">
                <el-icon v-if="row.type === 'recharge'"><Top /></el-icon>
                <el-icon v-else-if="row.type === 'consume'"><Bottom /></el-icon>
                <el-icon v-else><RefreshRight /></el-icon>
              </span>
              {{ typeNames[row.type] }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="AMOUNT" width="140">
          <template #default="{ row }">
            <span :class="['amount-value', getAmountClass(row.type)]">
              {{ row.type === 'consume' ? '-' : '+' }}¥{{ Math.abs(row.amount).toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="DESCRIPTION">
          <template #default="{ row }">
            <span class="description-text">{{ row.description || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="TIMESTAMP" width="180">
          <template #default="{ row }">
            <div class="time-cell">
              <span class="time-text">{{ formatTime(row.created_at) }}</span>
              <span class="time-relative">{{ formatRelative(row.created_at) }}</span>
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
          @size-change="fetchOrders"
          @current-change="fetchOrders"
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
import { Refresh, Top, Bottom, RefreshRight } from '@element-plus/icons-vue'
import api from '../../api'
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

const typeFilters = [
  { value: '', label: 'ALL', color: 'var(--text-muted)' },
  { value: 'recharge', label: 'RECHARGE', color: 'var(--neon-green)' },
  { value: 'consume', label: 'CONSUME', color: '#ef4444' },
  { value: 'refund', label: 'REFUND', color: '#eab308' }
]

const typeNames = {
  recharge: '充值',
  consume: '消费',
  refund: '退款'
}

// Computed stats
const totalRecharge = computed(() => {
  return orders.value
    .filter(o => o.type === 'recharge')
    .reduce((sum, o) => sum + Math.abs(o.amount || 0), 0)
})

const totalConsume = computed(() => {
  return orders.value
    .filter(o => o.type === 'consume')
    .reduce((sum, o) => sum + Math.abs(o.amount || 0), 0)
})

const totalRefund = computed(() => {
  return orders.value
    .filter(o => o.type === 'refund')
    .reduce((sum, o) => sum + Math.abs(o.amount || 0), 0)
})

function formatTime(time) {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm:ss') : '-'
}

function formatRelative(time) {
  return time ? dayjs(time).fromNow() : ''
}

function getAmountClass(type) {
  if (type === 'recharge') return 'positive'
  if (type === 'consume') return 'negative'
  return 'neutral'
}

async function fetchOrders() {
  loading.value = true
  try {
    const res = await api.admin.getOrders(page.value, pageSize.value, typeFilter.value)
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
.admin-orders-page {
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

.stat-chip.recharge .stat-value { color: var(--neon-green); }
.stat-chip.consume .stat-value { color: #ef4444; }
.stat-chip.refund .stat-value { color: #eab308; }

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
.order-id {
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 600;
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

.type-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 500;
}

.type-badge.recharge {
  background: rgba(34, 197, 94, 0.15);
  color: var(--neon-green);
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.type-badge.consume {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.type-badge.refund {
  background: rgba(234, 179, 8, 0.15);
  color: #eab308;
  border: 1px solid rgba(234, 179, 8, 0.3);
}

.type-icon {
  display: flex;
  align-items: center;
}

.amount-value {
  font-family: var(--font-mono);
  font-size: 15px;
  font-weight: 700;
}

.amount-value.positive {
  color: var(--neon-green);
}

.amount-value.negative {
  color: #ef4444;
}

.amount-value.neutral {
  color: #eab308;
}

.description-text {
  color: var(--text-secondary);
  font-size: 13px;
}

.time-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.time-text {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-secondary);
}

.time-relative {
  font-size: 11px;
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
