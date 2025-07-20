<template>
  <div class="min-h-screen bg-gray-100 pt-5">
    <div class="flex w-screen justify-between pb-10">
      <UserForm @user-added="handleUserAdded" />
      <UserList :users="users"
                @user-updated="handleUserUpdated"
                @user-removed="handleUserRemoved" />
    </div>
    <div class="flex w-screen justify-between pb-10">
      <TaskForm
          :users="users"
          @task-added="handleTaskAdded"/>
      <TaskList
          :tasks="tasks"
          :users="users"
          @task-updated="handleTaskUpdated"
          @task-removed="handleTaskRemoved"
      />
    </div>
    <div>
      <SprintChart :users="users || []" :tasks="tasks || []" />
    </div>
  </div>
</template>
<script setup>

import UserForm from './components/UserForm.vue'
import UserList from './components/UserList.vue'
import TaskForm from './components/TaskForm.vue'
import TaskList from './components/TaskList.vue'
import SprintChart from './components/SprintChart.vue'
import { ref, onMounted } from 'vue'
import { getUsers } from './api/users.js'
import { getTasks } from './api/tasks.js'

//для загрузки людей с бэка
const users = ref([])

//для загрузки задач с бэка
const tasks = ref([])

onMounted(async () => {
  try {
    const res = await getUsers()
    users.value = res.data
  } catch (err) {
    console.log('Some problems with download users: ',err)
  }
  try {
    const res = await getTasks()
    tasks.value = res.data
  } catch (err) {
    console.log('Some problems with download tasks: ',err)

  }
})

//Пользователи

//добавить пользователя
const handleUserAdded = (newUser) => {
  users.value.push(newUser)
}

//обновить существующих
const handleUserUpdated = (updatedUser) => {
  const idx = users.value.findIndex(u => u.id === updatedUser.id)
  if (idx !== -1) users.value[idx] = updatedUser
}

//убрать удаленного
const handleUserRemoved = (removedId) => {
  users.value = users.value.filter(u => u.id !== removedId)
}

//Задачи

//добавить задачу
const handleTaskAdded = (task) => {
  tasks.value.push(task)
}

//обновить существующую задачу
const handleTaskUpdated = (updatedTask) => {
  const idx = tasks.value.findIndex(t => t.id === updatedTask.id)
  if (idx !== -1) tasks.value[idx] = updatedTask
}

//убрать удаленную задачу
const handleTaskRemoved = (removedId) => {
  tasks.value = tasks.value.filter(t => t.id !== removedId)
}
</script>

