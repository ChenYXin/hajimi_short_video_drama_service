<template>
  <div class="dashboard">
    <!-- 欢迎区域 -->
    <el-card class="welcome-card" shadow="never">
      <div class="welcome-content">
        <div class="welcome-text">
          <h2>欢迎回来！</h2>
          <p>今天是个美好的一天，让我们开始工作吧</p>
        </div>
        <div class="welcome-icon">
          <el-icon size="48"><TrendCharts /></el-icon>
        </div>
      </div>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-number">{{ stats.totalDramas }}</div>
              <div class="stat-label">总短剧数</div>
              <div class="stat-trend">
                <el-icon color="#67c23a"><TrendCharts /></el-icon>
                <span>+12% 本月</span>
              </div>
            </div>
            <div class="stat-icon blue">
              <el-icon size="24"><VideoPlay /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-number">{{ stats.totalEpisodes }}</div>
              <div class="stat-label">总剧集数</div>
              <div class="stat-trend">
                <el-icon color="#67c23a"><TrendCharts /></el-icon>
                <span>+8% 本月</span>
              </div>
            </div>
            <div class="stat-icon green">
              <el-icon size="24"><Collection /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-number">{{ stats.totalUsers }}</div>
              <div class="stat-label">总用户数</div>
              <div class="stat-trend">
                <el-icon color="#67c23a"><TrendCharts /></el-icon>
                <span>+15% 本月</span>
              </div>
            </div>
            <div class="stat-icon purple">
              <el-icon size="24"><User /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-number">{{ stats.totalViews }}</div>
              <div class="stat-label">总观看数</div>
              <div class="stat-trend">
                <el-icon color="#67c23a"><TrendCharts /></el-icon>
                <span>+23% 昨日</span>
              </div>
            </div>
            <div class="stat-icon orange">
              <el-icon size="24"><View /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="charts-row">
      <el-col :xs="24" :lg="12">
        <el-card title="访问趋势" shadow="never">
          <div class="chart-placeholder">
            <p>图表区域 - 可集成 ECharts</p>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :lg="12">
        <el-card title="用户分布" shadow="never">
          <div class="chart-placeholder">
            <p>图表区域 - 可集成 ECharts</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { TrendCharts, VideoPlay, Collection, User, View } from '@element-plus/icons-vue'
import api from '@/utils/api'

const stats = ref({
  totalDramas: 0,
  totalEpisodes: 0,
  totalUsers: 0,
  totalViews: 0
})

const loadStats = async () => {
  try {
    // 这里应该调用实际的API
    // const response = await api.get('/admin/api/stats')
    // stats.value = response.data.data
    
    // 模拟数据
    stats.value = {
      totalDramas: 156,
      totalEpisodes: 1248,
      totalUsers: 8932,
      totalViews: 125678
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.dashboard {
  max-width: 1200px;
  margin: 0 auto;
}

.welcome-card {
  margin-bottom: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.welcome-card :deep(.el-card__body) {
  padding: 30px;
}

.welcome-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.welcome-text h2 {
  font-size: 28px;
  margin-bottom: 8px;
}

.welcome-text p {
  font-size: 16px;
  opacity: 0.9;
}

.welcome-icon {
  opacity: 0.3;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  margin-bottom: 20px;
}

.stat-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-number {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-trend {
  display: flex;
  align-items: center;
  font-size: 12px;
  color: #67c23a;
}

.stat-trend span {
  margin-left: 4px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-icon.blue {
  background-color: #409eff;
}

.stat-icon.green {
  background-color: #67c23a;
}

.stat-icon.purple {
  background-color: #9c27b0;
}

.stat-icon.orange {
  background-color: #ff9800;
}

.charts-row {
  margin-bottom: 20px;
}

.chart-placeholder {
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  border-radius: 4px;
  color: #909399;
}
</style>
