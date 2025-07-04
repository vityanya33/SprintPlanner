<template>
  <form @submit.prevent="handleSubmit" class="bg-white h-full p-4 rounded-xl w-full shadow-md space-y-4 max-w-md mx-auto mt-10 ml-5">
    <h2 class="text-xl font-semibold text-gray-700">Добавить задачу</h2>

    <div>
      <label class="block text-sm font-medium text-gray-600">Название задачи</label>
      <input
          v-model="form.title"
          type="text"
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
          placeholder="Например: Верстка лендинга"
          required
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-600">Исполнитель</label>
      <select
          v-model="form.userId"
          required
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
      >
        <option disabled value="">Выберите участника</option>
        <option v-for="user in props.users" :key="user.id" :value="user.id">
          {{ user.name }} ({{ user.role || 'Без роли' }})
        </option>
      </select>
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-600">Дата начала</label>
      <input
          v-model="form.startDate"
          type="date"
          required
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-600">Срок (дедлайн)</label>
      <input
          v-model="form.deadline"
          type="date"
          required
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
      />
    </div>

    <button
        type="submit"
        class="w-full bg-green-600 hover:bg-green-700 text-white font-semibold py-2 px-4 rounded-md transition"
    >
      Добавить задачу
    </button>
  </form>
</template>

<script setup>
import { reactive } from 'vue'
import { createTask } from '../api/tasks.js'
// import { usePlannerStore } from '../store/usePlannerStore'
// import { storeToRefs } from 'pinia'
//
// const store = usePlannerStore()
// const { users } = storeToRefs(store)
const emit = defineEmits(['task-added'])

const form = reactive({
  title: '',
  userId: '',
  startDate: '',
  deadline: '',
})



const props = defineProps({
  users: {
    type: Array,
    required: true,
  }
})

const handleSubmit = async () => {
  try {
    const res = await createTask({
      title: form.title,
      userId: form.userId,
      startDate: form.startDate,
      deadline: form.deadline,
    })

    form.title = ''
    form.userId = ''
    form.startDate = ''
    form.deadline = ''
    emit('task-added', res.data)
  } catch (err) {
    console.log(err)
  }
}
</script>
