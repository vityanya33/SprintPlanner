<template>
  <!-- Заголовок -->
  <h1 class="text-4xl font-extrabold uppercase tracking-widest text-lime-600 drop-shadow ml-25">
    Задачи
  </h1>

  <!-- Бургер-меню -->
  <AppMenu />

  <div class="flex w-screen justify-between pb-10">
    <TaskForm
        :users="users"
        @task-added="handleTaskAdded"
    />
    <TaskList
        :tasks="tasks"
        :users="users"
        @task-updated="handleTaskUpdated"
        @task-removed="handleTaskRemoved"
        @update-tasks="updateTasks"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import TaskForm from '../components/TaskForm.vue'
import TaskList from '../components/TaskList.vue'
import AppMenu from '../components/AppMenu.vue'

import { getTasks } from '../api/tasks.js'
import { getUsers } from '../api/users.js'

const tasks = ref([])
const users = ref([])

// загрузка при открытии страницы
onMounted(async () => {
  await updateTasks()
})

// добавление задачи
const handleTaskAdded = async () => {
  await updateTasks()
}

// обновление существующей
const handleTaskUpdated = async () => {
  await updateTasks()
}

// удаление задачи
const handleTaskRemoved = async () => {
  await updateTasks()
}

// общий метод для обновления списка задач и пользователей
const updateTasks = async () => {
  try {
    const resUsers = await getUsers()
    users.value = resUsers.data
    console.log('Пользователи обновлены:', users.value)
  } catch (err) {
    console.error('Ошибка при загрузке пользователей:', err)
  }

  try {
    const resTasks = await getTasks()
    tasks.value = resTasks.data
    console.log('Задачи обновлены:', tasks.value)
  } catch (err) {
    console.error('Ошибка при загрузке задач:', err)
  }
}
</script>
