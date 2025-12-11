<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const menuLinks = [
  { title: 'Anasayfa', route: '/home', icon: 'mdi-home-outline' },
  { title: 'Sözler', route: '/sozler', icon: 'mdi-format-quote-close' },
  { title: 'Şiirler', route: '/poems', icon: 'mdi-book-open-page-variant-outline' },
  { title: 'Güncel Şiirler', route: '/latest-poems', icon: 'mdi-clock-outline' },
  { title: 'Kitaplar', route: '/books', icon: 'mdi-book-outline' },
]

const handleNavigation = (linkRoute: string, event: MouseEvent) => {
  // Prevent default router-link behavior
  event.preventDefault()

  // Only navigate if we're not already on that route
  if (route.path !== linkRoute) {
    router.push(linkRoute).catch(() => {
      // Silently catch navigation duplicated errors
    })
  }
}
</script>

<template>
  <footer class="app-footer">
    <div class="footer-content">
      <!-- Brand Section -->
      <div class="footer-brand">
        <v-icon size="28" class="footer-icon">mdi-feather</v-icon>
        <span class="footer-title">Mihrimâh</span>
      </div>

      <!-- Menu Links Section -->
      <div class="footer-menu opacity-100">
        <router-link
          v-for="link in menuLinks"
          :key="link.route"
          :to="link.route"
          class="menu-link"
          @click="(e) => handleNavigation(link.route, e)"
        >
          <v-icon size="16" class="mr-1">{{ link.icon }}</v-icon>
          {{ link.title }}
        </router-link>
      </div>

      <!-- Divider -->
      <div class="footer-divider"></div>

      <!-- Social Links -->
      <div class="footer-social">
        <a
          href="https://github.com/osmanaydindev"
          target="_blank"
          rel="noopener noreferrer"
          class="social-link"
          aria-label="GitHub"
        >
          <v-icon size="20">mdi-github</v-icon>
          <span class="social-text">GitHub</span>
        </a>
        <a
          href="https://www.linkedin.com/in/osman-ayd%C4%B1n-a22860205/"
          target="_blank"
          rel="noopener noreferrer"
          class="social-link"
          aria-label="LinkedIn"
        >
          <v-icon size="20">mdi-linkedin</v-icon>
          <span class="social-text">LinkedIn</span>
        </a>
      </div>

      <!-- Designer Credit -->
      <div class="footer-designer">
        <v-icon size="16" class="mr-1">mdi-palette-outline</v-icon>
        <span>Designed and created by Osman Aydın</span>
      </div>

      <!-- Copyright -->
      <div class="footer-copyright">
        <span>&copy; {{ new Date().getFullYear() }} Mihrimâh. Tüm hakları saklıdır.</span>
      </div>
    </div>
  </footer>
</template>

<style scoped>
/* Footer Container */
.app-footer {
  width: 100%;
  background: linear-gradient(135deg, #1a1a1a 30%, #2d2d2d 100%);
  opacity: 90%;
  border-top: 1px solid #424242;
  padding: 48px 24px 32px;
  margin-top: auto;
  position: relative;
  overflow: hidden;
}

.app-footer::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg,
    transparent 0%,
    #27b0ab 20%,
    #3a9ab7 50%,
    #3f51b5 80%,
    transparent 100%
  );
  opacity: 0.6;
  animation: shimmer 3s ease-in-out infinite;
}

@keyframes shimmer {
  0%, 100% {
    opacity: 0.4;
  }
  50% {
    opacity: 0.8;
  }
}

/* Footer Content */
.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
}

/* Brand Section */
.footer-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  animation: fadeInUp 0.6s ease-out;
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

/* Menu Links Section */
.footer-menu {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 16px 24px;
  margin-top: 24px;
  animation: fadeInUp 0.6s ease-out 0.1s both;
}

.menu-link {
  display: flex;
  align-items: center;
  color: #bdbdbd;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.3px;
  transition: all 0.3s ease;
  padding: 6px 12px;
  border-radius: 6px;
  position: relative;
}

