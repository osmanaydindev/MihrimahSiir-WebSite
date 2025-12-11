<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useField, useForm } from "vee-validate";
import axios from "axios";

const { handleSubmit, handleReset } = useForm({
  validationSchema: {
    username(value) {
      if (!value) return "Bu alan boş bırakılamaz.";
      if (value.length < 3) return "Kullanıcı adı en az 3 karakter olmalıdır.";

      // Boşluk kontrolü
      if (/\s/.test(value)) {
        return "Kullanıcı adı boşluk içeremez.";
      }

      // Türkçe karakter kontrolü (ğ, ü, ş, ı, ö, ç, Ğ, Ü, Ş, İ, Ö, Ç)
      if (/[ğüşıöçĞÜŞİÖÇ]/.test(value)) {
        return "Kullanıcı adı Türkçe karakter içeremez.";
      }

      return true;
    },
    email(value) {
      if (!value) return "Bu alan boş bırakılamaz.";
      // Simple email validation
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!emailRegex.test(value)) return "Geçerli bir e-posta adresi giriniz.";
      return true;
    },
    password(value) {
      if (!value) return "Bu alan boş bırakılamaz.";
      if (value.length < 6) return "Şifre en az 6 karakter olmalıdır.";
      return true;
    },
    passwordConfirm(value) {
      if (!value) return "Bu alan boş bırakılamaz.";
      if (value !== password.value.value) return "Şifreler eşleşmiyor.";
      return true;
    },
  },
});

const loading = ref(false);
const router = useRouter();
const showPassword = ref(false);
const showPasswordConfirm = ref(false);
const popup = ref(false);
const popupMessage = ref("");
const popupColor = ref("red-lighten-1");

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value;
};

const togglePasswordConfirmVisibility = () => {
  showPasswordConfirm.value = !showPasswordConfirm.value;
};

const openPopup = (message, color = "red-lighten-1") => {
  loading.value = false;
  popupMessage.value = message;
  popupColor.value = color;
  popup.value = true;
};

const closePopup = () => {
  loading.value = false;
  popup.value = false;
};

const username = useField("username");
const email = useField("email");
const password = useField("password");
const passwordConfirm = useField("passwordConfirm");

const submitRegister = handleSubmit(async (values) => {
  try {
    loading.value = true;
    const resp = await axios.post("/register", {
      username: values.username,
      email: values.email,
      password: values.password,
    }, { withCredentials: true });

    loading.value = false;

    // Show success message
    openPopup("Kayıt başarılı! Giriş sayfasına yönlendiriliyorsunuz...", "green-lighten-1");

    // Redirect to login page after 2 seconds
    setTimeout(() => {
      router.push("/login");
    }, 2000);
  } catch (error) {
    loading.value = false;
    const errorMessage = error.response?.data?.message || "Kayıt sırasında bir hata oluştu. Lütfen tekrar deneyiniz.";
    openPopup(errorMessage);
  }
});
</script>

