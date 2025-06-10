<template>
  <div class="min-h-screen bg-gray-100 pt-5">
    <div class="flex w-screen justify-between pb-10">
      <UserForm @user-added="handleUserAdded" />
      <UserList :users="users"
                @user-updated="handleUserUpdated"
                @user-removed="handleUserRemoved" />
    </div>
    <div class="flex w-screen justify-between pb-10">
      <TaskForm />
      <TaskList />
    </div>
    <div>
      <SprintChart />
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

//для загрузки людей с бэка
const users = ref([])

onMounted(async () => {
  try {
    const res = await getUsers()
    users.value = res.data
  } catch (err) {
    console.log('Some problems with download users: ',err)
  }
})

//добавить пользователя
const handleUserAdded = (newUser) => {
  users.value.push(newUser)
}

//обновить существующих
const handleUserUpdated = (updatedUser) => {
  const idx = users.value.findIndex(user.id === updatedUser.id)
  if (idx !== -1) users.value[idx] = updatedUser
}

//убрать удаленного
const handleUserRemoved = (removedId) => {
  users.value = users.value.filter(u => u.id !== removedId)
}
</script>

