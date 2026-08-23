<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const password = ref('')
const passwordConfirm = ref('')
const showPassword = ref(false)
const loading = ref(false)
const message = ref('')
const success = ref(false)

const passwordError = () => {
  if (!password.value) return 'Şifre gerekli.'
  if (password.value.length < 8) return 'Şifre en az 8 karakter olmalıdır.'
  if (password.value.length > 128) return 'Şifre en fazla 128 karakter olabilir.'
  if (!/[a-z]/.test(password.value) || !/[A-Z]/.test(password.value) || !/[0-9]/.test(password.value) || !/[^\w\s]/.test(password.value)) {
    return 'Şifre büyük/küçük harf, sayı ve özel karakter içermelidir.'
  }
  if (password.value !== passwordConfirm.value) return 'Şifreler eşleşmiyor.'
  return ''
}

const submit = async () => {
  const err = passwordError()
  if (err) {
    message.value = err
    success.value = false
    return
  }
  loading.value = true
  try {
    const token = String(route.query.token || '')
    const res = await axios.post('/reset-password', { token, password: password.value })
    success.value = true
    message.value = res.data?.message || 'Şifren güncellendi.'
  } catch (error: any) {
    success.value = false
    message.value = error.response?.data?.message || 'Şifre güncellenemedi.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="reset-page">
    <v-card class="reset-card">
      <v-icon size="56" color="grey-lighten-1" class="mb-4">mdi-lock-reset</v-icon>
      <h1 class="title">Şifre Sıfırla</h1>
      <p class="subtitle">Yeni şifreni belirle.</p>
      <v-text-field
        v-model="password"
        label="Yeni Şifre"
        variant="outlined"
        density="comfortable"
        prepend-inner-icon="mdi-lock"
        :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
        :type="showPassword ? 'text' : 'password'"
        maxlength="128"
        autocomplete="new-password"
        @click:append-inner="showPassword = !showPassword"
      />
      <v-text-field
        v-model="passwordConfirm"
        label="Yeni Şifre Tekrar"
        variant="outlined"
        density="comfortable"
        prepend-inner-icon="mdi-lock-check"
        :type="showPassword ? 'text' : 'password'"
        maxlength="128"
        autocomplete="new-password"
      />
      <v-alert v-if="message" :type="success ? 'success' : 'error'" variant="tonal" class="mb-4">
        {{ message }}
      </v-alert>
      <v-btn :loading="loading" block color="grey-darken-3" class="text-none" @click="submit">
        Şifreyi Güncelle
      </v-btn>
      <v-btn v-if="success" block variant="text" class="text-none mt-3" @click="router.push('/login')">
        Giriş Sayfasına Git
      </v-btn>
    </v-card>
  </div>
</template>

<style scoped>
.reset-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.reset-card {
  width: min(100%, 480px);
  padding: 36px;
  background: linear-gradient(135deg, rgba(30, 30, 30, 0.92), rgba(42, 42, 42, 0.92)) !important;
  border: 1px solid #424242;
  border-radius: 12px;
}

.title {
  color: #fff;
  font-size: 28px;
  margin: 0 0 8px;
  text-align: center;
}

.subtitle {
  color: #bdbdbd;
  text-align: center;
  margin-bottom: 24px;
}
</style>