<template>
  <div class="register-page">
    <div class="register-container">
      <!-- Register Card -->
      <v-card class="register-card">
        <!-- Header -->
        <div class="register-header">
          <v-icon size="56" color="grey-lighten-1" class="mb-4">mdi-account-plus-outline</v-icon>
          <h1 class="register-title">Hesap Oluştur</h1>
          <p class="register-subtitle">Yeni bir hesap oluşturarak başlayın</p>
        </div>

        <!-- Form -->
        <form class="register-form" @submit.prevent="submitRegister">
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
              v-model="email.value.value"
              :error-messages="email.errorMessage.value"
              label="E-posta"
              type="email"
              variant="outlined"
              density="comfortable"
              prepend-inner-icon="mdi-email"
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

          <div class="form-field">
            <v-text-field
              v-model="passwordConfirm.value.value"
              :error-messages="passwordConfirm.errorMessage.value"
              label="Şifre Tekrar"
              variant="outlined"
              density="comfortable"
              prepend-inner-icon="mdi-lock-check"
              :append-inner-icon="showPasswordConfirm ? 'mdi-eye-off' : 'mdi-eye'"
              :type="showPasswordConfirm ? 'text' : 'password'"
              @click:append-inner="togglePasswordConfirmVisibility"
              clearable
              color="grey-lighten-1"
            />
          </div>

          <v-btn
            type="submit"
            :loading="loading"
            :disabled="loading"
            class="register-btn"
            size="large"
            block
            color="grey-darken-3"
          >
            <v-icon left class="mr-2">mdi-account-plus</v-icon>
            Kayıt Ol
          </v-btn>

          <v-divider class="my-6"></v-divider>

          <v-btn
            @click="router.push('/login')"
            class="login-btn"
            size="large"
            block
            variant="outlined"
            color="grey-lighten-1"
          >
            <v-icon left class="mr-2">mdi-login</v-icon>
            Zaten Hesabım Var
          </v-btn>
        </form>
      </v-card>
    </div>

    <!-- Success/Error Dialog -->
    <v-dialog v-model="popup" max-width="400">
      <v-card :class="popupColor === 'green-lighten-1' ? 'success-dialog' : 'error-dialog'">
        <v-card-title :class="popupColor === 'green-lighten-1' ? 'success-dialog-title' : 'error-dialog-title'">
          <v-icon
            size="48"
            :color="popupColor === 'green-lighten-1' ? 'success' : 'error'"
            class="mb-2"
          >
            {{ popupColor === 'green-lighten-1' ? 'mdi-check-circle' : 'mdi-alert-circle' }}
          </v-icon>
          <span>{{ popupColor === 'green-lighten-1' ? 'Başarılı' : 'Kayıt Başarısız' }}</span>
        </v-card-title>
        <v-card-text :class="popupColor === 'green-lighten-1' ? 'success-dialog-text' : 'error-dialog-text'">
          {{ popupMessage }}
        </v-card-text>
        <v-card-actions class="justify-center pb-4">
          <v-btn
            :color="popupColor === 'green-lighten-1' ? 'success' : 'error'"
            variant="outlined"
            @click="closePopup"
          >
            Tamam
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: inherit;
  padding: clamp(20px, 4vw, 40px);
  animation: fadeIn 0.6s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.register-container {
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

.register-card {
  background: linear-gradient(
    135deg,
    rgba(30, 30, 30, 0.8) 0%,
    rgba(42, 42, 42, 0.9) 100%
  ) !important;
  border: 1px solid #424242;
  border-radius: clamp(12px, 2vw, 16px) !important;
  padding: clamp(32px, 6vw, 48px);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  transition: all 0.3s ease;
}

.register-card:hover {
  box-shadow: 0 24px 70px rgba(0, 0, 0, 0.6);
  transform: translateY(-4px);
}

.register-header {
  text-align: center;
  margin-bottom: clamp(32px, 6vw, 40px);
}

.register-title {
  font-size: clamp(28px, 5vw, 36px);
  font-weight: 700;
  color: #ffffff;
  margin-bottom: clamp(8px, 1.5vw, 12px);
  letter-spacing: 0.5px;
}

.register-subtitle {
  font-size: clamp(14px, 2.5vw, 16px);
  color: #9e9e9e;
  font-weight: 400;
  margin: 0;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: clamp(16px, 3vw, 20px);
}

.form-field {
  width: 100%;
}

.register-btn {
  margin-top: clamp(8px, 2vw, 12px);
  text-transform: none;
  font-size: clamp(15px, 2.5vw, 16px);
  font-weight: 600;
  letter-spacing: 0.3px;
  height: clamp(44px, 8vw, 48px) !important;
  transition: all 0.3s ease;
}

.register-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.3);
}

.login-btn {
  text-transform: none;
  font-size: clamp(15px, 2.5vw, 16px);
  font-weight: 600;
  letter-spacing: 0.3px;
  height: clamp(44px, 8vw, 48px) !important;
  transition: all 0.3s ease;
}

.login-btn:hover {
  background: rgba(158, 158, 158, 0.1);
  transform: translateY(-2px);
}

/* Success Dialog */
.success-dialog {
  background: linear-gradient(135deg, #1e1e1e 0%, #2a2a2a 100%) !important;
  border: 1px solid #66bb6a;
  text-align: center;
}

.success-dialog-title {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: clamp(20px, 4vw, 24px);
  font-size: clamp(20px, 3.5vw, 24px);
  font-weight: 600;
  color: #66bb6a;
}

.success-dialog-text {
  text-align: center;
  font-size: clamp(14px, 2.5vw, 16px);
  color: #bdbdbd;
  padding: 0 clamp(20px, 4vw, 24px) clamp(16px, 3vw, 20px);
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
  .register-page {
    padding: 16px;
  }

  .register-card {
    padding: 24px;
  }
}
</style>
