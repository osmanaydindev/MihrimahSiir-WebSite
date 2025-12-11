<script setup lang="ts">
import {onMounted, ref} from "vue";
import {useAppStore} from "../../store/app";
import axios from "axios";
import BookCard from "../../components/frontend/book/BookCard.vue"
const store = useAppStore()
const loading = ref(false)

const getReadsBooks = async () => {
  try {
    const res = await axios.get(`/get-reads-books-ids/${store.getUserId}`)
    store.setReadBookIds(res.data)
  } catch (error) {
    console.error('Failed to fetch read book IDs:', error)
  }
}
const addBookToReads = async (a) => {
  await store.addBooksRead(a)
  try {
    await axios.post(`/add-book-to-reads/${store.getUserId}`, a)
      .then((res) => {
      }).catch((error) => {
      })

  }catch (e) {
  }
}
const removeBookFromReads = async (a) => {
  await store.removeBooksRead(a)
  try {
    await axios.post(`/delete-book-from-reads/${store.getUserId}`, a)
      .then((res) => {
      }).catch((error) => {
      })

  }catch (e) {
  }
}
const getAllBooks = async () => {
  if (store.allBooks.length <= 0){
    try {
      const res = await axios.get(`/get-books`)
      store.setAllBooks(res.data)
    } catch (error) {
    }
  }
}

onMounted(async ()=>{
  loading.value = true
  try {
    await getAllBooks()
    await getReadsBooks()
  } finally {
    loading.value = false
  }
})

</script>

<template>

  <div class="books-div page
    d-flex flex-column
    flex-sm-row flex-sm-wrap justify-sm-center
    ga-4 px-4 mt-8 pb-16"
  >
    <div class="px-4 pt-12 text-center">
      <h1>Kitaplar</h1>
    </div>

    <div
      v-if="!loading"
      v-for="(book,index) in store.allBooks"
      :key="index"
      class="w-100 w-sm-50 w-md-33 w-lg-25 w-xl-20 w-xxl-16 mx-0 flex justify-center"
    >
      <BookCard
        :book="book"
        @addBookToReads="addBookToReads"
        @removeBookFromReads="removeBookFromReads"
      />
    </div>

    <v-card v-else class="bg-gray-darken-2">
      <v-card-text>
        Kitaplar Yükleniyor...
      </v-card-text>
      <v-card-actions class="flex flex-row justify-center align-center">
        <v-progress-circular
          v-if="loading"
          indeterminate
          color="white"
          size="24"
        ></v-progress-circular>
      </v-card-actions>
    </v-card>
  </div>
</template>

<style scoped>

</style>
