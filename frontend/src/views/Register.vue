<template>
  <div class="register-container">
    <div class="register-card">
      <h2 class="register-title">注册</h2>
      <form @submit.prevent="handleRegister" class="register-form">
        <div class="form-group">
          <label for="email">邮箱</label>
          <input id="email" v-model="email" type="email" placeholder="请输入邮箱" required class="form-input" />
        </div>
        <div class="form-group">
          <label for="password">密码</label>
          <input id="password" v-model="password" type="password" placeholder="请输入密码（至少6位）" required minlength="6" class="form-input" />
        </div>
        <div class="form-group">
          <label for="confirmPassword">确认密码</label>
          <input id="confirmPassword" v-model="confirmPassword" type="password" placeholder="请再次输入密码" required class="form-input" />
        </div>
        <div v-if="error" class="error-message">{{ error }}</div>
        <div v-if="success" class="success-message">{{ success }}</div>
        <button type="submit" :disabled="loading" class="submit-button">{{ loading ? '注册中...' : '注册' }}</button>
      </form>
      <div class="register-footer">
        <span>已有账号？</span>
        <router-link to="/login" class="link">立即登录</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'Register',
  data() {
    return {
      email: '',
      password: '',
      confirmPassword: '',
      error: '',
      success: '',
      loading: false
    }
  },
  methods: {
    validateForm() {
      if (this.password !== this.confirmPassword) {
        this.error = '两次输入的密码不一致'
        return false
      }
      if (this.password.length < 6) {
        this.error = '密码长度至少为6位'
        return false
      }
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailRegex.test(this.email)) {
        this.error = '请输入有效的邮箱地址'
        return false
      }
      return true
    },
    async handleRegister() {
      this.error = ''
      this.success = ''
      if (!this.validateForm()) return
      this.loading = true
      try {
        const response = await axios.post('http://localhost:8080/register', {
          email: this.email,
          password: this.password
        })
        if (response.data.message) {
          this.success = '注册成功！正在跳转到登录页面...'
          setTimeout(() => { this.$router.push('/login') }, 2000)
        }
      } catch (error) {
        if (error.response && error.response.data.error) {
          this.error = error.response.data.error
        } else {
          this.error = '注册失败，请稍后重试'
        }
      } finally { this.loading = false }
    }
  }
}
</script>

<style scoped>
.register-container { min-height: 100vh; display:flex; justify-content:center; align-items:center; background: linear-gradient(135deg,#667eea 0%,#764ba2 100%); padding:20px }
.register-card { background:white; border-radius:12px; box-shadow:0 10px 40px rgba(0,0,0,0.1); padding:40px; width:100%; max-width:400px }
.success-message { color:#27ae60; font-size:14px; text-align:center; padding:8px; background:#efe; border-radius:6px }
</style>
