<template>
  <div class="episodes">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>剧集管理</span>
          <el-button type="primary" @click="showCreateDialog = true">
            <el-icon><Plus /></el-icon>
            新增剧集
          </el-button>
        </div>
      </template>

      <!-- 搜索筛选 -->
      <div class="search-bar">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-input
              v-model="searchForm.title"
              placeholder="搜索剧集标题"
              clearable
              @input="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </el-col>
          <el-col :span="6">
            <el-select
              v-model="searchForm.drama_id"
              placeholder="选择短剧"
              clearable
              @change="handleSearch"
            >
              <el-option
                v-for="drama in dramas"
                :key="drama.id"
                :label="drama.title"
                :value="drama.id"
              />
            </el-select>
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
        <el-table-column label="缩略图" width="80">
          <template #default="{ row }">
            <el-image
              v-if="row.thumbnail"
              :src="row.thumbnail"
              :preview-src-list="[row.thumbnail]"
              fit="cover"
              style="width: 60px; height: 40px; border-radius: 4px;"
              :preview-teleported="true"
            />
            <span v-else class="text-gray-400">无缩略图</span>
          </template>
        </el-table-column>
        <el-table-column prop="drama_title" label="所属短剧" min-width="120" />
        <el-table-column prop="title" label="剧集标题" min-width="150" />
        <el-table-column prop="episode_num" label="集数" width="80" />
        <el-table-column prop="duration" label="时长" width="100">
          <template #default="{ row }">
            {{ formatDuration(row.duration) }}
          </template>
        </el-table-column>
        <el-table-column label="播放" width="80">
          <template #default="{ row }">
            <el-button
              v-if="row.video_url"
              size="small"
              type="primary"
              @click="playVideo(row.video_url)"
            >
              播放
            </el-button>
            <span v-else class="text-gray-400">无视频</span>
          </template>
        </el-table-column>
        <el-table-column prop="view_count" label="观看数" width="100" />
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
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
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
      :title="editingItem ? '编辑剧集' : '新增剧集'"
      width="600px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="80px"
      >
        <el-form-item label="所属短剧" prop="drama_id">
          <el-select v-model="form.drama_id" style="width: 100%">
            <el-option
              v-for="drama in dramas"
              :key="drama.id"
              :label="drama.title"
              :value="drama.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="剧集标题" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="集数" prop="episode_num">
          <el-input-number v-model="form.episode_num" :min="1" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="时长(秒)" prop="duration">
          <el-input-number v-model="form.duration" :min="0" />
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
  drama_id: '',
  status: ''
})

const form = reactive({
  drama_id: '',
  title: '',
  episode_num: 1,
  description: '',
  duration: 0,
  status: 'draft'
})

const formRules = {
  drama_id: [{ required: true, message: '请选择所属短剧', trigger: 'change' }],
  title: [{ required: true, message: '请输入剧集标题', trigger: 'blur' }],
  episode_num: [{ required: true, message: '请输入集数', trigger: 'blur' }]
}

const tableData = ref([])
const dramas = ref([])
const pagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

const formatDuration = (seconds) => {
  if (!seconds) return '0分0秒'
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes}分${remainingSeconds}秒`
}

const playVideo = (videoUrl) => {
  if (videoUrl) {
    // 创建一个视频播放的弹窗
    const videoWindow = window.open('', '_blank', 'width=800,height=600,scrollbars=yes,resizable=yes')
    videoWindow.document.write(`
      <!DOCTYPE html>
      <html>
      <head>
        <title>视频播放</title>
        <style>
          body { margin: 0; padding: 20px; background: #000; display: flex; justify-content: center; align-items: center; min-height: 100vh; }
          video { max-width: 100%; max-height: 100%; }
        </style>
      </head>
      <body>
        <video controls autoplay style="width: 100%; height: auto;">
          <source src="${videoUrl}" type="video/mp4">
          您的浏览器不支持视频播放。
        </video>
      </body>
      </html>
    `)
    videoWindow.document.close()
  } else {
    ElMessage.warning('视频地址不存在')
  }
}

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

const loadDramas = async () => {
  try {
    // 模拟数据
    dramas.value = [
      { id: 1, title: '霸道总裁爱上我' },
      { id: 2, title: '重生之商业帝国' }
    ]
  } catch (error) {
    ElMessage.error('加载短剧列表失败')
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const response = await api.get('/api/admin/episodes', {
      params: {
        page: pagination.page,
        page_size: pagination.size,
        ...searchForm
      }
    })

    if (response.data.success) {
      // 处理剧集数据，添加短剧标题
      const episodes = response.data.data.episodes.map(episode => ({
        ...episode,
        drama_title: episode.drama ? episode.drama.title : '未知短剧'
      }))
      tableData.value = episodes
      pagination.total = response.data.data.total
    } else {
      ElMessage.error(response.data.message || '加载数据失败')
    }
  } catch (error) {
    console.error('加载剧集数据失败:', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleEdit = (row) => {
  editingItem.value = row
  Object.assign(form, row)
  showCreateDialog.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这个剧集吗？', '提示', {
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
    drama_id: '',
    title: '',
    episode_num: 1,
    description: '',
    duration: 0,
    status: 'draft'
  })
}

onMounted(() => {
  loadDramas()
  loadData()
})
</script>

<style scoped>
.episodes {
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
