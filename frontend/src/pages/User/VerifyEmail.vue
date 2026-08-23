<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const ok = ref(false)
const message = ref('E-posta doğrulanıyor...')

onMounted(async () => {
  try {
    const token = String(route.query.token || '')
    const res = await axios.get(`/verify-email?token=${encodeURIComponent(token)}`)
    ok.value = true
    message.value = res.data?.message || 'E-posta adresin doğrulandı.'
  } catch (error: any) {
    ok.value = false
    message.value = error.response?.data?.message || 'Doğrulama bağlantısı geçersiz veya süresi dolmuş.'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="auth-result-page">
    <v-card class="auth-result-card">
      <v-progress-circular v-if="loading" indeterminate color="grey-lighten-1" class="mb-5" />
      <v-icon v-else size="56" :color="ok ? 'success' : 'error'" class="mb-4">
        {{ ok ? 'mdi-check-circle' : 'mdi-alert-circle' }}
      </v-icon>
      <h1 class="title">{{ ok ? 'E-posta Doğrulandı' : 'Doğrulama Başarısız' }}</h1>
      <p class="message">{{ message }}</p>
      <v-btn color="grey-darken-3" class="text-none mt-5" @click="router.push('/login')">
        Giriş Sayfasına Git
      </v-btn>
    </v-card>
  </div>
</template>

<style scoped>
.auth-result-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.auth-result-card {
  width: min(100%, 460px);
  padding: 36px;
  text-align: center;
  background: linear-gradient(135deg, rgba(30, 30, 30, 0.92), rgba(42, 42, 42, 0.92)) !important;
  border: 1px solid #424242;
  border-radius: 12px;
}

.title {
  color: #fff;
  font-size: 28px;
  margin: 0 0 12px;
}

.message {
  color: #bdbdbd;
  margin: 0;
}
</style>
