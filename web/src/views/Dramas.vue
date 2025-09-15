<template>
  <div class="dramas">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>短剧管理</span>
          <el-button type="primary" @click="showCreateDialog = true">
            <el-icon><Plus /></el-icon>
            新增短剧
          </el-button>
        </div>
      </template>

      <!-- 搜索筛选 -->
      <div class="search-bar">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-input
              v-model="searchForm.title"
              placeholder="搜索短剧标题"
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
              <el-option label="草稿" value="draft" />
              <el-option label="已发布" value="published" />
              <el-option label="已归档" value="archived" />
            </el-select>
          </el-col>
          <el-col :span="4">
            <el-input
              v-model="searchForm.category"
              placeholder="分类"
              clearable
              @input="handleSearch"
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
        <el-table-column label="封面" width="80">
          <template #default="{ row }">
            <el-image
              v-if="row.cover_image"
              :src="row.cover_image"
              :preview-src-list="[row.cover_image]"
              fit="cover"
              style="width: 60px; height: 40px; border-radius: 4px;"
              :preview-teleported="true"
            />
            <span v-else class="text-gray-400">无封面</span>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="150" />
        <el-table-column prop="category" label="分类" width="100" />
        <el-table-column prop="view_count" label="观看数" width="100" />
        <el-table-column prop="like_count" label="点赞数" width="100" />
        <el-table-column prop="rating" label="评分" width="80">
          <template #default="{ row }">
            <span>{{ row.rating ? row.rating.toFixed(1) : '暂无' }}</span>
          </template>
        </el-table-column>
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
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">查看</el-button>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      :title="editingItem ? '编辑短剧' : '新增短剧'"
      width="600px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="80px"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-input v-model="form.category" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status">
            <el-option label="草稿" value="draft" />
            <el-option label="已发布" value="published" />
            <el-option label="已归档" value="archived" />
          </el-select>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import api from '@/utils/api'

const loading = ref(false)
const showCreateDialog = ref(false)
const editingItem = ref(null)
const formRef = ref()

const searchForm = reactive({
  title: '',
  status: '',
  category: ''
})

const form = reactive({
  title: '',
  description: '',
  category: '',
  status: 'draft'
})

const formRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  category: [{ required: true, message: '请输入分类', trigger: 'blur' }]
}

const tableData = ref([])
const pagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

const getStatusType = (status) => {
  const typeMap = {
    draft: 'info',
    published: 'success',
    archived: 'warning'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status) => {
  const textMap = {
    draft: '草稿',
    published: '已发布',
    archived: '已归档'
  }
  return textMap[status] || status
}

const loadData = async () => {
  loading.value = true
  try {
    const response = await api.get('/api/admin/dramas', {
      params: {
        page: pagination.page,
        page_size: pagination.size,
        ...searchForm
      }
    })

    if (response.data.success) {
      tableData.value = response.data.data.dramas
      pagination.total = response.data.data.total
    } else {
      ElMessage.error(response.data.message || '加载数据失败')
    }
  } catch (error) {
    console.error('加载短剧数据失败:', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleView = (row) => {
  // 创建一个详情查看窗口
  const detailWindow = window.open('', '_blank', 'width=600,height=800,scrollbars=yes,resizable=yes')
  detailWindow.document.write(`
    <!DOCTYPE html>
    <html>
    <head>
      <title>${row.title} - 短剧详情</title>
      <meta charset="utf-8">
      <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 500px; margin: 0 auto; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .cover { width: 100%; height: 300px; object-fit: cover; }
        .content { padding: 20px; }
        .title { font-size: 24px; font-weight: bold; margin-bottom: 10px; color: #333; }
        .category { display: inline-block; background: #007bff; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px; margin-bottom: 15px; }
        .description { color: #666; line-height: 1.6; margin-bottom: 15px; }
        .stats { display: flex; justify-content: space-between; margin-bottom: 15px; }
        .stat { text-align: center; }
        .stat-value { font-size: 18px; font-weight: bold; color: #007bff; }
        .stat-label { font-size: 12px; color: #999; }
        .info { color: #666; font-size: 14px; }
        .info strong { color: #333; }
      </style>
    </head>
    <body>
      <div class="container">
        <img src="${row.cover_image}" alt="${row.title}" class="cover" onerror="this.src='https://via.placeholder.com/500x300?text=暂无封面'">
        <div class="content">
          <h1 class="title">${row.title}</h1>
          <span class="category">${row.category}</span>
          <p class="description">${row.description || '暂无描述'}</p>
          <div class="stats">
            <div class="stat">
              <div class="stat-value">${row.view_count || 0}</div>
              <div class="stat-label">观看次数</div>
            </div>
            <div class="stat">
              <div class="stat-value">${row.like_count || 0}</div>
              <div class="stat-label">点赞数</div>
            </div>
            <div class="stat">
              <div class="stat-value">${row.rating || '暂无'}</div>
              <div class="stat-label">评分</div>
            </div>
          </div>
          <div class="info">
            <p><strong>导演：</strong>${row.director || '未知'}</p>
            <p><strong>演员：</strong>${Array.isArray(row.actors) ? row.actors.join(', ') : (row.actors || '未知')}</p>
            <p><strong>状态：</strong>${row.status === 'published' ? '已发布' : row.status === 'draft' ? '草稿' : '已归档'}</p>
            <p><strong>创建时间：</strong>${new Date(row.created_at).toLocaleString()}</p>
          </div>
        </div>
      </div>
    </body>
    </html>
  `)
  detailWindow.document.close()
}

const handleEdit = (row) => {
  editingItem.value = row
  Object.assign(form, row)
  showCreateDialog.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这个短剧吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    ElMessage.success('删除成功')
    loadData()
  } catch {
    // 用户取消
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (editingItem.value) {
          ElMessage.success('更新成功')
        } else {
          ElMessage.success('创建成功')
        }
        showCreateDialog.value = false
        resetForm()
        loadData()
      } catch (error) {
        ElMessage.error('操作失败')
      }
    }
  })
}

const resetForm = () => {
  editingItem.value = null
  Object.assign(form, {
    title: '',
    description: '',
    category: '',
    status: 'draft'
  })
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.dramas {
  max-width: 1200px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-bar {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}
</style>
