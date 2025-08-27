<template>
  <!-- Заголовок -->
  <h1 class="text-4xl font-extrabold uppercase tracking-widest text-lime-600 drop-shadow ml-25">
    Users
  </h1>

  <!-- Бургер-меню -->
  <AppMenu />

  <div class="flex w-screen justify-between pb-10">
    <UserForm
        @user-added="handleUserAdded"
        @users-uploaded="handleUsersUploaded"/>
    <UserList
        :users="users"
        @user-updated="handleUserUpdated"
        @user-removed="handleUserRemoved"
        @update-users="updateUsers"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import UserForm from '../components/UserForm.vue'
import UserList from '../components/UserList.vue'
import AppMenu from '../components/AppMenu.vue'
import { getUsers } from '../api/users.js'

const users = ref([])

// загрузка при открытии страницы
onMounted(async () => {
  await updateUsers()
})

// добавление пользователя
const handleUserAdded = async () => {
  await updateUsers()
}

// обновление существующего
const handleUserUpdated = async () => {
  await updateUsers()
}

// удаление пользователя
const handleUserRemoved = async () => {
  await updateUsers()
}

// загрузка XLS и XLSX документов
const handleUsersUploaded = async () => {
  await updateUsers()
}

// общий метод для обновления списка
const updateUsers = async () => {
  try {
    const res = await getUsers()
    users.value = res.data
    console.log('Users updated:', users.value)
  } catch (err) {
    console.error('Error with loading users:', err)
  }
}
</script>
