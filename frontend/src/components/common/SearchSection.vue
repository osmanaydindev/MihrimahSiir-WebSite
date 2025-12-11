<script setup lang="ts">
interface Props {
  modelValue: string
  label?: string
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  label: 'Ara',
  placeholder: 'Ara...'
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'search': []
  'clear': []
}>()

const handleInput = (value: string) => {
  emit('update:modelValue', value)
}

const handleClear = () => {
  emit('update:modelValue', '')
  emit('clear')
}

const handleSearch = () => {
  emit('search')
}
</script>

<template>
  <div class="search-section mb-8">
    <v-text-field
      class="search-field"
      :model-value="modelValue"
      @update:model-value="handleInput"
      :label="label"
      :placeholder="placeholder"
      outlined
      clearable
      rounded
      bg-color="grey-darken-2"
      variant="solo-inverted"
      @keyup.enter="handleSearch"
      @click:clear="handleClear"
    >
      <template v-slot:prepend-inner>
        <v-icon icon="mdi-magnify" size="20" />
      </template>
    </v-text-field>
  </div>
</template>

<style scoped>
.search-section {
  width: 100%;
  display: flex;
  justify-content: center;
}

.search-field {
  max-width: 500px;
  width: 100%;
}
</style>
