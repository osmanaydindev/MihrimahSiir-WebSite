<script setup>
import { onMounted, ref } from "vue";
import axios from "axios";
import { useField, useForm } from "vee-validate";
import Table from "../../components/shared/Table/Table.vue"
import SnackbarNotification from "../../components/common/SnackbarNotification.vue"
import { useNotification } from "../../composables/useNotification"

const search = ref('')
const edit = ref(false)
const dialog = ref(false)
const selectedAuthorId = ref(0)
const authors = ref([])
const { snackbar, error: showError, success } = useNotification()
const headers = [
  { key: 'id', title: 'ID' },
  { key: 'name', title: 'Adı' },
  { key: 'birth_year', title: 'Doğum' },
  { key: 'death_year', title: 'Ölüm' },
  { key: 'nationality', title: 'Uyruk' },
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
    bio() {
      return true
    },
    birth_year(value) {
      if (value && !/^[0-9]+$/.test(value)) return "Sadece sayı giriniz."
      return true
    },
    death_year(value) {
      if (value && !/^[0-9]+$/.test(value)) return "Sadece sayı giriniz."
      return true
    },
    nationality() {
      return true
    },
    image() {
      return true
    }
  }
})

const name = useField('name')
const bio = useField('bio')
const birth_year = useField('birth_year')
const death_year = useField('death_year')
const nationality = useField('nationality')
const image = useField('image')

const openCreateDialog = () => {
  edit.value = false
  dialog.value = true
}

const closeDialog = async () => {
  dialog.value = false
  await handleReset()
}

const openEditDialog = async (authorOrId) => {
  try {
    // If authorOrId is a number (ID), use it directly. Otherwise, get ID from object
    const authorId = typeof authorOrId === 'number' ? authorOrId : authorOrId.id

    // Fetch full author details from backend
    const response = await axios.get(`/get-author-by-id/${authorId}`)
    const fullAuthor = response.data

    // console.log('Full author data from backend:', fullAuthor)

    edit.value = true
    selectedAuthorId.value = fullAuthor.id
    name.value.value = fullAuthor.name || ''
    bio.value.value = fullAuthor.bio || ''
    birth_year.value.value = fullAuthor.birth_year || ''
    death_year.value.value = fullAuthor.death_year || ''
    nationality.value.value = fullAuthor.nationality || ''
    image.value.value = fullAuthor.image || ''
    dialog.value = true
  } catch (error) {
    showError(`Yazar bilgileri getirme başarısız: ${error.response?.data?.error || error.message}`)
    console.error(error)
  }
}

const deleteAuthor = async (id) => {
  if (!confirm('Bu yazarı silmek istediğinizden emin misiniz?')) return

  try {
    await axios.delete(`/delete-author/${id}`)
    await getAuthors()
  } catch (error) {
    // alert(`Yazar silme başarısız: ${error}`)
    // console.log(error)
  }
}

const submitCreateAuthor = handleSubmit(async () => {
  try {
    const data = {
      name: name.value.value,
      bio: bio.value.value || '',
      birth_year: birth_year.value.value ? parseInt(birth_year.value.value) : null,
      death_year: death_year.value.value ? parseInt(death_year.value.value) : null,
      nationality: nationality.value.value || '',
      image: image.value.value || ''
    }

    await axios.post('/create-author', data)
    await getAuthors()
    dialog.value = false
    await handleReset()
  } catch (error) {
    showError(`Yazar Oluşturma Başarısız: ${error.response?.data?.error || error.message}`)
    console.error(error)
  }
})

const submitEditAuthor = handleSubmit(async () => {
  try {
    const data = {
      name: name.value.value,
      bio: bio.value.value || '',
      birth_year: birth_year.value.value ? parseInt(birth_year.value.value) : null,
      death_year: death_year.value.value ? parseInt(death_year.value.value) : null,
      nationality: nationality.value.value || '',
      image: image.value.value || ''
    }

    await axios.put(`/update-author/${selectedAuthorId.value}`, data)
    await getAuthors()
    dialog.value = false
    await handleReset()
  } catch (error) {
    showError(`Yazar güncelleme başarısız: ${error.response?.data?.error || error.message}`)
    console.error(error)
  }
})

const getAuthors = async () => {
  try {
    const response = await axios.get('/get-all-authors-dropdown')
    authors.value = response.data
    // console.log('Authors from backend:', authors.value)
    // console.log('First author:', authors.value[0])
  } catch (error) {
    showError(`Yazarlar getirme başarısız: ${error.response?.data?.message || error.message}`)
    console.error(error)
  }
}

onMounted(async () => {
  await getAuthors()
})
</script>

<template>

    <v-card>
      <v-btn @click="openCreateDialog" prepend-icon="mdi-plus" class="mb-5 bg-grey-darken-3 text-white">
        Yeni Yazar Ekle
      </v-btn>
      <Table
        :headers="headers"
        :items="authors"
        :search="search"
        :deleteFunc="deleteAuthor"
        :openEditDialog="openEditDialog"
      />

      <!-- Create/Edit Dialog -->
      <v-dialog v-model="dialog" persistent max-width="600px">
        <v-card>
          <v-card-title>
            <span class="text-h5">{{ edit ? 'Yazar Düzenle' : 'Yeni Yazar' }}</span>
          </v-card-title>
          <v-card-text>
            <v-container>
              <v-row>
                <v-col cols="12">
                  <v-text-field
                    v-model="name.value.value"
                    :error-messages="name.errorMessage.value"
                    label="Adı *"
                    required
                  ></v-text-field>
                </v-col>

                <v-col cols="12">
                  <v-textarea
                    v-model="bio.value.value"
                    :error-messages="bio.errorMessage.value"
                    label="Biyografi"
                    rows="4"
                  ></v-textarea>
                </v-col>

                <v-col cols="12" sm="6">
                  <v-text-field
                    v-model="birth_year.value.value"
                    :error-messages="birth_year.errorMessage.value"
                    label="Doğum Yılı"
                    type="number"
                  ></v-text-field>
                </v-col>

                <v-col cols="12" sm="6">
                  <v-text-field
                    v-model="death_year.value.value"
                    :error-messages="death_year.errorMessage.value"
                    label="Ölüm Yılı"
                    type="number"
                  ></v-text-field>
                </v-col>

                <v-col cols="12">
                  <v-text-field
                    v-model="nationality.value.value"
                    :error-messages="nationality.errorMessage.value"
                    label="Uyruk"
                  ></v-text-field>
                </v-col>

                <v-col cols="12">
                  <v-text-field
                    v-model="image.value.value"
                    :error-messages="image.errorMessage.value"
                    label="Resim URL"
                  ></v-text-field>
                </v-col>
              </v-row>
            </v-container>
            <small>* gerekli alanlar</small>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="grey-darken-1" variant="text" @click="closeDialog">
              İptal
            </v-btn>
            <v-btn
              color="primary"
              variant="text"
              @click="edit ? submitEditAuthor() : submitCreateAuthor()"
            >
              {{ edit ? 'Güncelle' : 'Kaydet' }}
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-card>

    <!-- Snackbar for notifications -->
    <SnackbarNotification
      v-model="snackbar.show"
      :message="snackbar.message"
      :color="snackbar.color"
    />
</template>

<style scoped>
</style>
