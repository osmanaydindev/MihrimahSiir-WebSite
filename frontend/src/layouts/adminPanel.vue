<script lang="ts" setup>
import axios from "axios";
import { ref } from 'vue'
import {useRouter} from "vue-router"
import {useAppStore} from "../store/app";
const router = useRouter()
const store = useAppStore()
const logout = async()=>{
  await axios.post("/logout")
  await router.push("/login")
  localStorage.removeItem('user')
  localStorage.removeItem('token')
  store.$reset()
}
const drawer = ref(true)
const rail = ref(true)
</script>

<template>
  <v-app>
      <v-card class="h-100">
        <v-layout class="h-100">
          <v-navigation-drawer
            v-model="drawer"
            :rail="rail"
            permanent
            @click="rail = false"
            class="h-100"
          >
            <v-list-item
              prepend-icon="mdi-account-circle"
              @click="router.push('/panel')"
              nav
            >
              <template v-slot:append>
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
              <router-link to="/home" tag="v-list-item" class="text-white text-decoration-none">
                <v-list-item prepend-icon="mdi-home-switch-outline" title="Sayfaya Git" value="home"></v-list-item>
              </router-link>
              <router-link to="/log-management" tag="v-list-item" class="text-white text-decoration-none">
                <v-list-item prepend-icon="mdi-information" title="Giriş Yönet" value="latest-poems"></v-list-item>
              </router-link>
              <router-link to="/admin-management" tag="v-list-item" class="text-white text-decoration-none">
                <v-list-item prepend-icon="mdi-account-group-outline" title="Admin Yönet" value="admins"></v-list-item>
              </router-link>
              <router-link to="/poem-management" tag="v-list-item" class="text-white text-decoration-none">
                <v-list-item prepend-icon="mdi-text-box-plus-outline" title="Şiir Yönet" value="poems"></v-list-item>
              </router-link>
              <router-link to="/book-management" tag="v-list-item" class="text-white text-decoration-none">
                <v-list-item prepend-icon="mdi-book" title="Kitap Yönet" value="books"></v-list-item>
              </router-link>
              <router-link to="/author-management" tag="v-list-item" class="text-white text-decoration-none">
                <v-list-item prepend-icon="mdi-account-edit" title="Yazar Yönet" value="authors"></v-list-item>
              </router-link>
              <router-link to="/reminder-management" tag="v-list-item" class="text-white text-decoration-none">
                <v-list-item prepend-icon="mdi-bell" title="Hatırlatıcı Yönet" value="reminders"></v-list-item>
              </router-link>
              <router-link to="/homepage-management" tag="v-list-item" class="text-white text-decoration-none">
                <v-list-item prepend-icon="mdi-home-edit" title="Anasayfa Yönet" value="homepage"></v-list-item>
              </router-link>
              <router-link to="/mihrimah-card-management" tag="v-list-item" class="text-white text-decoration-none">
                <v-list-item prepend-icon="mdi-card-text" title="Mihrimah Kart Yönet" value="mihrimah-card"></v-list-item>
              </router-link>
              <v-list-item
                prepend-icon="mdi-logout"
                @click="logout"
                title="Çıkış Yap"
              />
            </v-list>
          </v-navigation-drawer>
          <v-main class="h-100" style="height: 250px">
            <div class="main d-flex w-75 justify-center my-12 mr-auto ml-auto">
              <router-view />
            </div>
          </v-main>
        </v-layout>
      </v-card>


  </v-app>
</template>

