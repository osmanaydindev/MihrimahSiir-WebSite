<script setup>
import {useRoute} from "vue-router";
const route = useRoute()
// console.log("route", route)
const props = defineProps({
  id: Number,
  deleteFunc: Function,
  openEditDialog: Function
})
</script>

<template>
  <td class="d-flex justify-start align-center">
    <v-dialog width="500">
      <template v-slot:activator="{ props }">
        <v-btn
          style="color: #E57373;"
          v-bind="props"
          variant="text"
        >
          <v-icon size="x-large" left>mdi-delete</v-icon>
          <v-tooltip activator="parent" location="top">Sil</v-tooltip>
        </v-btn>
      </template>

      <template v-slot:default="{ isActive }">
        <v-card
          class="px-4 py-8 d-flex flex-column justify-center align-center"
        >
          <v-card-title>Silme işlemini onayalıyor musunuz?</v-card-title>
          <v-card-actions class="d-flex justify-end mt-2">
            <v-btn
              text="Sil"
              class="rounded-full"
              color="error"
              @click="deleteFunc(id); isActive.value = false;"
            />
            <v-btn
              text="Vazgeç"
              @click="isActive.value = false"
            />
          </v-card-actions>
        </v-card>
      </template>
    </v-dialog>
<!--    <button
        style="color: #E57373;"
        @click="deleteFunc(id)"
    >
      <v-icon size="x-large" left>mdi-delete</v-icon>
      <v-tooltip activator="parent" location="top">Sil</v-tooltip>
    </button>-->

    <button
        v-if="route.path!=='/log-management'"
        style="color: #A5D6A7"
        @click="openEditDialog(id)"
    >
      <v-icon size="x-large" left>mdi-note-edit</v-icon>
      <v-tooltip activator="parent" location="top">Düzenle</v-tooltip>
    </button>
  </td>
</template>

<style scoped>

</style>
