<template>
  <div class="bg-white p-4 rounded-xl shadow-md space-y-4 ml-15 mr-10 mx-auto mt-10 flex-1 w-[70%] h-62 overflow-y-auto custom-scrollbar">
    <h2 class="text-xl font-semibold mb-4 text-gray-700">Участники спринта</h2>

    <div v-if="users.length" class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
          v-for="user in props.users"
          :key="user.id"
          class="bg-white rounded-xl shadow p-4 border-l-4 border-blue-500 relative"
      >
        <div v-if="editingId !== user.id">
          <p class="text-xs text-gray-400">ID: {{ user.id }}</p>
          <h3 class="text-lg font-bold text-gray-800">{{ user.name }}</h3>
          <p class="text-sm text-gray-600">{{ user.role || 'Без роли' }}</p>

          <div class="absolute top-2 right-2 space-x-2">
            <button @click="startEdit(user)" class="text-sm transition-transform duration-300 hover:scale-200">✏️</button>
            <button @click="remove(user.id)" class="text-sm transition-transform duration-300 hover:scale-200">🗑️</button>
          </div>
        </div>

        <div v-else>
          <input v-model="editName" class="mb-2 w-full border rounded px-2 py-1" />
          <input v-model="editRole" class="mb-2 w-full border rounded px-2 py-1" />
          <div class="flex gap-2">
            <button @click="confirmEdit(user.id)" class="bg-blue-500 text-white px-3 py-1 rounded">Сохранить</button>
            <button @click="cancelEdit" class="bg-gray-300 px-3 py-1 rounded">Отмена</button>
          </div>
        </div>
      </div>
    </div>

    <p v-else class="text-gray-500 italic">Нет участников</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { updateUser, deleteUser} from '../api/users.js'
//import axios from 'axios'

//import { usePlannerStore } from '../store/usePlannerStore'
//import { storeToRefs } from 'pinia'

//для загрузки с Pinia
//const store = usePlannerStore()
//const { users }  = storeToRefs(store)

const props = defineProps({
  users: {
    type: Array,
    required: true,
  }
})
const emit = defineEmits(['user-updated', "user-removed"])

const editingId = ref(null)
const editName = ref('')
const editRole = ref('')

//для изменения пользователей
const startEdit = (u) => {
  editingId.value = u.id
  editName.value = u.name
  editRole.value = u.role
}

const confirmEdit = async (id) => {
  try {
    await updateUser(id, {
      name: editName.value.trim(),
      role: editRole.value.trim()
    })

    //обновление локального списка
    const user = props.users.find(u => u.id === id)
    if (user) {
      emit('user-updated', { id, name: editName.value.trim(), role: editRole.value.trim() })
    }

    editingId.value = null
  } catch (err) {
  console.log('Some problems with editing users: ', err)
  }
}

//для отмены редактирования
const cancelEdit = () => {
  editingId.value = null
}

//для удаления пользователя
const remove = async (id) => {
  if (!confirm('Удалить пользователя? ')) return
  try {
    await deleteUser(id)
    emit('user-removed', id)
  } catch (err) {
    console.log('Some problems with delete users: ', err)
  }
}
</script>
