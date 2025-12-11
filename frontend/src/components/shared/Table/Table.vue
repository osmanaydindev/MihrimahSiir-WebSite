<script setup lang="js">
import {useRoute} from "vue-router";
import {ref, onMounted, computed} from "vue";
import DeleteAndUpdate from "@/components/shared/Table/DeleteAndUpdate.vue";
import {useAppStore} from "@/store/app";

const store = useAppStore()

const props = defineProps({
  headers:Array,
  items:Object,
  search:String,
  deleteFunc:Function,
  openEditDialog:Function,
})

// Handle both ref and reactive items
const tableItems = computed(() => {
  // If items has .value property (ref), use it
  if (props.items && 'value' in props.items) {
    return props.items.value
  }
  // Otherwise, items is already the array (reactive or plain object)
  return props.items
})
const route         = useRoute()
const pageName      = ref("")
const whichTable    = ()=>{
  return route.fullPath
}
onMounted(()=>{
  pageName.value = whichTable()
})
const getRoleName = (role) => {
  switch (role) {
    case 1:
      return "Admin";
    case 2:
      return "Kullanıcı";
    case 3:
      return "Misafir";
    default:
      return "Bilinmeyen Rol";
  }
}
const imagePath = (image) => {
  return 'http://127.0.0.1:8080/uploads/'+image
}
const stripHtml = (html) => {
  if (!html) return ''
  const tmp = document.createElement('div')
  tmp.innerHTML = html
  return tmp.textContent || tmp.innerText || ''
}
</script>

<template>
  <!--  TABLO ICERIK  -->
  <v-data-table
      :headers="headers"
      :items="tableItems"
      :search="search"
      class="elevation-1 bg-grey-darken-1"
      style="border-radius: 0;"
  >
    <template v-slot:item="row" >
      <tr class="bg-white">
        <td>{{row.item.id}}</td>

        <td v-if="pageName==='/poem-management'">{{row.item.title}}</td>
        <td v-if="pageName==='/poem-management'">{{row.item.author_data?.name || row.item.author}}</td>


        <td v-if="pageName==='/admin-management'">{{row.item.username}}</td>
        <td v-if="pageName==='/admin-management'">{{row.item.email}}</td>
        <td v-if="pageName==='/admin-management'">{{getRoleName(row.item.role_id)}}</td>

        <td v-if="pageName==='/log-management'">{{row.item.username}}</td>
        <td v-if="pageName==='/log-management'">{{row.item.login_date}}</td>
        <td v-if="pageName==='/log-management'">{{getRoleName(row.item.role_id)}}</td>

        <td v-if="pageName==='/book-management'">{{row.item.name}}</td>
        <td v-if="pageName==='/book-management'">{{row.item.author_data?.name || row.item.author}}</td>
        <td v-if="pageName==='/book-management'">
          <v-img
            style="width: 35px; height: 35px;"
            :src="imagePath(row.item.image)"
          ></v-img>
        </td>
        <td v-if="pageName==='/book-management'">{{row.item.page}}</td>
        <td v-if="pageName==='/book-management'">{{getRoleName(row.item.community)}}</td>

        <td v-if="pageName==='/author-management'">{{row.item.name}}</td>
        <td v-if="pageName==='/author-management'">{{row.item.birth_year || '-'}}</td>
        <td v-if="pageName==='/author-management'">{{row.item.death_year || '-'}}</td>
        <td v-if="pageName==='/author-management'">{{row.item.nationality || '-'}}</td>

        <td v-if="pageName==='/reminder-management'" v-html="row.item.text" ></td>
        <td v-if="pageName==='/reminder-management'">{{getRoleName(row.item.permission)}}</td>

        <td v-if="pageName==='/homepage-management'">{{getRoleName(row.item.permission)}}</td>

        <td v-if="pageName==='/mihrimah-card-management'">{{stripHtml(row.item.title)}}</td>
        <td v-if="pageName==='/mihrimah-card-management'">{{stripHtml(row.item.content)}}</td>






        <delete-and-update
            :id="row.item.id"
            :deleteFunc="deleteFunc"
            :open-edit-dialog="openEditDialog"
        />
      </tr>
    </template>
  </v-data-table>
</template>

<style scoped>

</style>
