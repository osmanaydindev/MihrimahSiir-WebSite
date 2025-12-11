<script setup lang="ts">
const props = defineProps(['comments']);
</script>

<template>
  <v-card
    v-if="!comments || comments.length === 0"
    class="text-center pa-6 rounded-xl"
    variant="outlined"
  >
    <v-icon
      icon="mdi-comment-outline"
      size="48"
      class="mb-2 text-grey"
    ></v-icon>
    <div class="text-body-1 text-grey">Henüz yorum yapılmamış</div>
  </v-card>
</template>

<style scoped>

</style>
