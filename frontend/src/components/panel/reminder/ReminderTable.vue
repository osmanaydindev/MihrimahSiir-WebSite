<script setup >
import {onMounted, reactive, ref} from "vue";
import axios from "axios";
import {useField, useForm} from "vee-validate";
import Table from "../../shared/Table/Table.vue"
import SnackbarNotification from "../../common/SnackbarNotification.vue"
import { useNotification } from "../../../composables/useNotification"

const search = ref('')
const edit = ref(false)
const dialog = ref(false)
const isOpenPicker = ref(false)
const selectedReminderId = ref(0)
const reminders = reactive([])
const { snackbar, error: showError, success } = useNotification()
const headers = [
  { key: 'id', title: 'ID' },
  { key: 'text', title: 'Metin' },
  { key: 'permission', title: 'İzin' },
  {
    align: 'start',
    key: 'islemler',
    sortable: false,
    title: 'İşlemler',
  }
]

const permissions = [
  { id: 1, name: 'Admin' },
  { id: 2, name: 'Kullanıcı' },
  { id: 3, name: 'Misafir' },
]

const { handleSubmit, handleReset } = useForm({
  validationSchema: {
    text(value) {
      if (!value) return "Bu alan boş kalamaz."
      return true
    },
    permission(value) {
      if (!value) return "Bu alan boş kalamaz."
      return true
    },
  }
})

const text = useField('text')
const permission = useField('permission')

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
  selectedReminderId.value = id
  try {
    await axios.get(`/reminders/${id}`)
      .then((res) => {
        console.log("res->",res)
        text.value.value = res.data.text
        permission.value.value = res.data.permission
        console.log(res.data.content)

      }).catch((error) => {
        showError(`Poem getirme basarisiz: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  }catch (e) {
    console.log("hic girmedi")
    console.log(e)
  }
}
const deletePoem = async (id) => {
  try {
    await axios.delete(`/reminders/${id}`)
      .then((res) => {
        console.log("res reminders data ", res.data)
        reminders.value = res.data
      }).catch((error) => {
        showError(`reminders silme basarisiz: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  }catch (e) {
    console.log("hic girmedi")
    console.log(e)
  }
}
const submitCreatePoem = handleSubmit(async (values) => {
  try {
    await axios.post('/reminders', values)
      .then((res) => {
        console.log(res.data)
        reminders.value = res.data
        dialog.value = false
      }).catch((error) => {
        showError(`reminders olusturma basarisiz: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  }catch (e) {
    console.log("hic girmedi")
    console.log(e)
  }
  await handleReset()
})
const submitEditPoem = handleSubmit(async (values) => {
  try {
    await axios.put(`/reminders/${selectedReminderId.value}`, values)
      .then((res) => {
        reminders.value = res.data
        dialog.value = false
      }).catch((error) => {
        showError(`reminders guncelleme basarisiz: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  }catch (e) {
    console.log("hic girmedi")
    console.log(e)
  }
})
onMounted(async () => {
  try {
    await axios.get('/reminders')
      .then((res) => {
        reminders.value = res.data
        console.log(reminders.value)
      }).catch((error) => {
        showError(`reminders getirme basarisiz: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  }catch (e) {
    console.log("hic girmedi")
    console.log(e)
  }
})
</script>

<template>
  <!-- TABLO -->
  <v-card
    title="Hatırlatıcılar"
    class="mt-8 sm:w-100 md:w-75 mx-auto rounded-xl overflow-y-auto bg-grey-darken-2"
  >
    <!-- ARAMA KISMI   -->
    <template v-slot:text>
      <div class="d-flex md:flex-row flex-column justify-center align-start mb-0">
        <v-text-field
          v-model="search"
          variant="solo"
          prepend-inner-icon="mdi-magnify"
          class="rounded-none overflow-hidden w-100 md:w-[25%]"
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
              style="text-transform: none;"
              v-bind="props"
              variant="flat"
              prepend-icon="mdi-plus-circle"
              size="large"
              height="56"
              class="md:mx-2 text-none"
              @click="openCreateDialog"
            >
            reminders Ekle
            </v-btn>
          </template>
          <v-card>
            <v-card-title>
              <span class="text-h5">Şiir</span>
            </v-card-title>
            <v-card-text>
              <v-col cols="12">
                <form
                  class=" align-stretch d-flex flex-column gap-10"
                  @submit.prevent="submitCreatePoem"
                  autocomplete="off"
                >
                  <v-row>
                    <VCol
                      cols="12"
                      class="h-100 border rounded-lg px-3 ma-3"
                    >
                      <VCardText
                        readonly
                        class="pa-0 mb-3 font-weight-medium"
                      >
                        Metin
                      </VCardText>
                      <QuillEditor
                        v-model:content="text.value.value"
                        content-type="html"
                        theme="snow"
                      />
                    </VCol>

                    <VCol cols="12" class="px-3 ma-3">
                      <VSelect
                        v-model="permission.value.value"
                        :items="permissions"
                        item-title="name"
                        item-value="id"
                        label="İzin Seviyesi"
                        :error-messages="permission.errorMessage.value"
                        clearable
                      />
                    </VCol>

                    <v-col cols="12" >
                      <v-row class="justify-end">
                        <v-btn
                          style="text-transform: none;"
                          height="40"
                          class="rounded-lg mr-4"
                          prepend-icon="mdi-close"
                          @click="closeDialog()"
                        >
                          Kapat
                        </v-btn>
                        <v-btn
                          v-if="!edit"
                          style="text-transform: none;"
                          height="40"
                          class="rounded-lg mb-2"
                          prepend-icon="mdi-plus"
                          @click="submitCreatePoem"
                        >
                          Kaydet
                        </v-btn>
                        <v-btn
                          v-else
                          style="text-transform: none;"
                          height="40"
                          class="rounded-lg mb-2"
                          prepend-icon="mdi-plus"
                          @click="submitEditPoem"
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
      :items="reminders"
      :search="search"
      :deleteFunc="deletePoem"
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

