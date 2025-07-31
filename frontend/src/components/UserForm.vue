<template>
  <form @submit.prevent="handleSubmit" class="hover:scale-102 transition-transform duration-300 form h-full p-4 rounded-xl w-full shadow-md space-y-4 max-w-md mx-auto mt-10 ml-10">
    <h2 class="text-xl font-semibold text-gray-700">Добавить участника</h2>
    <div>
      <label class="block text-sm font-bold text-gray-600">Имя</label>
      <input
          v-model.trim="form.name"
          type="text"
          class="bg-white mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 text-emerald-700 focus:ring focus:ring-blue-200 p-2"
          placeholder="Иван Иванов"
          required
      />
    </div>
    <div>
      <label class="block text-sm font-bold text-gray-600">Роль</label>
      <input
          v-model.trim="form.role"
          type="text"
          class="bg-white mt-1 block w-full rounded-md text-emerald-700 border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200 p-2"
          placeholder="Разработчик, дизайнер и т.д."
      />
    </div>
    <div>
      <label class="block text-sm font-bold text-gray-600">Ресурс в часах</label>
      <input
          v-model.trim="form.resource"
          type="number"
          min="1"
          class="bg-white mt-1 block w-full rounded-md text-emerald-700 border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200 p-2"
          placeholder="100"
      />
    </div>
    <button
        type="submit"
        class="w-full mt-3 bg-amber-600 hover:bg-amber-700 hover:-translate-y-1.5 duration-700 text-white font-semibold py-2 px-4 rounded-md transition"
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
  role: '',
  resource: null
})

const handleSubmit = async () => {
  try {
    const res = await createUser({
      name: form.name,
      role: form.role,
      resource: form.resource
    })

    form.name = ''
    form.role = ''
    form.resource = null
    emit('user-added', res.data)

  } catch (err) {
    console.log(err)
  }
}
</script>
<style scoped>
.form {
  background-color: #FFD8BE;
}
</style>
