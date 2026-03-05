<template>
  <div class="login-page">
    <n-card title="登录" class="login-card">
      <n-form ref="formRef" :model="formValue" :rules="rules" label-placement="left">
        <n-form-item label="邮箱" path="email">
          <n-input
            v-model:value="formValue.email"
            placeholder="请输入邮箱"
            @keyup.enter="handleLogin"
          />
        </n-form-item>
        <n-form-item label="密码" path="password">
          <n-input
            v-model:value="formValue.password"
            type="password"
            placeholder="请输入密码"
            show-password-on="click"
            @keyup.enter="handleLogin"
          />
        </n-form-item>
        <n-button
          type="primary"
          block
          :loading="loading"
          :disabled="!isFormValid"
          @click="handleLogin"
        >
          登录
        </n-button>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NCard, NForm, NFormItem, NInput, NButton, useMessage, type FormInst, type FormRules } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const message = useMessage()
const authStore = useAuthStore()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)

const formValue = ref({
  email: '',
  password: '',
})

const rules: FormRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少 6 位', trigger: 'blur' },
  ],
}

const isFormValid = computed(() => {
  return formValue.value.email.length > 0 &&
         formValue.value.password.length >= 6 &&
         /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formValue.value.email)
})

async function handleLogin() {
  if (!isFormValid.value) return

  try {
    loading.value = true

    // 表单验证
    await formRef.value?.validate()

    // TODO: 调用后端 API 进行实际认证
    // const response = await api.post('/auth/login', formValue.value)

    // 模拟登录成功 - 实际项目中应替换为真实 API 调用
    authStore.setToken('mock-jwt-token')
    authStore.setUserInfo('mock-user-id', 'mock-tenant-id')

    message.success('登录成功')

    // 重定向到原本要访问的页面或首页
    const redirect = route.query.redirect as string
    router.push(redirect || '/dashboard')
  } catch (error) {
    if (error instanceof Error) {
      message.error(error.message)
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.login-card {
  width: 400px;
}
</style>
