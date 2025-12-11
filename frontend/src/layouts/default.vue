<script setup lang="ts">
import {computed, onMounted, onUnmounted, ref, watch} from "vue"
import {useAppStore}  from "../store/app";
import axios          from "axios";
import router         from "../router";
import Footer         from "../components/shared/Footer.vue";
import {wsService}    from "../services/websocket";
import { getImageUrl } from "../utils/imageHelper";

const store       = useAppStore()
const drawer             = ref(false)
const logoutPopup        = ref(false)
const group              = ref(null)
const likedPoems         = ref([])
const bookmarksId        = ref([])
const readsBooks         = ref([])
const role_id            = ref(0)
const pendingRequestsCount = ref(0)

const items   = computed(() => {
  const allItems = [
    {
      title: 'Admin Paneli',
      value: 'panel',
      route: '/panel',
      icon: 'mdi-view-dashboard-outline',
      requiresRole: [1],
    },
    {
      title: 'Anasayfa',
      value: 'home',
      route: '/home',
      icon: 'mdi-home-outline',
    },
    {
      title: 'Sözler',
      value: 'sozler',
      route: '/sozler',
      icon: 'mdi-format-quote-close',
      requiresRole: [1, 2, 3],
    },
    {
      title: 'Şiirler',
      value: 'poems-group',
      icon: 'mdi-book-open-page-variant-outline',
      subLinks: [
        {
          title: 'Tüm Şiirler',
          value: 'all-poems',
          route: '/poems',
          icon: 'mdi-text-box-multiple-outline',
        },
        {
          title: 'Güncel Şiirler',
          value: 'latest-poems',
          route: '/latest-poems',
          icon: 'mdi-clock-outline',
        },
        {
          title: 'Popüler Şiirler',
          value: 'popular-poems',
          route: '/popular-poems',
          icon: 'mdi-fire',
        },
        {
          title: 'Beğenilen Şiirler',
          value: 'liked-poems',
          route: '/liked-poems',
          icon: 'mdi-heart-outline',
        },
        {
          title: 'Kaydedilen Şiirler',
          value: 'bookmark',
          route: '/bookmark',
          icon: 'mdi-bookmark-outline',
        }
      ],
    },
    {
      title: 'Kitaplar',
      value: 'books-group',
      icon: 'mdi-book-outline',
      subLinks: [
        {
          title: 'Tüm Kitaplar',
          value: 'all-books',
          route: '/books',
          icon: 'mdi-bookshelf',
        },
        {
          title: 'Okunanlar',
          value: 'reads-books',
          route: '/reads-books',
          icon: 'mdi-book-check-outline',
        },
        {
          title: 'Okunacaklar',
          value: 'book-to-reads',
          route: '/book-to-reads',
          icon: 'mdi-book-clock-outline',
        },
      ],
    },
    {
      title: 'Yazarlar',
      value: 'authors',
      route: '/authors',
      icon: 'mdi-account-edit',
    },
    {
      title: 'Arkadaş Yönet',
      value: 'friends',
      route: '/friends',
      icon: 'mdi-account-group',
    },
  ];

  // requiresRole varsa ve role_id uyuşmuyorsa filtrele
  return allItems.filter(item => {
    if (item.requiresRole) {
      return item.requiresRole.includes(role_id.value);
    }
    return true;
  });
})
const logout  = async()=>{
  logoutPopup.value = true
  await axios.post("/logout")
  logoutPopup.value = false
  await router.push("/login")
  localStorage.removeItem('user')
  localStorage.removeItem('token')
  store.$reset()
}

const fetchPendingRequestsCount = async () => {
  try {
    const res = await axios.get('/get-pending-requests-count')
    pendingRequestsCount.value = res.data.count || 0
  } catch (error) {
    console.error('Failed to fetch pending requests count:', error)
    pendingRequestsCount.value = 0
  }
}

onMounted(async ()=>{

  // const resp    = await axios.post("/user", {Data:localStorage.getItem('token')},{ withCredentials: true });
  const resAuth = await axios.get('/auth-check', { withCredentials: true });
  const user    = resAuth.data.user;

  role_id.value     = user.role_id;
  // likedPoems.value  = resp.data.admin.admin_liked_poems || []

  // bookmarksId.value = resp.data.admin.admin_bookmark_poems || []
  // readsBooks.value  = resp.data.admin.user_books_read || []

  store.setUser(user);
  // store.setLikedPoems(likedPoems.value)
  // store.setBookmarks(bookmarksId.value)

  // Fetch pending friend requests count
  await fetchPendingRequestsCount()

  // Connect to WebSocket
  const token = localStorage.getItem('token')
  if (token) {
    wsService.connect(token)

    // Listen for friend request notifications
    wsService.on('friend_request_received', (payload: any) => {
      console.log('[WebSocket] Friend request received:', payload)
      pendingRequestsCount.value = payload.count || 0
      console.log('[WebSocket] Badge updated to:', pendingRequestsCount.value)
    })

    wsService.on('friend_request_update', (payload: any) => {
      console.log('[WebSocket] Friend request update:', payload)
      pendingRequestsCount.value = payload.count || 0
      console.log('[WebSocket] Badge updated to:', pendingRequestsCount.value)
    })

    wsService.on('friend_request_accepted', (payload: any) => {
      console.log('[WebSocket] Friend request accepted by:', payload.username)
    })
  }
})

