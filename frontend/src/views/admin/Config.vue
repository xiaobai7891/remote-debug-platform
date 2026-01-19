<template>
  <div class="admin-config-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="12" cy="12" r="3"/>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
          </svg>
        </div>
        <div class="header-text">
          <h1 class="page-title">
            <span class="title-prefix">//</span>
            SYSTEM_CONFIG
          </h1>
          <p class="page-subtitle">系统参数配置</p>
        </div>
      </div>
    </div>

    <!-- Config Sections -->
    <div class="config-grid">
      <!-- Port Pool Config -->
      <div class="config-section">
        <div class="section-header">
          <div class="section-icon port">
            <el-icon><Connection /></el-icon>
          </div>
          <div class="section-title">
            <h3>PORT_POOL</h3>
            <p>端口池配置</p>
          </div>
        </div>
        <div class="section-content">
          <div class="config-item">
            <div class="item-label">
              <span class="label-icon">></span>
              MIN_PORT
            </div>
            <div class="item-input">
              <el-input-number
                v-model="config.port_pool_min"
                :min="1024"
                :max="65535"
                class="cyber-number"
              />
            </div>
            <div class="item-hint">端口池起始值 (1024-65535)</div>
          </div>
          <div class="config-item">
            <div class="item-label">
              <span class="label-icon">></span>
              MAX_PORT
            </div>
            <div class="item-input">
              <el-input-number
                v-model="config.port_pool_max"
                :min="1024"
                :max="65535"
                class="cyber-number"
              />
            </div>
            <div class="item-hint">端口池结束值 (1024-65535)</div>
          </div>
          <div class="config-preview">
            <span class="preview-label">可用端口范围</span>
            <span class="preview-value">
              :{{ config.port_pool_min }} - :{{ config.port_pool_max }}
              <span class="preview-count">({{ config.port_pool_max - config.port_pool_min + 1 }} 个)</span>
            </span>
          </div>
        </div>
        <div class="corner-decor top-left"></div>
        <div class="corner-decor top-right"></div>
        <div class="corner-decor bottom-left"></div>
        <div class="corner-decor bottom-right"></div>
      </div>

      <!-- Pricing Config -->
      <div class="config-section">
        <div class="section-header">
          <div class="section-icon pricing">
            <el-icon><Wallet /></el-icon>
          </div>
          <div class="section-title">
            <h3>PRICING</h3>
            <p>计费配置</p>
          </div>
        </div>
        <div class="section-content">
          <div class="config-item">
            <div class="item-label">
              <span class="label-icon">></span>
              PORT_PER_DAY
            </div>
            <div class="item-input price">
              <span class="currency">¥</span>
              <el-input-number
                v-model="config.port_per_day"
                :min="0"
                :precision="2"
                :controls="false"
                class="cyber-number price"
              />
              <span class="unit">/天</span>
            </div>
            <div class="item-hint">每端口每日费用</div>
          </div>
          <div class="pricing-preview">
            <div class="price-card">
              <span class="price-period">日费</span>
              <span class="price-value">¥{{ config.port_per_day.toFixed(2) }}</span>
            </div>
            <div class="price-card">
              <span class="price-period">周费</span>
              <span class="price-value">¥{{ (config.port_per_day * 7).toFixed(2) }}</span>
            </div>
            <div class="price-card">
              <span class="price-period">月费</span>
              <span class="price-value">¥{{ (config.port_per_day * 30).toFixed(2) }}</span>
            </div>
          </div>
        </div>
        <div class="corner-decor top-left"></div>
        <div class="corner-decor top-right"></div>
        <div class="corner-decor bottom-left"></div>
        <div class="corner-decor bottom-right"></div>
      </div>
    </div>

    <!-- Action Bar -->
    <div class="action-bar">
      <div class="action-info">
        <el-icon><Warning /></el-icon>
        <span>配置修改仅在内存中生效，重启服务后会恢复为配置文件中的值。如需永久生效，请修改 config.yaml 文件。</span>
      </div>
      <div class="action-buttons">
        <button class="btn-reset" @click="fetchConfig">
          <el-icon><Refresh /></el-icon>
          重置
        </button>
        <button class="btn-save" @click="handleSave" :disabled="saveLoading">
          <el-icon v-if="saveLoading" class="is-loading"><Loading /></el-icon>
          <el-icon v-else><Check /></el-icon>
          {{ saveLoading ? '保存中...' : '保存配置' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Connection, Wallet, Warning, Refresh, Loading, Check } from '@element-plus/icons-vue'
import api from '../../api'

const config = reactive({
  port_pool_min: 10000,
  port_pool_max: 60000,
  port_per_day: 1.0
})

const saveLoading = ref(false)

async function fetchConfig() {
  const res = await api.admin.getConfig()
  if (res.code === 0) {
    config.port_pool_min = res.data.port_pool?.min || 10000
    config.port_pool_max = res.data.port_pool?.max || 60000
    config.port_per_day = res.data.pricing?.port_per_day || 1.0
  }
}

async function handleSave() {
  saveLoading.value = true
  try {
    const res = await api.admin.updateConfig({
      port_pool_min: config.port_pool_min,
      port_pool_max: config.port_pool_max,
      port_per_day: config.port_per_day
    })
    if (res.code === 0) {
      ElMessage.success('配置已保存')
    } else {
      ElMessage.error(res.message)
    }
  } finally {
    saveLoading.value = false
  }
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped>
.admin-config-page {
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

/* Config Grid */
.config-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

@media (max-width: 1200px) {
  .config-grid {
    grid-template-columns: 1fr;
  }
}

/* Config Section */
.config-section {
  position: relative;
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-xl);
  padding: 24px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border-subtle);
}

.section-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.section-icon.port {
  background: linear-gradient(135deg, rgba(0, 240, 255, 0.2), rgba(0, 240, 255, 0.05));
  border: 1px solid rgba(0, 240, 255, 0.3);
  color: var(--neon-cyan);
}

.section-icon.pricing {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.2), rgba(34, 197, 94, 0.05));
  border: 1px solid rgba(34, 197, 94, 0.3);
  color: var(--neon-green);
}

