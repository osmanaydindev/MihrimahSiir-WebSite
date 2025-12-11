<script setup lang="ts">
import { ref } from "vue";
import {useRouter} from "vue-router";
import { useField, useForm } from "vee-validate";
import {useAppStore} from "../../store/app";
import axios from "axios";

const store = useAppStore();

const { handleSubmit, handleReset } = useForm({
  validationSchema: {
    username (value:any) {
      if(!value) return "Bu alan boş bırakılmamaz."
      return true
    },
    password (value:any) {
      if (value && value.length>5) return true
      return 'Şifre 5 karakterden uzun olmalıdır.'
    },
  },
})
const loading = ref(false);
const router = useRouter()
const showPassword = ref(false);
const popup = ref(false);
const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value;
};
const openPopup = () =>{
  loading.value = false;
  popup.value = true;
}
const closePopup = () => {
  loading.value = false;
  popup.value = false;

}
const username = useField('username')
const password = useField('password')
const submitLogin = handleSubmit(async (values) => {
  try {
    loading.value = true;
    const resp = await axios.post('/login', values, { withCredentials: true });
    localStorage.setItem("token", resp.data.cookie.value);
    loading.value = false;
    await router.push("/home");
  } catch (error) {
    openPopup();
  }

});
</script>

<template>
  <div class="login-page">
    <div class="login-container">
      <!-- Login Card -->
      <v-card class="login-card">
        <!-- Header -->
        <div class="login-header">
          <v-icon size="56" color="grey-lighten-1" class="mb-4">mdi-account-circle</v-icon>
          <h1 class="login-title">Hoş Geldiniz</h1>
          <p class="login-subtitle">Hesabınıza giriş yapın</p>
        </div>

        <!-- Form -->
        <form class="login-form" @submit.prevent="submitLogin">
          <div class="form-field">
            <v-text-field
              v-model="username.value.value"
              :error-messages="username.errorMessage.value"
              label="Kullanıcı Adı"
              variant="outlined"
              density="comfortable"
              prepend-inner-icon="mdi-account"
              clearable
              color="grey-lighten-1"
            />
          </div>

          <div class="form-field">
            <v-text-field
              v-model="password.value.value"
              :error-messages="password.errorMessage.value"
              label="Şifre"
              variant="outlined"
              density="comfortable"
              prepend-inner-icon="mdi-lock"
              :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
              :type="showPassword ? 'text' : 'password'"
              @click:append-inner="togglePasswordVisibility"
              clearable
              color="grey-lighten-1"
            />
          </div>

          <v-btn
            type="submit"
            :loading="loading"
            :disabled="loading"
            class="login-btn"
            size="large"
            block
            color="grey-darken-3"
          >
            <v-icon left class="mr-2">mdi-login</v-icon>
            Giriş Yap
          </v-btn>

          <v-divider class="my-6"></v-divider>

          <v-btn
            @click="router.push('/register')"
            class="register-btn"
            size="large"
            block
            variant="outlined"
            color="grey-lighten-1"
          >
            <v-icon left class="mr-2">mdi-account-plus</v-icon>
            Hesap Oluştur
          </v-btn>
        </form>
      </v-card>
    </div>

    <!-- Error Dialog -->
    <v-dialog v-model="popup" max-width="400">
      <v-card class="error-dialog">
        <v-card-title class="error-dialog-title">
          <v-icon size="48" color="error" class="mb-2">mdi-alert-circle</v-icon>
          <span>Giriş Başarısız</span>
        </v-card-title>
        <v-card-text class="error-dialog-text">
          Lütfen kullanıcı adınızı ve şifrenizi kontrol ediniz.
        </v-card-text>
        <v-card-actions class="justify-center pb-4">
          <v-btn color="error" variant="outlined" @click="closePopup">
            Tamam
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  //background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 50%, #1a1a1a 100%);
  background: inherit;
  padding: clamp(20px, 4vw, 40px);
  animation: fadeIn 0.6s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.login-container {
  width: 100%;
  max-width: 480px;
  animation: slideUp 0.8s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-card {
  //background: linear-gradient(135deg, #1e1e1e 0%, #2a2a2a 100%) !important;
  background: linear-gradient(
    135deg,
    rgba(30, 30, 30, 0.9) 0%,
    rgba(42, 42, 42, 0.9) 100%
  ) !important;
  border: 1px solid #424242;
  border-radius: clamp(12px, 2vw, 16px) !important;
  padding: clamp(32px, 6vw, 48px);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  transition: all 0.3s ease;
}

.login-card:hover {
  box-shadow: 0 24px 70px rgba(0, 0, 0, 0.6);
  transform: translateY(-4px);
}

.login-header {
  text-align: center;
  margin-bottom: clamp(32px, 6vw, 40px);
}

.login-title {
  font-size: clamp(28px, 5vw, 36px);
  font-weight: 700;
  color: #ffffff;
  margin-bottom: clamp(8px, 1.5vw, 12px);
  letter-spacing: 0.5px;
}

.login-subtitle {
  font-size: clamp(14px, 2.5vw, 16px);
  color: #9e9e9e;
  font-weight: 400;
  margin: 0;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: clamp(16px, 3vw, 20px);
}

.form-field {
  width: 100%;
}

.login-btn {
  margin-top: clamp(8px, 2vw, 12px);
  text-transform: none;
  font-size: clamp(15px, 2.5vw, 16px);
  font-weight: 600;
  letter-spacing: 0.3px;
  height: clamp(44px, 8vw, 48px) !important;
  transition: all 0.3s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.3);
}

.register-btn {
  text-transform: none;
  font-size: clamp(15px, 2.5vw, 16px);
  font-weight: 600;
  letter-spacing: 0.3px;
  height: clamp(44px, 8vw, 48px) !important;
  transition: all 0.3s ease;
}

.register-btn:hover {
  background: rgba(158, 158, 158, 0.1);
  transform: translateY(-2px);
}

/* Error Dialog */
.error-dialog {
  background: linear-gradient(135deg, #1e1e1e 0%, #2a2a2a 100%) !important;
  border: 1px solid #ef5350;
  text-align: center;
}

.error-dialog-title {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: clamp(20px, 4vw, 24px);
  font-size: clamp(20px, 3.5vw, 24px);
  font-weight: 600;
  color: #ef5350;
}

.error-dialog-text {
  text-align: center;
  font-size: clamp(14px, 2.5vw, 16px);
  color: #bdbdbd;
  padding: 0 clamp(20px, 4vw, 24px) clamp(16px, 3vw, 20px);
}

/* Responsive adjustments */
@media (max-width: 480px) {
  .login-page {
    padding: 16px;
  }

  .login-card {
    padding: 24px;
  }
}
</style>
