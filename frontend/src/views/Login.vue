<template>
  <div class="login-container">
    <div class="login-card">
      <h2 class="login-title">登录</h2>
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="email">邮箱</label>
          <input id="email" v-model="email" type="email" placeholder="请输入邮箱" required class="form-input" />
        </div>
        <div class="form-group">
          <label for="password">密码</label>
          <input id="password" v-model="password" type="password" placeholder="请输入密码" required class="form-input" />
        </div>
        <div v-if="error" class="error-message">{{ error }}</div>
        <button type="submit" :disabled="loading" class="submit-button">{{ loading ? '登录中...' : '登录' }}</button>
      </form>
      <div class="login-footer">
        <span>还没有账号？</span>
        <router-link to="/register" class="link">立即注册</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'Login',
  data() {
    return {
      email: '',
      password: '',
      error: '',
      loading: false
    }
  },
  methods: {
    async handleLogin() {
      this.error = ''
      this.loading = true
      try {
        const response = await axios.post('http://localhost:8080/login', {
          email: this.email,
          password: this.password
        })
        if (response.data.token) {
          localStorage.setItem('token', response.data.token)
          localStorage.setItem('userEmail', this.email)
          this.$router.push('/dashboard')
        }
      } catch (error) {
        if (error.response && error.response.data.error) {
          this.error = error.response.data.error
        } else {
          this.error = '登录失败，请稍后重试'
        }
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.login-container { min-height: 100vh; display: flex; justify-content: center; align-items: center; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 20px; }
.login-card { background: white; border-radius: 12px; box-shadow: 0 10px 40px rgba(0,0,0,0.1); padding: 40px; width: 100%; max-width: 400px; }
.login-title { font-size: 28px; font-weight: 600; color: #333; text-align: center; margin-bottom: 30px; }
.login-form { display: flex; flex-direction: column; gap: 20px; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-input { padding: 12px 16px; border: 2px solid #e0e0e0; border-radius: 8px; font-size: 16px; }
.error-message { color: #e74c3c; font-size: 14px; text-align: center; padding: 8px; background: #fee; border-radius: 6px; }
.submit-button { padding: 14px; background: linear-gradient(135deg,#667eea 0%,#764ba2 100%); color: white; border: none; border-radius: 8px; font-size: 16px; font-weight: 600; cursor: pointer; }
.link { color: #667eea; text-decoration: none; font-weight: 500; margin-left: 5px; }
</style>
