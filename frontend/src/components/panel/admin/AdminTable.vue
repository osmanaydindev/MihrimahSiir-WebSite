<script setup >
import { onMounted, reactive, ref} from "vue";
import axios from "axios";
import {useField, useForm} from "vee-validate";
import Table from "../../shared/Table/Table.vue"
import SnackbarNotification from "../../common/SnackbarNotification.vue"
import { useNotification } from "../../../composables/useNotification"

const search = ref('')
const edit = ref(false)
const dialog = ref(false)
const visible = ref(false)
const { snackbar, error: showError, success } = useNotification()
const headers = [
  { key: 'id', title: 'ID' },
  { key: 'username', title: 'Kullanıcı Adı' },
  { key: 'email', title: 'Email' },
  { key: 'role_id', title: 'Rol' },
  {
    align: 'start',
    key: 'islemler',
    sortable: false,
    title: 'İşlemler',
  }
]
const roles = [
  { id: "1", name: 'Admin' },
  { id: "2", name: 'Kullanıcı' },
]

const admins = reactive([])
const selectedAdminId = ref(0)
const { handleSubmit, handleReset } = useForm({
  validationSchema: {
    username (value){
      if(!value) return "Bu alan boş kalamaz."
      else if (value.length<3) return "en az 3 karakter gir"
      else if (value.length>20) return "en fazla 30 karakter gir"
      return true
    },
    email (value){
      if(!value) return "Bu alan boş kalamaz."
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if(!emailRegex.test(value)) return "Geçerli bir email adresi girin"
      return true
    },
    password (value){
      if(!value) return "Bu alan boş kalamaz."
      return true
    },
    role_id (value) {
      if (!value) return "Bu alan boş kalamaz."
      return true
    },
  },
})
const username = useField('username')
const email = useField('email')
const password = useField('password')
const role_id = useField('role_id')
const openCreateDialog = () => {
  edit.value = false
  dialog.value = true
}
const closeDialog =  async () => {
  dialog.value = false
  await handleReset()
}
const openEditDialog = async (id) => {
  edit.value = true
  dialog.value = true
  selectedAdminId.value = id
  // console.log("id =>", id)
  try {
    await axios.get(`/get-admin/${id}`)
      .then((res) => {
        username.value.value = res.data.username
        email.value.value = res.data.email
        role_id.value.value = res.data.role_id
        // console.log("res admin data ", res.data)
      }).catch((error) => {
        showError(`User getirme basarisiz: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  }catch (e) {
  }
}
const submitCreateAdmin = handleSubmit(async (values) => {
  try {
    await axios.post('/create-admin', values)
      .then((res) => {
        // console.log(res.data)
        admins.value = res.data
        dialog.value = false
      }).catch((error) => {
        showError(`kullanıcı olusturma basarisiz: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  }catch (e) {
    // console.log("hic girmedi")
  }
  await handleReset()
})
const submitEditAdmin = handleSubmit(async (values) => {
  try {
    await axios.put(`/update-admin/${selectedAdminId.value}`, values)
      .then((res) => {
        admins.value = res.data
        dialog.value = false
      }).catch((error) => {
        showError(`kullanıcı guncelleme basarisiz: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  }catch (e) {
  }
})
const deleteAdmin = async (id) => {
  try {
    await axios.delete(`/delete-admin/${id}`)
      .then((res) => {
        admins.value = res.data
      }).catch((error) => {
        showError(`Ülke silme basarisiz: ${error.response?.data?.message || error.message}`)
      })
  }catch (e) {
  }
}
onMounted(async () => {
  try {
    await axios.get('/get-admins-management')
      .then((res) => {
        admins.value = res.data
      }).catch((error) => {
        showError(`admin getirme basarisiz: ${error.response?.data?.message || error.message}`)
      })
  }catch (e) {
  }
})

</script>

<template>
  <!-- TABLO -->
  <v-card
    title="Kullanıcılar"
    class="mt-8 sm:w-100 md:w-75 mx-auto rounded-xl overflow-y-auto bg-grey-darken-2"
  >
    <!-- ARAMA KISMI   -->
    <template v-slot:text>
      <div class="d-flex flex-column sm:flex-row justify-center align-start mb-0">
        <v-text-field
          v-model="search"
          variant="solo"
          prepend-inner-icon="mdi-magnify"
          class="rounded-none overflow-hidden w-100 sm:w-[25%]"
          label="Ara"
          clearable
        ></v-text-field>

        <v-dialog
          v-model="dialog"
          persistent
          width="1024"
          style="height: auto"
        >
          <template v-slot:activator="{ props }">
            <v-btn
              v-bind="props"
              variant="flat"
              prepend-icon="mdi-plus-circle"
              size="large"
              height="56"
              class="sm:mx-2 text-none"
              @click="openCreateDialog"
            >
              Kullanıcı Ekle
            </v-btn>
          </template>
          <v-card>
            <v-card-title>
              <span class="text-h5">Kullanıcı</span>
            </v-card-title>
            <v-card-text>
              <v-col cols="12">
                <form
                  class=" align-stretch d-flex flex-column gap-10"
                  @submit.prevent="submitCreateAdmin"
                  autocomplete="off"
                >
                  <v-row>
                    <v-col cols="12">
                      <v-text-field
                        v-model="username.value.value"
                        label="Kullanıcı Adı"
                        :error-messages="username.errorMessage.value"
                        autocomplete="off"
                      />
                    </v-col>

                    <v-col cols="12">
                      <v-text-field
                        v-model="email.value.value"
                        label="Email"
                        :error-messages="email.errorMessage.value"
                        autocomplete="off"
                        type="email"
                      />
                    </v-col>

                    <v-col>
                      <VSelect
                        v-model="role_id.value.value"
                        :items="roles"
                        item-title="name"
                        item-value="id"
                        label="Rol"
                        clearable
                      />
                    </v-col>

                    <v-col cols="12">
                      <v-text-field
                        :append-inner-icon="visible ? 'mdi-eye-off' : 'mdi-eye'"
                        :type="visible ? 'text' : 'password'"
                        placeholder="Şifre"
                        prepend-inner-icon="mdi-lock-outline"
                        variant="outlined"
                        v-model="password.value.value"
                        :error-messages="password.errorMessage.value"
                        @click:append-inner="visible = !visible"
                      />
                    </v-col>



                    <v-col cols="12" >
                      <v-row class="justify-end">
                        <v-btn
                          height="40"
                          class="rounded-lg mr-4"
                          prepend-icon="mdi-close"
                          @click="closeDialog()"
                        >
                          Kapat
                        </v-btn>
                        <v-btn
                          v-if="!edit"
                          height="40"
                          class="rounded-lg mb-2"
                          prepend-icon="mdi-plus"
                          @click="submitCreateAdmin"
                        >
                          Kaydet
                        </v-btn>
                        <v-btn
                          v-else
                          height="40"
                          class="rounded-lg mb-2"
                          prepend-icon="mdi-plus"
                          @click="submitEditAdmin"
                        >
                          Guncelle
                        </v-btn>
                      </v-row>
                    </v-col>
                  </v-row>
                </form>
              </v-col>
            </v-card-text>
          </v-card>
        </v-dialog>
      </div>
    </template>
    <!--  TABLO ICERIK  -->
    <Table
      :headers="headers"
      :items="admins"
      :search="search"
      :deleteFunc="deleteAdmin"
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

