<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="brand">Todo App</div>
      <h2 class="title">创建账户</h2>

      <form @submit.prevent="handleRegister" class="form">
        <label class="label">邮箱</label>
        <input v-model="email" type="email" placeholder="例如: you@example.com" required class="input" />

        <label class="label">密码</label>
        <input v-model="password" type="password" placeholder="至少 6 位" required minlength="6" class="input" />

        <label class="label">确认密码</label>
        <input v-model="confirmPassword" type="password" placeholder="再次输入密码" required class="input" />

        <div v-if="error" class="toast error">{{ error }}</div>
        <div v-if="success" class="toast success">{{ success }}</div>

        <button type="submit" :disabled="loading" class="btn primary">{{ loading ? '注册中...' : '注册' }}</button>
      </form>

      <div class="alt">
        已有账号？ <router-link to="/login">去登录</router-link>
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
        const res = await axios.post('http://localhost:8080/register', {
          email: this.email,
          password: this.password
        })
        if (res.data.message) {
          this.success = '注册成功，2秒后跳转到登录页'
          setTimeout(() => { this.$router.push('/login') }, 2000)
        }
      } catch (err) {
        this.error = (err.response && err.response.data && err.response.data.error) ? err.response.data.error : '注册失败，请稍后重试'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.auth-page { min-height: 100vh; display:flex; align-items:center; justify-content:center; background: linear-gradient(135deg,#eef2ff 0%, #f7f7ff 100%); padding:24px }
.auth-card { width:100%; max-width:420px; background:#fff; border-radius:12px; box-shadow: 0 8px 30px rgba(24,39,75,0.08); padding:28px; }
.brand { font-weight:700; color:#4f46e5; font-size:18px; margin-bottom:6px }
.title { margin:0 0 16px 0; font-size:22px; color:#0f172a }
.form { display:flex; flex-direction:column; gap:12px }
.label { font-size:13px; color:#475569 }
.input { padding:12px 14px; border-radius:8px; border:1px solid #e6e9ef; font-size:14px }
.input:focus { outline:none; border-color:#6366f1; box-shadow:0 0 0 4px rgba(99,102,241,0.08) }
.btn { padding:12px 16px; border-radius:10px; font-weight:600; border:none; cursor:pointer }
.btn.primary { background: linear-gradient(90deg,#6366f1,#7c3aed); color:#fff }
.alt { margin-top:14px; font-size:13px; color:#64748b; text-align:center }
.alt a { color:#6366f1; font-weight:600 }
.toast { padding:10px; border-radius:8px; text-align:center; font-size:13px }
.toast.error { background:#fff5f5; color:#c53030; border:1px solid #f8d7da }
.toast.success { background:#f0fdf4; color:#166534; border:1px solid #bbf7d0 }
</style>