onUnmounted(() => {
  wsService.disconnect()
})

watch(group, () => {
  drawer.value = false
})

</script>

<template>
  <!-- App Bar / Header -->


  <!-- Navigation Drawer / Sidebar -->
  <v-navigation-drawer
    v-model="drawer"
    temporary
    class="sidebar-drawer"
  >
    <!-- Sidebar Header -->
    <div class="sidebar-header">
      <v-icon size="32" color="grey-lighten-1">mdi-feather</v-icon>
      <h3 class="sidebar-title">Mihrimâh</h3>
    </div>

    <v-divider class="sidebar-divider"></v-divider>

    <!-- Navigation List -->
    <v-list class="sidebar-list" nav>


      <!-- Main Menu Items -->
      <template v-for="(link, index) in items" :key="index">
        <!-- Items with sublinks -->
        <template v-if="link.subLinks">
          <v-list-group :value="link.value" class="sidebar-group">
            <template v-slot:activator="{ props }">
              <v-list-item
                v-bind="props"
                :prepend-icon="link.icon"
                :title="link.title"
                :value="link.value"
                class="sidebar-item"
              ></v-list-item>
            </template>
            <v-list-item
              v-for="(subLink, subIndex) in link.subLinks"
              :key="subIndex"
              active-class="active-item"
              class="sidebar-subitem"
              :prepend-icon="subLink.icon"
              :to="subLink.route"
              :title="subLink.title"
              :value="subLink.value"
            ></v-list-item>
          </v-list-group>
        </template>

        <!-- Items without sublinks -->
        <template v-else>
          <router-link
            active-class="active"
            class="text-decoration-none"
            :to="link.route"
          >
            <v-list-item
              :prepend-icon="link.icon"
              :title="link.title"
              :value="link.value"
              class="sidebar-item"
            >
              <template v-if="link.value === 'friends' && pendingRequestsCount > 0" v-slot:append>
                <v-badge
                  :content="pendingRequestsCount"
                  color="error"
                  inline
                />
              </template>
            </v-list-item>
          </router-link>
        </template>
      </template>

      <v-divider class="sidebar-divider my-2"></v-divider>

      <!-- Logout -->
      <v-list-item
        @click="logout"
        class="sidebar-item logout-item"
        title="Çıkış Yap"
        value="logout"
        prepend-icon="mdi-logout"
      />
    </v-list>
  </v-navigation-drawer>
  <!-- Logout Dialog -->
  <v-dialog
    v-model="logoutPopup"
    max-width="380"
    persistent
  >
    <v-card class="logout-dialog">
      <v-card-text class="d-flex align-center justify-center ga-3 pa-6">
        <v-progress-circular
          color="grey-lighten-1"
          indeterminate
          size="24"
          width="2"
        ></v-progress-circular>
        <span class="text-grey-lighten-1 text-body-1">Çıkış yapılıyor...</span>
      </v-card-text>
    </v-card>
  </v-dialog>

  <v-main class="pt-0">
    <v-app-bar class="app-header" id="header" elevation="0">
      <template v-slot:prepend>
        <v-app-bar-nav-icon
          variant="text"
          @click.stop="drawer = !drawer"
          class="nav-icon"
        ></v-app-bar-nav-icon>
      </template>

      <v-app-bar-title class="app-title">
        <router-link
          active-class="active"
          to="/home"
          class="text-decoration-none title-link"
        >
          <v-icon size="24" class="mr-2">mdi-feather</v-icon>
          Mihrimâh
        </router-link>
      </v-app-bar-title>

      <template v-slot:append>
        <v-btn
          @click="$router.push('/friends')"
          class="friend-notification-btn"
        >
          <v-badge
            :content="pendingRequestsCount"
            :model-value="pendingRequestsCount > 0"
            color="error"
            overlap
          >
            <v-icon>mdi-account-group</v-icon>
          </v-badge>
        </v-btn>

        <!-- Profile Avatar -->
        <v-btn
          v-if="store.user && store.user.username"
          @click="$router.push(`/profile/${store.user.username}`)"
          class="profile-avatar-btn"
          icon
        >
          <v-avatar size="36">
            <v-img
              :src="getImageUrl(store.user.profile_image)"
              cover
            />
          </v-avatar>
        </v-btn>
      </template>
    </v-app-bar>
    <router-view class="d-flex w-full px-6 py-8 mr-auto ml-auto content-area"/>
  </v-main>
  <!-- Footer Component -->
  <Footer />
