<template>
  <div class="max-w-md mx-auto mt-10 ml-10">
    <!-- Кнопка для показа формы добавления вручную -->
    <button
        @click="toggleForm"
        class="flex items-center gap-2 bg-amber-600 hover:bg-amber-700 text-white font-semibold py-2 px-4 rounded-md transition duration-500"
    >
      <span v-if="!isOpen">Add user manually</span>
      <span v-else>Hide the form</span>
      <span class="text-xl font-bold">+</span>
    </button>

    <!-- Форма добавления вручную -->
    <transition name="fade">
      <form
          v-if="isOpen"
          @submit.prevent="handleSubmit"
          class="hover:scale-102 transition-transform duration-300 form p-4 rounded-xl w-full shadow-md space-y-4 mt-4"
      >
        <h2 class="text-xl font-semibold text-gray-700">Add user</h2>
        <div>
          <label class="block text-sm font-bold text-gray-600">Name</label>
          <input v-model.trim="form.name" type="text"
                 class="bg-white mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 text-emerald-700 focus:ring focus:ring-blue-200 p-2"
                 placeholder="Ivan Ivanov" required/>
        </div>
        <div>
          <label class="block text-sm font-bold text-gray-600">Role</label>
          <input v-model.trim="form.role" type="text"
                 class="bg-white mt-1 block w-full rounded-md text-emerald-700 border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200 p-2"
                 placeholder="Developer"/>
        </div>
        <div>
          <label class="block text-sm font-bold text-gray-600">Working resource</label>
          <input v-model.trim="form.resource" type="number" min="1"
                 class="bg-white mt-1 block w-full rounded-md text-emerald-700 border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200 p-2"
                 placeholder="100"/>
        </div>
        <button type="submit"
                class="w-full mt-3 bg-amber-600 hover:bg-amber-700 hover:-translate-y-1.5 duration-700 text-white font-semibold py-2 px-4 rounded-md transition">
          Add
        </button>
      </form>
    </transition>

    <!-- Кнопка для показа формы загрузки XLS -->
    <button
        @click="toggleUploadForm"
        class="flex items-center gap-2 bg-green-600 hover:bg-green-700 text-white font-semibold py-2 px-4 rounded-md transition duration-500 mt-6"
    >
      <span v-if="!isUploadOpen">Upload users with XLS</span>
      <span v-else>Hide upload</span>
      <span class="text-xl font-bold">+</span>
    </button>

    <!-- Форма загрузки XLS -->
    <transition name="fade">
      <form
          v-if="isUploadOpen"
          @submit.prevent="handleUpload"
          class="hover:scale-102 transition-transform duration-300 form p-4 rounded-xl w-full shadow-md space-y-4 mt-4"
      >
        <h2 class="text-xl font-semibold text-gray-700">Upload XLS document</h2>

        <!--Скрываю оригинальный input для более красивого оформления-->
        <div class="relative">
          <input
              id="fileInput"
              type="file"
              name="file"
              @change="handleFileChange"
              accept=".xls,.xlsx"
              class="hidden"
          />
        </div>

        <!--Кнопка для загрузки файла-->
        <label
            for="fileInput"
            class="block cursor-pointer w-full bg-white border border-gray-300 rounded-md py-2 px-3 text-gray-700 hover:bg-gray-200 transition"
        >
          📂 Select XLS/XLSX file
        </label>

        <!--Показать имя файла-->
        <p v-if="selectedFile" class="mt-2 text-sm text-gray-600">
          Selected: {{ selectedFile.name }}
        </p>

        <button type="submit"
                class="w-full mt-3 bg-green-600 hover:bg-green-700 hover:-translate-y-1.5 duration-700 text-white font-semibold py-2 px-4 rounded-md transition">
          Upload
        </button>
      </form>
    </transition>
  </div>
</template>

<script setup>
import {reactive, ref} from 'vue'
import {createUser, uploadUsersXLS} from '../api/users.js'

const emit = defineEmits(['user-added', 'users-uploaded'])

const isOpen = ref(false)
const isUploadOpen = ref(false)
const selectedFile = ref(null)

const form = reactive({
  name: '',
  role: '',
  resource: null
})

const toggleForm = () => {
  isOpen.value = !isOpen.value
}

const toggleUploadForm = () => {
  isUploadOpen.value = !isUploadOpen.value
}

const handleSubmit = async () => {
  try {
    const res = await createUser({
      name: form.name,
      role: form.role,
      resource: form.resource
    })
    form.name = ''
    form.role = ''
    form.resource = null
    emit('user-added', res.data)
    isOpen.value = false
  } catch (err) {
    console.log(err)
  }
}

const handleFileChange = (e) => {
  selectedFile.value = e.target.files[0]
}

const handleUpload = async () => {
  if (!selectedFile.value) return
  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)

    console.log([...formData]) // проверим, что файл там есть

    const res = await uploadUsersXLS(formData)
    emit('users-uploaded', res.data)
    isUploadOpen.value = false
  } catch (err) {
    console.log(err)
  }
}
</script>

<style scoped>
.form {
  background-color: #FFD8BE;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
