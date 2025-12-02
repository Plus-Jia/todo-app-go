<template>
  <div class="dashboard-container">
    <header class="dashboard-header">
      <h1 class="dashboard-title">任务管理</h1>
      <div class="user-info">
        <span class="user-email">{{ userEmail }}</span>
        <button @click="handleLogout" class="logout-button">退出登录</button>
      </div>
    </header>

    <main class="dashboard-main">
      <div class="task-form-card">
        <h2 class="card-title">{{ editingTask ? '编辑任务' : '创建新任务' }}</h2>
        <form @submit.prevent="handleSubmit" class="task-form">
          <div class="form-row">
            <div class="form-group">
              <label for="title">任务标题 *</label>
              <input id="title" v-model="taskForm.title" type="text" placeholder="请输入任务标题" required class="form-input" />
            </div>
            <div class="form-group">
              <label for="deadline">截止日期 *</label>
              <input id="deadline" v-model="taskForm.deadline" type="datetime-local" required class="form-input" />
            </div>
          </div>
          <div class="form-group">
            <label for="desc">任务描述</label>
            <textarea id="desc" v-model="taskForm.desc" placeholder="请输入任务描述（可选）" rows="3" class="form-textarea"></textarea>
          </div>
          <div class="form-actions">
            <button type="submit" class="submit-button" :disabled="loading">
              {{ editingTask ? '更新任务' : '创建任务' }}
            </button>
            <button v-if="editingTask" type="button" @click="cancelEdit" class="cancel-button">
              取消
            </button>
          </div>
        </form>
        <div v-if="error" class="error-message">{{ error }}</div>
      </div>

      <div class="task-list-section">
        <div class="task-filter">
          <button 
            :class="{ active: filter === 'all' }" 
            @click="filter = 'all'"
            class="filter-button"
          >
            全部任务
          </button>
          <button 
            :class="{ active: filter === 'pending' }" 
            @click="filter = 'pending'"
            class="filter-button"
          >
            未完成
          </button>
          <button 
            :class="{ active: filter === 'completed' }" 
            @click="filter = 'completed'"
            class="filter-button"
          >
            已完成
          </button>
        </div>

        <div v-if="loadingTasks" class="loading">加载中...</div>
        <div v-else>
          <div v-if="filteredTasks.length === 0" class="no-tasks">
            {{ filter === 'all' ? '暂无任务' : filter === 'pending' ? '暂无未完成任务' : '暂无已完成任务' }}
          </div>
          <div v-else class="task-list">
            <div v-for="task in filteredTasks" :key="task.id" class="task-item">
              <div class="task-content">
                <h3 :class="{ 'task-title': true, 'completed': task.done }">{{ task.title }}</h3>
                <p v-if="task.desc" class="task-description">{{ task.desc }}</p>
                <p class="task-deadline">截止时间: {{ formatDateTime(task.deadline) }}</p>
              </div>
              <div class="task-actions">
                <button 
                  @click="toggleTask(task)" 
                  :class="{ 'complete-button': !task.done, 'uncomplete-button': task.done }"
                >
                  {{ task.done ? '标记未完成' : '标记完成' }}
                </button>
                <button @click="editTask(task)" class="edit-button">编辑</button>
                <button @click="deleteTask(task.id)" class="delete-button">删除</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Dashboard',
  data() {
    return {
      tasks: [],
      loading: false,
      loadingTasks: false,
      error: '',
      filter: 'all',
      taskForm: { title: '', desc: '', deadline: '' },
      editingTask: null,
      userEmail: ''
    }
  },
  computed: {
    filteredTasks() {
      if (this.filter === 'all') return this.tasks
      else if (this.filter === 'pending') return this.tasks.filter(task => !task.done)
      else return this.tasks.filter(task => task.done)
    }
  },
  mounted() {
    this.checkAuth()
    this.fetchTasks()
    this.loadUserInfo()
  },
  methods: {
    checkAuth() {
      const token = localStorage.getItem('token')
      if (!token) this.$router.push('/login')
    },
    loadUserInfo() {
      const user = localStorage.getItem('user')
      if (user) {
        try {
          const userData = JSON.parse(user)
          this.userEmail = userData.email || ''
        } catch (e) {
          console.error('Error parsing user data:', e)
        }
      }
    },
    getAuthHeaders() {
      const token = localStorage.getItem('token')
      return { 
        Authorization: token ? `Bearer ${token}` : '',
        'Content-Type': 'application/json'
      }
    },
    async fetchTasks() {
      this.loadingTasks = true
      this.error = ''
      try {
        const response = await axios.get('http://localhost:8080/tasks', { 
          headers: this.getAuthHeaders() 
        })
        
        if (response.data && response.data.tasks) {
          this.tasks = response.data.tasks.map(task => ({
            ...task,
            deadline: task.deadline ? new Date(task.deadline).toISOString().slice(0, 16) : ''
          }))
        }
      } catch (error) {
        console.error('Error fetching tasks:', error)
        if (error.response && error.response.status === 401) {
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          this.$router.push('/login')
        } else {
          this.error = '获取任务列表失败'
        }
      } finally { 
        this.loadingTasks = false 
      }
    },
    async handleSubmit() {
      if (!this.taskForm.title || !this.taskForm.deadline) { 
        this.error = '请填写必填字段'; 
        return 
      }
      
      this.loading = true
      this.error = ''
      
      try {
        const taskData = { 
          title: this.taskForm.title, 
          description: this.taskForm.desc || '', 
          deadline: new Date(this.taskForm.deadline).toISOString() 
        }
        
        if (this.editingTask) {
          await axios.put(
            `http://localhost:8080/tasks/${this.editingTask.id}`, 
            { ...taskData, done: this.editingTask.done }, 
            { headers: this.getAuthHeaders() }
          )
        } else {
          await axios.post(
            'http://localhost:8080/tasks', 
            taskData, 
            { headers: this.getAuthHeaders() }
          )
        }
        
        this.resetForm()
        await this.fetchTasks()
      } catch (error) {
        console.error('Error saving task:', error)
        if (error.response && error.response.data && error.response.data.error) {
          this.error = error.response.data.error
        } else {
          this.error = this.editingTask ? '更新任务失败' : '创建任务失败'
        }
      } finally { 
        this.loading = false 
      }
    },
    async toggleTask(task) {
      try {
        await axios.put(
          `http://localhost:8080/tasks/${task.id}`, 
          { 
            title: task.title, 
            description: task.desc, 
            deadline: new Date(task.deadline).toISOString(), 
            done: !task.done 
          }, 
          { headers: this.getAuthHeaders() }
        )
        await this.fetchTasks()
      } catch (error) {
        console.error('Error toggling task:', error)
        this.error = '更新任务状态失败'
      }
    },
    editTask(task) { 
      this.editingTask = task
      this.taskForm = { 
        title: task.title, 
        desc: task.desc || '', 
        deadline: task.deadline 
      }
    },
    cancelEdit() { 
      this.editingTask = null
      this.resetForm() 
    },
    resetForm() { 
      this.taskForm = { title: '', desc: '', deadline: '' }
      this.editingTask = null 
    },
    async deleteTask(taskId) {
      if (!confirm('确定要删除这个任务吗？')) return
      
      try { 
        await axios.delete(
          `http://localhost:8080/tasks/${taskId}`, 
          { headers: this.getAuthHeaders() }
        )
        await this.fetchTasks() 
      } catch (error) {
        console.error('Error deleting task:', error)
        this.error = '删除任务失败'
      }
    },
    handleLogout() { 
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      this.$router.push('/login') 
    },
    formatDateTime(dateTimeString) { 
      if (!dateTimeString) return ''
      
      const date = new Date(dateTimeString)
      return date.toLocaleString('zh-CN', { 
        year: 'numeric', 
        month: '2-digit', 
        day: '2-digit', 
        hour: '2-digit', 
        minute: '2-digit' 
      })
    }
  }
}
</script>

