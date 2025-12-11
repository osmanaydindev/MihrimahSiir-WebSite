<script setup lang="ts">
const props = defineProps({
  title: String,
  subtitle: String,
  icon: String
})
</script>

<template>
  <v-card
      class="w-75 mx-auto rounded-xl"
      bg-opacity="10"
      width="600"
      :prepend-icon="icon"
  >
    <template v-slot:title>
      {{title}}
    </template>

    <v-card-text>
      {{subtitle}}
    </v-card-text>
  </v-card>
</template>

<style scoped>

</style>
