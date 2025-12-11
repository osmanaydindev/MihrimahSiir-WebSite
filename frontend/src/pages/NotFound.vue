<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/home')
}

const goBack = () => {
  router.back()
}
</script>

<template>
  <div class="error-page-container">
    <div class="error-content">
      <!-- Animated 404 -->
      <div class="error-number-wrapper">
        <h1 class="error-number">404</h1>
        <div class="error-shadow">404</div>
      </div>

      <!-- Error Message -->
      <div class="error-message">
        <h2 class="error-title">Sayfa Bulunamadı</h2>
        <p class="error-description">
          Aradığınız sayfa mevcut değil veya taşınmış olabilir.
        </p>
      </div>

      <!-- Icon -->
      <div class="error-icon-wrapper">
        <v-icon size="120" color="grey-lighten-1" class="error-icon">
          mdi-file-question-outline
        </v-icon>
      </div>

      <!-- Action Buttons -->
      <div class="error-actions">
        <v-btn
          color="blue-lighten-1"
          variant="flat"
          size="large"
          rounded="lg"
          @click="goHome"
          class="action-btn text-none"
        >
          <v-icon size="20" class="mr-2">mdi-home</v-icon>
          Ana Sayfaya Dön
        </v-btn>

        <v-btn
          color="grey-lighten-1"
          variant="outlined"
          size="large"
          rounded="lg"
          @click="goBack"
          class="action-btn text-none"
        >
          <v-icon size="20" class="mr-2">mdi-arrow-left</v-icon>
          Geri Dön
        </v-btn>
      </div>

      <!-- Quick Links -->
      <div class="quick-links">
        <p class="quick-links-title">Hızlı Bağlantılar</p>
        <div class="quick-links-grid">
          <v-btn
            variant="text"
            size="small"
            @click="router.push('/poems')"
            class="quick-link-btn text-none"
          >
            <v-icon size="16" class="mr-1">mdi-book-open-page-variant</v-icon>
            Şiirler
          </v-btn>
          <v-btn
            variant="text"
            size="small"
            @click="router.push('/books')"
            class="quick-link-btn text-none"
          >
            <v-icon size="16" class="mr-1">mdi-book-multiple</v-icon>
            Kitaplar
          </v-btn>
          <v-btn
            variant="text"
            size="small"
            @click="router.push('/sozler')"
            class="quick-link-btn text-none"
          >
            <v-icon size="16" class="mr-1">mdi-format-quote-close</v-icon>
            Sözler
          </v-btn>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Container */
.error-page-container {
  min-height: calc(100vh - 200px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  animation: fadeIn 0.6s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

/* Content */
.error-content {
  max-width: 600px;
  width: 100%;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 32px;
}

/* 404 Number */
.error-number-wrapper {
  position: relative;
  animation: bounceIn 0.8s ease-out;
}

@keyframes bounceIn {
  0% {
    opacity: 0;
    transform: scale(0.3);
  }
  50% {
    transform: scale(1.05);
  }
  70% {
    transform: scale(0.9);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

.error-number {
  font-size: 150px;
  font-weight: 900;
  background: linear-gradient(135deg, #2797b0 0%, #3f51b5 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
  line-height: 1;
  position: relative;
  z-index: 2;
}

.error-shadow {
  position: absolute;
  top: 10px;
  left: 10px;
  font-size: 150px;
  font-weight: 900;
  color: rgba(39, 151, 176, 0.1);
  z-index: 1;
  line-height: 1;
}

/* Error Message */
.error-message {
  animation: fadeInUp 0.8s ease-out 0.2s both;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.error-title {
  font-size: 36px;
  font-weight: 700;
  color: #ffffff;
  margin: 0 0 16px 0;
  letter-spacing: 0.5px;
}

.error-description {
  font-size: 18px;
  color: #bdbdbd;
  margin: 0;
  line-height: 1.6;
  letter-spacing: 0.3px;
}

/* Icon */
.error-icon-wrapper {
  animation: float 3s ease-in-out infinite;
  opacity: 0.6;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.error-icon {
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.3));
}

/* Action Buttons */
.error-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  justify-content: center;
  animation: fadeInUp 0.8s ease-out 0.4s both;
}

.action-btn {
  font-weight: 500;
  letter-spacing: 0.5px;
  transition: all 0.3s ease;
  min-width: 180px;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.4);
}

/* Quick Links */
.quick-links {
  animation: fadeInUp 0.8s ease-out 0.6s both;
  width: 100%;
}

.quick-links-title {
  font-size: 14px;
  color: #9e9e9e;
  margin: 0 0 12px 0;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 600;
}

.quick-links-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
}

.quick-link-btn {
  font-weight: 500;
  color: #bdbdbd;
  transition: all 0.3s ease;
}

.quick-link-btn:hover {
  color: #ffffff;
  background: rgba(255, 255, 255, 0.05);
}

/* Responsive Design */
@media (max-width: 768px) {
  .error-page-container {
    padding: 24px 12px;
  }

  .error-number,
  .error-shadow {
    font-size: 100px;
  }

  .error-title {
    font-size: 28px;
  }

  .error-description {
    font-size: 16px;
  }

  .error-icon-wrapper {
    display: none;
  }

  .action-btn {
    min-width: 150px;
  }
}

@media (max-width: 480px) {
  .error-number,
  .error-shadow {
    font-size: 80px;
  }

  .error-title {
    font-size: 24px;
  }

  .error-description {
    font-size: 14px;
  }

  .error-actions {
    flex-direction: column;
    width: 100%;
  }

  .action-btn {
    width: 100%;
    min-width: unset;
  }

  .quick-links-grid {
    flex-direction: column;
  }

  .quick-link-btn {
    width: 100%;
  }
}
</style>
