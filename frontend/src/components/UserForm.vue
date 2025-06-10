<template>
  <form @submit.prevent="handleSubmit" class="bg-white h-full p-4 rounded-xl w-full shadow-md space-y-4 max-w-md mx-auto mt-10 ml-5">
    <h2 class="text-xl font-semibold text-gray-700">Добавить участника</h2>
    <div>
      <label class="block text-sm font-medium text-gray-600">Имя</label>
      <input
          v-model.trim="form.name"
          type="text"
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
          placeholder="Иван Иванов"
          required
      />
    </div>
    <div>
      <label class="block text-sm font-medium text-gray-600">Роль</label>
      <input
          v-model.trim="form.role"
          type="text"
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
          placeholder="Разработчик, дизайнер и т.д."
      />
    </div>
    <button
        type="submit"
        class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-md transition"
    >
      Добавить
    </button>
  </form>
</template>

<script setup>
import { reactive } from 'vue'
import { createUser } from '../api/users.js'
//import { usePlannerStore } from '../store/usePlannerStore'
//const store = usePlannerStore()
const emit = defineEmits(['user-added'])

const form = reactive({
  name: '',
  role: ''
})

const handleSubmit = async () => {
  try {
    const res = await createUser({
      name: form.name,
      role: form.role,
    })

    form.name = ''
    form.role = ''
    emit('user-added', res.data)

  } catch (err) {
    console.log(err)
  }
}
</script>
