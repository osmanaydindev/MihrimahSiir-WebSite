<script lang="ts" setup>
import axios from "axios";
import { computed, onMounted, ref, watch } from 'vue'
import {useRouter} from "vue-router"
import {useAppStore} from "../store/app";
import { useDisplay } from "vuetify";
const router = useRouter()
const store = useAppStore()
const { mobile } = useDisplay()

// Bekleyen kitap isteği rozeti
const pendingBookRequests = ref(0)
onMounted(async () => {
  try {
    const res = await axios.get('/get-book-request-count')
    pendingBookRequests.value = res.data?.pending || 0
  } catch (e) {
    pendingBookRequests.value = 0
  }
})
const logout = async()=>{
  await axios.post("/logout")
  await router.push("/login")
  localStorage.removeItem('user')
  localStorage.removeItem('token')
  store.$reset()
}
const drawer = ref(true)
const rail = ref(false)
const navItems = [
  { to: '/home', icon: 'mdi-home-switch-outline', title: 'Sayfaya Git', value: 'home' },
  { to: '/log-management', icon: 'mdi-information', title: 'Giriş Yönet', value: 'logs' },
  { to: '/admin-management', icon: 'mdi-account-group-outline', title: 'Admin Yönet', value: 'admins' },
  { to: '/poem-management', icon: 'mdi-text-box-plus-outline', title: 'Şiir Yönet', value: 'poems' },
  { to: '/book-management', icon: 'mdi-book', title: 'Kitap Yönet', value: 'books' },
  { to: '/book-request-management', icon: 'mdi-book-plus-outline', title: 'Kitap İstekleri', value: 'book-requests' },
  { to: '/author-management', icon: 'mdi-account-edit', title: 'Yazar Yönet', value: 'authors' },
  { to: '/reminder-management', icon: 'mdi-bell', title: 'Hatırlatıcı Yönet', value: 'reminders' },
  { to: '/homepage-management', icon: 'mdi-home-edit', title: 'Anasayfa Yönet', value: 'homepage' },
  { to: '/mihrimah-card-management', icon: 'mdi-card-text', title: 'Mihrimah Kart Yönet', value: 'mihrimah-card' },
]
const drawerMode = computed(() => mobile.value ? 'temporary' : 'permanent')

watch(mobile, (isMobile) => {
  drawer.value = !isMobile
  rail.value = false
}, { immediate: true })
</script>

<template>
  <v-app>
      <v-card class="admin-shell">
        <v-layout class="admin-layout">
          <v-app-bar
            v-if="mobile"
            class="admin-app-bar"
            elevation="0"
            density="comfortable"
          >
            <template v-slot:prepend>
              <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
            </template>
            <v-app-bar-title>Admin Paneli</v-app-bar-title>
          </v-app-bar>

          <v-navigation-drawer
            v-model="drawer"
            :rail="!mobile && rail"
            :temporary="drawerMode === 'temporary'"
            :permanent="drawerMode === 'permanent'"
            class="admin-drawer"
            @click="!mobile && (rail = false)"
          >
            <v-list-item
              prepend-icon="mdi-view-dashboard-outline"
              title="Admin Paneli"
              @click="router.push('/panel')"
              nav
            >
              <template v-if="!mobile" v-slot:append>
                <v-btn
                  style="text-transform: none;"
                  icon="mdi-chevron-left"
                  variant="text"
                  @click.stop="rail = !rail"
                ></v-btn>
              </template>
            </v-list-item>

            <v-divider></v-divider>

            <v-list density="compact" nav>
              <router-link
                v-for="item in navItems"
                :key="item.value"
                :to="item.to"
                class="text-white text-decoration-none"
                @click="mobile && (drawer = false)"
              >
                <v-list-item :prepend-icon="item.icon" :title="item.title" :value="item.value">
                  <template v-if="pendingBookRequests > 0" v-slot:append>
                    <v-badge
                      v-if="item.value === 'book-requests'"
                      :content="pendingBookRequests"
                      color="error"
                      inline
                    />
                  </template>
                </v-list-item>
              </router-link>
              <v-list-item
                prepend-icon="mdi-logout"
                @click="logout"
                title="Çıkış Yap"
              />
            </v-list>
          </v-navigation-drawer>
          <v-main class="admin-main">
            <div class="admin-content">
              <router-view />
            </div>
          </v-main>
        </v-layout>
      </v-card>


  </v-app>
</template>

<style scoped>
.admin-shell,
.admin-layout {
  min-height: 100vh;
  background: #121212;
}

.admin-app-bar,
.admin-drawer {
  background: linear-gradient(180deg, #1a1a1a 0%, #2d2d2d 100%) !important;
  border-color: #424242 !important;
}

.admin-main {
  min-height: 100vh;
  background: #121212;
}

.admin-content {
  width: min(100%, 1280px);
  margin: 0 auto;
  padding: 24px;
}

@media (max-width: 600px) {
  .admin-content {
    padding: 12px;
  }
}

@media (max-width: 960px) {
  .admin-content :deep(.v-card) {
    width: 100% !important;
    max-width: 100% !important;
  }

  .admin-content :deep(.v-card-title) {
    white-space: normal;
    line-height: 1.25;
  }

  .admin-content :deep(.v-card-text) {
    padding: 12px;
  }
}
</style>
