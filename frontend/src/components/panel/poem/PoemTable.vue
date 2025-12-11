<script setup >
import {onMounted, reactive, ref, watch} from "vue";
import axios from "axios";
import {useField, useForm} from "vee-validate";
import Table from "../../shared/Table/Table.vue"
import SnackbarNotification from "../../common/SnackbarNotification.vue"
import { useNotification } from "../../../composables/useNotification"

const { snackbar, error: showError } = useNotification()

const search = ref('')
const edit = ref(false)
const dialog = ref(false)
const isOpenPicker = ref(false)
const selectedPoemId = ref(0)
const poems = reactive([])
const authors = ref([])
const headers = [
  { key: 'id', title: 'ID' },
  { key: 'title', title: 'Başlık' },
  { key: 'author', title: 'Yazar' },
  {
    align: 'start',
    key: 'islemler',
    sortable: false,
    title: 'İşlemler',
  }
]

const { handleSubmit, handleReset } = useForm({
  validationSchema: {
    title(value) {
      if (!value) return "Bu alan boş kalamaz."
      else if (value.length < 3) return "en az 3 karakter gir"
      else if (value.length > 20) return "en fazla 30 karakter gir"
      return true
    },
    author(value) {
      // Optional now - kept for backward compatibility
      return true
    },
    author_id(value) {
      // Optional - can be null
      return true
    },
    content(value) {
      if (!value) return "Bu alan boş kalamaz."
      return true
    },
    created_at(value) {
      return true
    },
    community(value) {
      return true
    },
  }
})

const title = useField('title')
const author = useField('author')
const author_id = useField('author_id')
const content = useField('content')
const created_at = useField('created_at')
const community = useField('community')

const openCreateDialog = async () => {
  edit.value = false
  await handleReset()
  dialog.value = true
  community.value.value = 2 // Default: Herkese Açık
}
const closeDialog =  async () => {
  dialog.value = false
  await handleReset()
}
const openEditDialog = async (id) => {
  edit.value = true
  dialog.value = true
  selectedPoemId.value = id
  try {
    await axios.get(`/get-poem-by-id/${id}`)
      .then((res) => {
        title.value.value = res.data.title
        author.value.value = res.data.author
        author_id.value.value = res.data.author_id || null
        content.value.value = res.data.content
        created_at.value.value = res.data.created_at
        community.value.value = res.data.community || 2
      }).catch((error) => {
        showError(`Şiir getirme başarısız: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  }catch (e) {
    console.log("hic girmedi")
    console.log(e)
  }
}
const deletePoem = async (id) => {
  try {
    await axios.delete(`/delete-poem/${id}`)
      .then((res) => {
        console.log("res poem data ", res.data)
        poems.value = res.data
      }).catch((error) => {
        showError(`Şiir silme başarısız: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  }catch (e) {
    console.log("hic girmedi")
    console.log(e)
  }
}
const submitCreatePoem = handleSubmit(async (values) => {
  try {
    await axios.post('/create-poem', values)
      .then((res) => {
        console.log(res.data)
        poems.value = res.data
        dialog.value = false
      }).catch((error) => {
        showError(`Şiir oluşturma başarısız: ${error.response?.data?.message || error.message}`)
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
    await axios.put(`/update-poem/${selectedPoemId.value}`, values)
      .then((res) => {
        poems.value = res.data
        dialog.value = false
      }).catch((error) => {
        showError(`Şiir güncelleme başarısız: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  }catch (e) {
    console.log("hic girmedi")
    console.log(e)
  }
})
onMounted(async () => {
  try {
    // Fetch poems
    await axios.get('/get-all-poems')
      .then((res) => {
        poems.value = res.data
        console.log(poems.value)
      }).catch((error) => {
        showError(`Şiir getirme başarısız: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })

    // Fetch authors for dropdown
    await axios.get('/get-all-authors-dropdown')
      .then((res) => {
        authors.value = res.data
      }).catch((error) => {
        console.error('Failed to fetch authors:', error)
      })
  }catch (e) {
    console.log("hic girmedi")
    console.log(e)
  }
})
watch(created_at.value, (value) => {
  console.log("created_at.value.value", value)
})
</script>

<template>
  <!-- TABLO -->
  <v-card
    title="Şiirler"
    class="mt-8 sm:w-100 md:w-75 mx-auto rounded-xl overflow-y-auto bg-grey-darken-2"
  >
    <!-- ARAMA KISMI   -->
    <template v-slot:text>
      <div class="d-flex flex-column md:flex-row justify-center align-start mb-0">
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
              Şiir Ekle
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
                    <v-col cols="12">
                      <v-text-field
                        v-model="title.value.value"
                        label="Şiir Başlığı"
                        :error-messages="title.errorMessage.value"
                        autocomplete="off"
                      />
                    </v-col>
                    <v-col cols="12">
                      <v-autocomplete
                        v-model="author_id.value.value"
                        :items="authors"
                        item-title="name"
                        item-value="id"
                        label="Yazar Seçin"
                        :error-messages="author.errorMessage.value"
                        clearable
                        no-data-text="Yazar bulunamadı"
                      >
                        <template v-slot:append>
                          <v-icon>mdi-account</v-icon>
                        </template>
                      </v-autocomplete>
                    </v-col>
                    <v-col cols="12" v-if="edit">
                      <v-text-field
                        readonly
                        v-model="created_at.value.value"
                        label="Oluşturulma Tarih"
                        autocomplete="off"
                      />
                    </v-col>
                    <v-col cols="12">
                      <v-select
                        v-model="community.value.value"
                        :items="[
                          { title: 'Özel (Sadece Admin ve Üyeler)', value: 1 },
                          { title: 'Herkese Açık', value: 2 }
                        ]"
                        label="Görünürlük"
                        item-title="title"
                        item-value="value"
                        variant="outlined"
                      />
                    </v-col>
                    <VCol
                      cols="12"
                      class="h-100 border rounded-lg px-3 ma-3"
                    >
                      <VCardText
                        readonly
                        class="pa-0 mb-3 font-weight-medium"
                      >
                        Açıklama
                      </VCardText>
                      <QuillEditor
                        v-model:content="content.value.value"
                        content-type="html"
                        theme="snow"
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
      :items="poems"
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

