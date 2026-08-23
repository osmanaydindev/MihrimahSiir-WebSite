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
  <td class="table-actions">
    <v-dialog width="500">
      <template v-slot:activator="{ props }">
        <v-btn
          style="color: #E57373;"
          v-bind="props"
          variant="text"
          icon
          size="small"
        >
          <v-icon size="large">mdi-delete</v-icon>
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

    <v-btn
        v-if="route.path!=='/log-management'"
        color="success"
        variant="text"
        icon
        size="small"
        @click="openEditDialog(id)"
    >
      <v-icon size="large">mdi-note-edit</v-icon>
      <v-tooltip activator="parent" location="top">Düzenle</v-tooltip>
    </v-btn>
  </td>
</template>

<style scoped>
.table-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 96px;
}
</style>
