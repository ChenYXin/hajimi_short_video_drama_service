<template>
  <div class="users">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="8">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-number">{{ userStats.active }}</div>
              <div class="stat-label">活跃用户</div>
            </div>
            <div class="stat-icon green">
              <el-icon size="20"><User /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="8">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-number">{{ userStats.inactive }}</div>
              <div class="stat-label">非活跃用户</div>
            </div>
            <div class="stat-icon yellow">
              <el-icon size="20"><Warning /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="8">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-number">{{ userStats.total }}</div>
              <div class="stat-label">总用户数</div>
            </div>
            <div class="stat-icon blue">
              <el-icon size="20"><UserFilled /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
        </div>
      </template>

      <!-- 搜索筛选 -->
      <div class="search-bar">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-input
              v-model="searchForm.keyword"
              placeholder="搜索用户名或邮箱"
              clearable
              @input="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </el-col>
          <el-col :span="4">
            <el-select
              v-model="searchForm.status"
              placeholder="状态"
              clearable
              @change="handleSearch"
            >
              <el-option label="活跃" value="active" />
              <el-option label="非活跃" value="inactive" />
              <el-option label="已封禁" value="banned" />
            </el-select>
          </el-col>
          <el-col :span="6">
            <el-date-picker
              v-model="searchForm.dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              @change="handleSearch"
            />
          </el-col>
        </el-row>
      </div>

      <!-- 数据表格 -->
      <el-table
        :data="tableData"
        :loading="loading"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户" min-width="200">
          <template #default="{ row }">
            <div class="user-info">
              <el-avatar :size="40" :src="row.avatar">
                {{ row.username.charAt(0).toUpperCase() }}
              </el-avatar>
              <div class="user-details">
                <div class="username">{{ row.username }}</div>
                <div class="email">{{ row.email }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="手机号" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="getStatusType(row.status)"
              size="small"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_login_at" label="最后登录" width="180" />
        <el-table-column prop="created_at" label="注册时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'active'"
              size="small"
              type="warning"
              @click="handleDeactivate(row)"
            >
              禁用
            </el-button>
            <el-button
              v-else
              size="small"
              type="success"
              @click="handleActivate(row)"
            >
              启用
            </el-button>
            <el-button
              size="small"
              type="danger"
              @click="handleBan(row)"
            >
              封禁
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { User, UserFilled, Warning, Search } from '@element-plus/icons-vue'
import api from '@/utils/api'

const loading = ref(false)

const searchForm = reactive({
  keyword: '',
  status: '',
  dateRange: []
})

const userStats = ref({
  active: 0,
  inactive: 0,
  total: 0
})

const tableData = ref([])
const pagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

const getStatusType = (status) => {
  const typeMap = {
    active: 'success',
    inactive: 'warning',
    banned: 'danger'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status) => {
  const textMap = {
    active: '活跃',
    inactive: '非活跃',
    banned: '已封禁'
  }
  return textMap[status] || status
}

const loadStats = async () => {
  try {
    // 模拟数据
    userStats.value = {
      active: 7856,
      inactive: 1076,
      total: 8932
    }
  } catch (error) {
    ElMessage.error('加载统计数据失败')
  }
}

const loadData = async () => {
  loading.value = true
  try {
    // 模拟数据
    tableData.value = [
      {
        id: 1,
        username: 'user001',
        email: 'user001@example.com',
        phone: '13800138001',
        avatar: '',
        status: 'active',
        last_login_at: '2024-01-15 14:30:00',
        created_at: '2024-01-10 09:15:00'
      },
      {
        id: 2,
        username: 'user002',
        email: 'user002@example.com',
        phone: '13800138002',
        avatar: '',
        status: 'inactive',
        last_login_at: '2024-01-12 16:20:00',
        created_at: '2024-01-08 11:30:00'
      }
    ]
    pagination.total = 2
  } catch (error) {
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleActivate = async (row) => {
  try {
    await ElMessageBox.confirm('确定要启用这个用户吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })
    ElMessage.success('用户已启用')
    loadData()
    loadStats()
  } catch {
    // 用户取消
  }
}

const handleDeactivate = async (row) => {
  try {
    await ElMessageBox.confirm('确定要禁用这个用户吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    ElMessage.success('用户已禁用')
    loadData()
    loadStats()
  } catch {
    // 用户取消
  }
}

const handleBan = async (row) => {
  try {
    await ElMessageBox.confirm('确定要封禁这个用户吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'error'
    })
    ElMessage.success('用户已封禁')
    loadData()
    loadStats()
  } catch {
    // 用户取消
  }
}

onMounted(() => {
  loadStats()
  loadData()
})
</script>

<style scoped>
.users {
  max-width: 1200px;
  margin: 0 auto;
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
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-icon.green {
  background-color: #67c23a;
}

.stat-icon.yellow {
  background-color: #e6a23c;
}

.stat-icon.blue {
  background-color: #409eff;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-bar {
  margin-bottom: 20px;
}

.user-info {
  display: flex;
  align-items: center;
}

.user-details {
  margin-left: 12px;
}

.username {
  font-weight: 500;
  color: #303133;
}

.email {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}
</style>
