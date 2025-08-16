<template>
  <div class="min-h-screen phone pt-5">
    <RouterView :tasks="tasks" :users="users" />
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { RouterView } from 'vue-router'
import { getUsers } from './api/users.js'
import { getTasks } from './api/tasks.js'

//Для загрузки задач и юзеров
const tasks = ref([])
const users = ref([])

onMounted(async () => {
  try {
    const resTasks = await getTasks()
    tasks.value = resTasks.data
  } catch (err) {
    console.log('Some problems with download tasks: ', err)
  }

  try {
    const resUsers = await getUsers()
    users.value = resUsers.data
  } catch (err) {
    console.log('Some problems with download users: ',err)
  }
})
</script>

<style scoped>
.phone {
  background-color: #E2DADB;
}
</style>
