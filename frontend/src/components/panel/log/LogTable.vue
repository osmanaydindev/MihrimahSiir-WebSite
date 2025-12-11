<script setup >
import { onMounted, reactive, ref} from "vue";
import axios from "axios";
import Table from "../../shared/Table/Table.vue";
import SnackbarNotification from "../../common/SnackbarNotification.vue"
import { useNotification } from "../../../composables/useNotification"


const headers = [
  { key: 'id', title: 'ID' },
  { key: 'username', title: 'Kullanıcı Adı' },
  { key: 'login_date', title: 'Giriş Tarihi'},
  { key: 'role_id', title:'Role'},
  { key: 'actions', title: 'İşlemler', sortable: false },
]
const roles = [
  { id: "1", name: 'Admin' },
  { id: "2", name: 'Kullanıcı' },
]

const logs = reactive([])
const search = ref('')
const { snackbar, error: showError, success } = useNotification()
const deleteLog = async (id) => {
  try {
    await axios.delete(`/delete-log/${id}`)
      .then((res) => {
        logs.value = res.data
        console.log("log silme basarili")
      }).catch((error) => {
        showError(`log silme basarisiz: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  }catch (e) {
    console.log(e)
  }
}
const openEditDialog = () => {
  console.log('openEditDialog')
}
onMounted(async () => {
  try {
    await axios.get('/get-logs')
      .then((res) => {
        logs.value = res.data
        console.log("logs", logs.value)
      }).catch((error) => {
        showError(`admin getirme basarisiz: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  }catch (e) {
    console.log(e)
  }
})

</script>

<template>
  <!-- TABLO -->
  <v-card
    title="Giriş Logları"
    class="mt-8 sm:w-100 md:w-75 mx-auto rounded-xl overflow-y-auto bg-grey-darken-2"
  >
    <!-- ARAMA KISMI   -->
    <template v-slot:text>

    </template>
    <!--  TABLO ICERIK  -->
    <Table
      :headers="headers"
      :items="logs"
      :search="search"
      :deleteFunc="deleteLog"
      :openEditDialog="openEditDialog"
    />
  </v-card>

  <!-- Snackbar for notifications -->
  <SnackbarNotification
    v-model="snackbar.show"
    :message="snackbar.message"
    :color="snackbar.color"
  />
</template>

