<script setup >
import {onMounted, reactive, ref} from "vue";
import axios from "axios";
import {useField, useForm} from "vee-validate";
import Table from "../../components/shared/Table/Table.vue"
import SnackbarNotification from "../../components/common/SnackbarNotification.vue"
import { useNotification } from "../../composables/useNotification"

const search = ref('')
const edit = ref(false)
const dialog = ref(false)
const selectedBookId = ref(0)
const books = reactive([])
const authors = ref([])
const { snackbar, error: showError, success } = useNotification()
const headers = [
  { key: 'id', title: 'ID' },
  { key: 'name', title: 'Adı' },
  { key: 'author', title: 'Yazar' },
  { key: 'image', title: 'Resim' },
  { key: 'page', title: 'Sayfa' },
  { key: 'community', title: 'İzin' },
  {
    align: 'start',
    key: 'islemler',
    sortable: false,
    title: 'İşlemler',
  }
]

const { handleSubmit, handleReset } = useForm({
  validationSchema: {
    name(value) {
      if (!value) return "Bu alan boş kalamaz."
      return true
    },
    author(value) {
      // Optional - kept for backward compatibility
      return true
    },
    author_id(value) {
      // Optional - can be null
      return true
    },
    created_at(value) {
      return true
    },
    page(value) {
      if (!value) return "Bu alan boş kalamaz."
      return true
    },
    image(value){
      if (!value) return "Bu alan boş kalamaz."
      return true
    },
    slug(value){
      return true
    },
    is_deleted(value) {
      return true
    },
    community(value) {
      if (!value) return "Bu alan boş kalamaz."
      return true
    }
  }
})

const name = useField('name')
const author = useField('author')
const author_id = useField('author_id')
const created_at = useField('created_at')
const page = useField('page')
const image = useField('image')
const slug = useField('slug')
const is_deleted = useField('is_deleted')
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
  selectedBookId.value = id
  try {
    await axios.get(`/get-book-by-id/${id}`)
      .then((res) => {
        name.value.value = res.data.name
        author.value.value = res.data.author
        author_id.value.value = res.data.author_id || null
        page.value.value = res.data.page
        image.value.value = res.data.image
        created_at.value.value = res.data.created_at
        community.value.value = res.data.community || 2
      }).catch((error) => {
        showError(`Kitap getirme başarısız: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  } catch (e) {
    console.log("hic girmedi")
    console.log(e)
  }
}
const deleteBook = async (id) => {
  try {
    await axios.delete(`/delete-book/${id}`)
      .then((res) => {
        console.log("res book data ", res.data)
        Object.assign(books, res.data)
      }).catch((error) => {
        showError(`Kitap silme başarısız: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  }catch (e) {
    console.log("hic girmedi")
    console.log(e)
  }
}
const submitCreateBook = handleSubmit(async (values) => {
  try {
    const data = {
      name: name.value.value,
      author: author.value.value || '',
      author_id: author_id.value.value || null,
      page: parseInt(page.value.value),
      image: image.value.value,
      community: community.value.value || 2
    }

    await axios.post('/create-book', data)
      .then((res) => {
        console.log(res.data)
        Object.assign(books, res.data)
        dialog.value = false
      }).catch((error) => {
        showError(`Kitap Oluşturma Başarısız: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  } catch (e) {
    console.log("hic girmedi")
    console.log(e)
  }
  await handleReset()
})
const submitEditBook = handleSubmit(async (values) => {
  try {
    const data = {
      name: name.value.value,
      author: author.value.value || '',
      author_id: author_id.value.value || null,
      page: parseInt(page.value.value),
      image: image.value.value,
      community: community.value.value || 2
    }

    await axios.put(`/update-book/${selectedBookId.value}`, data)
      .then((res) => {
        Object.assign(books, res.data)
        dialog.value = false
      }).catch((error) => {
        showError(`Kitap güncelleme başarısız: ${error.response?.data?.message || error.message}`)
        console.error(error)
      })
  } catch (e) {
    console.log("hic girmedi")
    console.log(e)
  }
  await handleReset()
})
onMounted(async () => {
  try {
    // Fetch books
    await axios.get('/get-books')
      .then((res) => {
        Object.assign(books, res.data)
        console.log(books)
      }).catch((error) => {
        showError(`kitap getirme basarisiz: ${error.response?.data?.message || error.message}`)
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

</script>

<template>
  <!-- TABLO -->
  <v-card
    title="Kitaplar"
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
              Kitap Ekle
            </v-btn>
          </template>
          <v-card>
            <v-card-title>
              <span class="text-h5">Kitap</span>
            </v-card-title>
            <v-card-text>
              <v-col cols="12">
                <form
                  class=" align-stretch d-flex flex-column gap-10"
                  @submit.prevent="submitCreateBook"
                  autocomplete="off"
                >
                  <v-row>
                    <v-col cols="12">
                      <v-text-field
                        v-model="name.value.value"
                        label="Kitap Adı"
                        :error-messages="name.errorMessage.value"
                        autocomplete="off"
                      />
                    </v-col>
                    <!-- Image -->
                    <VCol cols="12">
<!--                      <VFileInput-->
<!--                        v-if="!image.value.value"-->
<!--                        v-model="image.value.value"-->
<!--                        clearable-->
<!--                        show-size-->
<!--                        label="Resim"-->
<!--                        type="file"-->
<!--                        variant="outlined"-->
<!--                        accept=".jpg, .jpeg, .png, .gif"-->
<!--                      />-->
<!--                      <div-->
<!--                        v-else-->
<!--                        style="display: flex; align-items: center;"-->
<!--                      >-->
<!--                        <VBtn-->
<!--                          style="width: 280px; margin-top: 5px;"-->
<!--                          @click="image.value.value = ''"-->
<!--                        >-->
<!--                          <VIcon icon="bi-pencil-square" />-->
<!--                          <span>Resmi Değiştir</span>-->
<!--                        </VBtn>-->
<!--                      </div>-->
                      <v-text-field
                        v-model="image.value.value"
                        label="Resim"
                        :error-messages="image.errorMessage.value"
                        autocomplete="off"
                      />
                    </VCol>

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
                    <v-col cols="12">
                      <v-text-field
                        v-model="page.value.value"
                        label="Sayfa Sayısı"
                        :error-messages="page.errorMessage.value"
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
                        :error-messages="community.errorMessage.value"
                        variant="outlined"
                      />
                    </v-col>

                    <v-col cols="12" v-if="edit">
                      <v-text-field
                        readonly
                        v-model="created_at.value.value"
                        label="Oluşturulma Tarih"
                        autocomplete="off"
                      />
                    </v-col>

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
                          @click="submitCreateBook"
                        >
                          Kaydet
                        </v-btn>
                        <v-btn
                          v-else
                          style="text-transform: none;"
                          height="40"
                          class="rounded-lg mb-2"
                          prepend-icon="mdi-plus"
                          @click="submitEditBook"
                        >
                          Güncelle
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
      :items="books"
      :search="search"
      :deleteFunc="deleteBook"
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

