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
      <label class="block text-sm font-medium text-gray-600">Исполнители</label>
      <Multiselect
          v-model="form.userIds"
          :options="props.users"
          :multiple="true"
          :close-on-select="false"
          :clear-on-select="false"
          :preserve-search="true"
          placeholder="Выберите участника"
          label="name"
          track-by="id"
          class="multiselect-control"
      />
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
  userIds: [],
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
      userIds: form.userIds.map(u => u.id),
      startDate: form.startDate,
      deadline: form.deadline,
    })

    form.title = ''
    form.userIds = []
    form.startDate = ''
    form.deadline = ''
    emit('task-added', res.data)
  } catch (err) {
    console.log(err)
  }
}
</script>
