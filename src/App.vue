<template>
  <div class="min-h-screen bg-gray-100 pt-5">
    <div class="flex w-screen justify-between pb-10">
      <UserForm />
      <UserList />
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
</script>

