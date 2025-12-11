<script setup>
import { onMounted, reactive, ref} from "vue";
import axios from "axios";
import {useField, useForm} from "vee-validate";
import Table from "../../shared/Table/Table.vue"
import { QUILL_TOOLBAR_BASIC, QUILL_TOOLBAR_FULL } from "../../../constants/quillEditorConfig"
import SnackbarNotification from "../../common/SnackbarNotification.vue"
import { useNotification } from "../../../composables/useNotification"

const { snackbar, error: showError } = useNotification()

const search = ref('')
const edit = ref(false)
const dialog = ref(false)
const headers = [
  { key: 'id', title: 'ID', width: '80px' },
  { key: 'title', title: 'Başlık', width: '200px' },
  { key: 'content', title: 'İçerik', width: '400px' },
  {
    align: 'start',
    key: 'islemler',
    sortable: false,
    title: 'İşlemler',
    width: '150px'
  }
]

const mihrimahCards = reactive([])
const selectedCardId = ref(0)
const { handleSubmit, handleReset } = useForm({
  validationSchema: {
    title (value){
      if(!value || value.trim() === '' || value === '<p><br></p>') return "Bu alan boş kalamaz."
      return true
    },
    content (value){
      if(!value || value.trim() === '' || value === '<p><br></p>') return "Bu alan boş kalamaz."
      return true
    },
  },
})
const title = useField('title')
const content = useField('content')

const openCreateDialog = () => {
  edit.value = false
  dialog.value = true
}
const closeDialog = async () => {
  dialog.value = false
  await handleReset()
}
const openEditDialog = async (id) => {
  edit.value = true
  dialog.value = true
  selectedCardId.value = id
  try {
    await axios.get(`/get-mihrimah-card/${id}`)
      .then((res) => {
        title.value.value = res.data.title
        content.value.value = res.data.content
      }).catch((error) => {
        showError(`Mihrimah kart getirme başarısız: ${error.response?.data?.message || error.message}`)
      })
  }catch (e) {
    console.error(e)
  }
}
const submitCreateCard = handleSubmit(async (values) => {
  try {
    await axios.post('/create-mihrimah-card', values)
      .then((res) => {
        mihrimahCards.value = res.data
        dialog.value = false
      }).catch((error) => {
        showError(`Mihrimah kart oluşturma başarısız: ${error.response?.data?.message || error.message}`)
      })
  }catch (e) {
    console.error(e)
  }
  await handleReset()
})
const submitEditCard = handleSubmit(async (values) => {
  try {
    await axios.put(`/update-mihrimah-card/${selectedCardId.value}`, values)
      .then((res) => {
        mihrimahCards.value = res.data
        dialog.value = false
      }).catch((error) => {
        showError(`Mihrimah kart güncelleme başarısız: ${error.response?.data?.message || error.message}`)
      })
  }catch (e) {
    console.error(e)
  }
})
const deleteCard = async (id) => {
  try {
    await axios.delete(`/delete-mihrimah-card/${id}`)
      .then((res) => {
        mihrimahCards.value = res.data
      }).catch((error) => {
        showError(`Mihrimah kart silme başarısız: ${error.response?.data?.message || error.message}`)
      })
  }catch (e) {
    console.error(e)
  }
}
onMounted(async () => {
  try {
    await axios.get('/get-mihrimah-cards')
      .then((res) => {
        mihrimahCards.value = res.data
      }).catch((error) => {
        showError(`Mihrimah kartlar getirme başarısız: ${error.response?.data?.message || error.message}`)
      })
  }catch (e) {
    console.error(e)
  }
})

</script>

<template>
  <!-- TABLO -->
  <v-card
    title="Mihrimah Kart Yönetimi"
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
              Mihrimah Kart Ekle
            </v-btn>
          </template>
          <v-card>
            <v-card-title>
              <span class="text-h5">Mihrimah Kart</span>
            </v-card-title>
            <v-card-text>
              <v-col cols="12">
                <form
                  class=" align-stretch d-flex flex-column gap-10"
                  @submit.prevent="submitCreateCard"
                  autocomplete="off"
                >
                  <v-col cols="12">
                    <v-col cols="12">
                      <label class="text-subtitle-2 mb-2 d-block">Başlık (Buton Yazısı)</label>
                      <QuillEditor
                        v-model:content="title.value.value"
                        theme="snow"
                        contentType="html"
                        :toolbar="QUILL_TOOLBAR_BASIC"
                      />
                      <div v-if="title.errorMessage.value" class="text-error text-caption mt-1">
                        {{ title.errorMessage.value }}
                      </div>
                    </v-col>

                    <v-col cols="12">
                      <label class="text-subtitle-2 mb-2 d-block">İçerik (Popup İçeriği)</label>
                      <QuillEditor
                        v-model:content="content.value.value"
                        theme="snow"
                        contentType="html"
                        :toolbar="QUILL_TOOLBAR_FULL"
                        style="min-height: 200px;"
                      />
                      <div v-if="content.errorMessage.value" class="text-error text-caption mt-1">
                        {{ content.errorMessage.value }}
                      </div>
                    </v-col>

                    <v-col cols="12" >
                      <v-row class="d-flex justify-center items-center flex-nowrap">
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
                          @click="submitCreateCard"
                        >
                          Kaydet
                        </v-btn>
                        <v-btn
                          v-else
                          height="40"
                          class="rounded-lg mb-2"
                          prepend-icon="mdi-plus"
                          @click="submitEditCard"
                        >
                          Güncelle
                        </v-btn>
                      </v-row>
                    </v-col>
                  </v-col>
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
      :items="mihrimahCards"
      :search="search"
      :deleteFunc="deleteCard"
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
