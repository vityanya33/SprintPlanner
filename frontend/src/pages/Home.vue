<template>
  <!-- Заголовок -->
  <h1 class="text-4xl font-extrabold uppercase tracking-widest text-lime-600 drop-shadow ml-25">
    Главная
  </h1>

  <!-- Бургер-меню -->
  <AppMenu />

  <!-- График -->
  <div>
    <SprintChart :users="users || []" :tasks="tasks || []" />
  </div>

  <!-- D&D доска -->
  <BoardStatic :users="users" :tasks="tasks" @tasks-updated="tasks = $event" />

  <!-- Адрес для страницы -->
  <router-view />
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppMenu from '../components/AppMenu.vue'
import SprintChart from '../components/SprintChart.vue'
import BoardStatic from '../components/BoardStatic.vue'
import TaskList from '../components/TaskList.vue'
import UserList from '../components/UserList.vue'
import { getUsers } from '../api/users.js'
import { getTasks } from '../api/tasks.js'

const users = ref([])
const tasks = ref([])

const updateUsers = async () => {
  try {
    const res = await getUsers()
    users.value = res.data
  } catch (error) {
    console.log('Error in loading users:',error)
  }
}

const updateTasks = async () => {
  try {
    const res = await getTasks()
    tasks.value = res.data
  } catch (error) {
    console.log('Error in loading tasks:',error)
  }
}

onMounted(async () => {
  await updateUsers()
  await updateTasks()
})
</script>