.menu-link::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  width: 0;
  height: 2px;
  background: linear-gradient(90deg, #2797b0, #3a9ab7);
  transform: translateX(-50%);
  transition: width 0.3s ease;
}

.menu-link:hover {
  color: #ffffff;
  background: rgba(255, 255, 255, 0.05);
  transform: translateY(-2px);
}

.menu-link:hover::after {
  width: 80%;
}

.menu-link.router-link-active {
  color: #2797b0;
  background: rgba(39, 151, 176, 0.1);
}

.footer-icon {
  color: #bdbdbd;
  transition: all 0.3s ease;
}

.footer-brand:hover .footer-icon {
  color: #2797b0;
  transform: rotate(15deg) scale(1.1);
}

.footer-title {
  font-size: 24px;
  font-weight: 600;
  color: #ffffff;
  letter-spacing: 0.5px;
  transition: all 0.3s ease;
}

.footer-brand:hover .footer-title {
  color: #e0e0e0;
}

/* Divider */
.footer-divider {
  width: 80px;
  height: 2px;
  background: linear-gradient(90deg,
    transparent 0%,
    rgba(189, 189, 189, 0.5) 50%,
    transparent 100%
  );
  animation: fadeInUp 0.6s ease-out 0.2s both;
}

/* Social Links */
.footer-social {
  display: flex;
  gap: 24px;
  animation: fadeInUp 0.6s ease-out 0.4s both;
}

.social-link {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: #bdbdbd;
  text-decoration: none;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.social-link::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg,
    transparent 0%,
    rgba(255, 255, 255, 0.1) 50%,
    transparent 100%
  );
  transition: left 0.5s ease;
}

.social-link:hover::before {
  left: 100%;
}

.social-link:hover {
  transform: translateY(-4px);
  border-color: rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
}

.social-link:hover .social-text {
  color: #ffffff;
}

.social-link:active {
  transform: translateY(-2px);
}

.social-text {
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.3px;
  transition: color 0.3s ease;
}

/* Designer Credit */
.footer-designer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  font-size: 14px;
  color: #9e9e9e;
  font-weight: 500;
  letter-spacing: 0.3px;
  animation: fadeInUp 0.6s ease-out 0.5s both;
  margin-top: 12px;
  padding: 8px 16px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.05);
  transition: all 0.3s ease;
}

.footer-designer:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
  color: #bdbdbd;
  transform: translateY(-2px);
}

/* Copyright */
.footer-copyright {
  font-size: 13px;
  color: #757575;
  letter-spacing: 0.2px;
  text-align: center;
  animation: fadeInUp 0.6s ease-out 0.7s both;
  margin-top: 8px;
}

/* Responsive Design */
@media (max-width: 768px) {
  .app-footer {
    padding: 36px 20px 24px;
  }

  .footer-brand {
    gap: 10px;
  }

  .footer-icon {
    font-size: 24px;
  }

  .footer-title {
    font-size: 20px;
  }

  .footer-menu {
    gap: 12px 16px;
    margin-top: 20px;
  }

  .menu-link {
    font-size: 13px;
    padding: 6px 10px;
  }

  .footer-social {
    flex-direction: column;
    gap: 12px;
    width: 100%;
    max-width: 280px;
  }

  .social-link {
    justify-content: center;
    width: 100%;
  }

  .footer-designer {
    font-size: 13px;
    padding: 8px 12px;
  }

  .footer-copyright {
    font-size: 12px;
    padding: 0 12px;
  }
}

@media (max-width: 480px) {
  .footer-title {
    font-size: 18px;
  }

  .footer-divider {
    width: 60px;
  }

  .footer-menu {
    gap: 10px 12px;
  }

  .menu-link {
    font-size: 12px;
    padding: 5px 8px;
  }

  .footer-designer {
    font-size: 12px;
    padding: 6px 10px;
  }
}
</style>
