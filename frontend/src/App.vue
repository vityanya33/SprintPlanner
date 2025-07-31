<template>
  <div class="min-h-screen phone pt-5">
    <div class="flex w-screen justify-between pb-10 shadow-2xl">
      <UserForm @user-added="handleUserAdded" />
      <UserList :users="users"
                :tasks="tasks"
                @user-updated="handleUserUpdated"
                @user-removed="handleUserRemoved" />
    </div>
    <div class="flex w-screen justify-between pb-10 shadow-2xl">
      <TaskForm
          :users="users"
          @task-added="handleTaskAdded"
          @update-users="updateUsers"/>
      <TaskList
          :tasks="tasks"
          :users="users"
          @task-updated="handleTaskUpdated"
          @task-removed="handleTaskRemoved"
          @update-users="updateUsers"/>
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
import {getTasks, getTask, updateTask} from './api/tasks.js'

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
const handleUserUpdated = async () => {
  alert('Обновление пользователя!')
  const res = await getUsers()
  users.value = res.data
  console.log('Пользователи обновлены:', users.value)
}


//убрать удаленного
const handleUserRemoved = async () => {
  const res = await getUsers()
  users.value = res.data
}

//Задачи

//добавить задачу
const handleTaskAdded = (task) => {
  tasks.value.push(task)
}

//обновить существующую задачу
const handleTaskUpdated = (updatedTask) => {
  const index = tasks.value.findIndex(t => t.id === updatedTask.id)
  if (index !== -1) {
    tasks.value.splice(index, 1, { ...updatedTask })  // вот так нужно
  }
}
//убрать удаленную задачу
const handleTaskRemoved = async () => {
  try {
    const res = await getTasks()
    tasks.value = res.data
  } catch (err) {
    console.error('Ошибка при удалении задачи:', err)
  }
}

const updateUsers = async () => {
  const res = await getUsers()
  users.value = res.data
}
</script>
<style scoped>
.phone {
  background-color: #E2DADB;
}
</style>