.section-title h3 {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: 1px;
}

.section-title p {
  font-size: 12px;
  color: var(--text-muted);
  margin: 4px 0 0;
}

.section-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* Config Item */
.config-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.item-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 500;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.label-icon {
  color: var(--neon-orange);
}

.item-input {
  display: flex;
  align-items: center;
}

.item-input.price {
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  padding: 0 12px;
}

.item-input.price:focus-within {
  border-color: var(--neon-green);
  box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.1);
}

.currency {
  font-family: var(--font-mono);
  font-size: 18px;
  font-weight: 700;
  color: var(--neon-green);
}

.unit {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
  margin-left: 4px;
}

.cyber-number {
  width: 100%;
}

.cyber-number :deep(.el-input__wrapper) {
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  padding: 10px 14px;
}

.cyber-number :deep(.el-input__wrapper:focus-within) {
  border-color: var(--neon-cyan);
  box-shadow: 0 0 0 3px rgba(0, 240, 255, 0.1);
}

.cyber-number.price :deep(.el-input__wrapper) {
  background: transparent !important;
  border: none !important;
  box-shadow: none !important;
  padding: 8px 0;
}

.cyber-number.price :deep(.el-input__inner) {
  font-family: var(--font-mono);
  font-size: 24px;
  font-weight: 700;
  text-align: center;
}

.item-hint {
  font-size: 11px;
  color: var(--text-muted);
}

/* Config Preview */
.config-preview {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  background: rgba(0, 240, 255, 0.1);
  border: 1px solid rgba(0, 240, 255, 0.2);
  border-radius: var(--radius-md);
  margin-top: 8px;
}

.preview-label {
  font-size: 12px;
  color: var(--text-secondary);
}

.preview-value {
  font-family: var(--font-mono);
  font-size: 14px;
  font-weight: 600;
  color: var(--neon-cyan);
}

.preview-count {
  font-size: 11px;
  color: var(--text-muted);
  margin-left: 8px;
}

/* Pricing Preview */
.pricing-preview {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-top: 8px;
}

.price-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 16px 12px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  transition: all 0.2s ease;
}

.price-card:hover {
  border-color: var(--neon-green);
  background: rgba(34, 197, 94, 0.05);
}

.price-period {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 1px;
  text-transform: uppercase;
}

.price-value {
  font-family: var(--font-mono);
  font-size: 18px;
  font-weight: 700;
  color: var(--neon-green);
}

/* Corner Decorations */
.corner-decor {
  position: absolute;
  width: 14px;
  height: 14px;
  border-color: var(--neon-orange);
  opacity: 0.3;
}

.corner-decor.top-left { top: 8px; left: 8px; border-top: 2px solid; border-left: 2px solid; }
.corner-decor.top-right { top: 8px; right: 8px; border-top: 2px solid; border-right: 2px solid; }
.corner-decor.bottom-left { bottom: 8px; left: 8px; border-bottom: 2px solid; border-left: 2px solid; }
.corner-decor.bottom-right { bottom: 8px; right: 8px; border-bottom: 2px solid; border-right: 2px solid; }

/* Action Bar */
.action-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  background: var(--bg-card);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-xl);
}

.action-info {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: var(--text-muted);
}

.action-info .el-icon {
  color: #eab308;
  font-size: 18px;
}

.action-buttons {
  display: flex;
  gap: 12px;
}

.btn-reset {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-reset:hover {
  border-color: var(--text-muted);
  color: var(--text-primary);
}

.btn-save {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 24px;
  background: linear-gradient(135deg, var(--neon-orange), #ea580c);
  border: none;
  border-radius: var(--radius-md);
  color: white;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-save:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.btn-save:disabled {
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