</template>

<style scoped>
/* App Header */
.app-header {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%) !important;
  border-bottom: 1px solid #424242;
  backdrop-filter: blur(10px);
}

.nav-icon {
  transition: all 0.3s ease;
}

.nav-icon:hover {
  transform: rotate(90deg);
}

.app-title {
  font-weight: 600;
  letter-spacing: 0.5px;
}

.title-link {
  color: #ffffff !important;
  display: flex;
  align-items: center;
  transition: all 0.3s ease;
  font-size: 20px;
  font-weight: 600;
}

.title-link:hover {
  color: #e0e0e0 !important;
  transform: translateX(4px);
}

/* Sidebar Drawer */
.sidebar-drawer {
  background: linear-gradient(180deg, #1a1a1a 0%, #2d2d2d 100%) !important;
  border-right: 1px solid #424242;
}

/* Sidebar Header */
.sidebar-header {
  padding: 24px 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  background: linear-gradient(135deg, #212121 0%, #2d2d2d 100%);
}

.sidebar-title {
  font-size: 20px;
  font-weight: 600;
  color: #ffffff;
  letter-spacing: 0.5px;
  margin: 0;
}

.sidebar-divider {
  border-color: rgba(255, 255, 255, 0.1) !important;
  margin: 0 !important;
}

/* Sidebar List */
.sidebar-list {
  background: transparent !important;
  padding: 8px 0;
  overflow-y: auto !important;
  max-height: calc(100vh - 100px) !important;
}

.sidebar-list::-webkit-scrollbar {
  width: 6px;
}

.sidebar-list::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 3px;
}

.sidebar-list::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 3px;
}

.sidebar-list::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}

/* Sidebar Items */
.sidebar-item {
  margin: 4px 12px;
  border-radius: 8px;
  transition: all 0.3s ease;
  color: #e0e0e0 !important;
}

.sidebar-item :deep(.v-list-item__prepend) {
  margin-right: 12px !important;
}

.sidebar-item:hover {
  background: rgba(255, 255, 255, 0.05) !important;
  transform: translateX(4px);
}

.sidebar-subitem {
  margin: 2px 8px 2px 8px;
  padding-left: 12px !important;
  border-radius: 8px;
  transition: all 0.3s ease;
  color: #bdbdbd !important;
  font-size: 14px;
}

.sidebar-subitem :deep(.v-list-item__prepend) {
  margin-right: 0 !important;
  padding-left: 36px !important;
}

.sidebar-subitem :deep(.v-list-item-title) {
  padding-left: 0 !important;
}

.sidebar-subitem:hover {
  background: rgba(255, 255, 255, 0.05) !important;
  transform: translateX(4px);
  color: #e0e0e0 !important;
}

.active-item {
  background: linear-gradient(90deg, rgba(158, 158, 158, 0.2) 0%, rgba(158, 158, 158, 0.1) 100%) !important;
  border-left: 3px solid #9e9e9e;
  color: #ffffff !important;
}

/* Sidebar Group */
.sidebar-group {
  margin: 4px 0;
}

/* Logout Item */
.logout-item {
  color: #ef5350 !important;
}

.logout-item:hover {
  background: rgba(239, 83, 80, 0.1) !important;
}

/* Logout Dialog */
.logout-dialog {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  border-radius: 12px;
}

/* Button */
.button {
  text-transform: none;
}

/* Friend Notification Button */
.friend-notification-btn {
  color: #e0e0e0 !important;
  transition: all 0.3s ease;
}

.friend-notification-btn:hover {
  color: #ffffff !important;
  background: rgba(255, 255, 255, 0.05) !important;
}

/* Profile Avatar Button */
.profile-avatar-btn {
  margin-left: 8px;
  transition: all 0.3s ease;
}

.profile-avatar-btn:hover {
  transform: scale(1.1);
  background: rgba(255, 255, 255, 0.05) !important;
}

.profile-avatar-btn .v-avatar {
  border: 2px solid #424242;
  transition: all 0.3s ease;
}

.profile-avatar-btn:hover .v-avatar {
  border-color: #9e9e9e;
}



.content-area {
  flex: 1 0 auto;
  min-height: calc(100vh - 200px);
}

/* Main Layout - Ensures footer stays at bottom */
.v-main {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: transparent !important;
}

.v-main > :first-child {
  flex: 1 0 auto;
}

/* Responsive */
@media (max-width: 768px) {
  .title-link {
    font-size: 18px;
  }

  .sidebar-header {
    padding: 20px 16px;
  }

  .sidebar-title {
    font-size: 18px;
  }
}
</style>