<style scoped>
.dashboard-container {
  min-height: 100vh;
  background: #f5f7fa;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.dashboard-title {
  margin: 0;
  color: #333;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.user-email {
  color: #666;
}

.logout-button {
  padding: 8px 16px;
  background: #dc3545;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.dashboard-main {
  max-width: 1200px;
  margin: 20px auto;
  padding: 0 20px;
}

.task-form-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  margin-bottom: 30px;
}

.card-title {
  margin-top: 0;
  margin-bottom: 20px;
  color: #333;
}

.form-row {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
}

.form-group {
  flex: 1;
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #555;
}

.form-input, .form-textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.form-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.submit-button, .cancel-button {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.submit-button {
  background: #007bff;
  color: white;
}

.submit-button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.cancel-button {
  background: #6c757d;
  color: white;
}

.task-list-section {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.task-filter {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  border-bottom: 1px solid #eee;
  padding-bottom: 15px;
}

.filter-button {
  padding: 8px 16px;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 4px;
}

.filter-button.active {
  background: #007bff;
  color: white;
}

.loading, .no-tasks {
  text-align: center;
  padding: 40px;
  color: #666;
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.task-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 15px;
  border: 1px solid #eee;
  border-radius: 6px;
}

.task-content {
  flex: 1;
}

.task-title {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #333;
}

.task-title.completed {
  text-decoration: line-through;
  color: #999;
}

.task-description {
  margin: 0 0 8px 0;
  color: #666;
  font-size: 14px;
}

.task-deadline {
  margin: 0;
  font-size: 13px;
  color: #999;
}

.task-actions {
  display: flex;
  gap: 8px;
}

.task-actions button {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
}

.complete-button {
  background: #28a745;
  color: white;
}

.uncomplete-button {
  background: #ffc107;
  color: #212529;
}

.edit-button {
  background: #007bff;
  color: white;
}

.delete-button {
  background: #dc3545;
  color: white;
}

.error-message {
  margin-top: 15px;
  padding: 10px;
  background: #f8d7da;
  color: #721c24;
  border-radius: 4px;
  border: 1px solid #f5c6cb;
}
</style>