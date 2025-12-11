<script setup lang="ts">
import {useRouter, useRoute} from "vue-router";
const router = useRouter()
const route = useRoute()
</script>

<template>
  <div class="mt-8 pb-16 d-flex flex-column align-center align-content-start flex-sm-row flex-sm-wrap justify-sm-center ga-5">
    <v-btn
      v-if="route.path==='/latest-poems'"
      class="button"
      variant="plain"
      @click="router.push('/poems')"
    >Tüm Şiirleri Gör
    </v-btn>
    <v-btn
      class="button"
      variant="plain"
      @click="router.push('/liked-poems')"
    >Beğenilen Şiirleri Gör
    </v-btn>
    <v-btn
      class="button"
      variant="plain"
      @click="router.push('/bookmark')"
    >Kaydedilen Şiirleri Gör
    </v-btn>
  </div>
</template>

<style scoped>

</style>
