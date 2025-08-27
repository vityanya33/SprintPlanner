<template>
  <div class="form p-4 hover:scale-102 transition-transform duration-600 rounded-xl shadow-md space-y-4 ml-15 mr-10 mx-auto mt-10 flex-1 w-[70%] h-[615px] overflow-y-auto scrollbar-custom">
    <h2 class="text-xl font-semibold mb-4 text-gray-700">Sprint participants</h2>

    <div v-if="users.length" class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
          v-for="user in props.users"
          :key="user.id"
          class="bg-white rounded-xl shadow p-4 border-l-4 border-amber-600 relative"
      >
        <div v-if="editingId !== user.id">
          <p class="text-xs text-gray-400">ID: {{ user.id }}</p>
          <h3 class="text-lg font-bold text-gray-800">{{ user.name }}</h3>
          <p class="text-sm text-gray-600">
            🎓 Role: {{ user.role }}<br />
            ⏱️ Working resource: {{ user.resource }}<br />
            ⏱️ Busy: {{ user.busy }}<br />
            ⏱️ Free: {{ user.free }}
          </p>

          <div class="absolute top-2 right-2 space-x-2">
            <button @click="startEdit(user)" class="text-sm transition-transform duration-300 hover:scale-200">✏️</button>
            <button @click="remove(user.id)" class="text-sm transition-transform duration-300 hover:scale-200">🗑️</button>
          </div>
        </div>

        <div v-else>
          <input v-model="editName" class="mb-2 w-full border rounded px-2 py-1" />
          <input v-model="editRole" class="mb-2 w-full border rounded px-2 py-1" />
          <input v-model="editResource" type="number" class="mb-2 w-full border rounded px-2 py-1" />
          <div class="flex gap-2">
            <button @click="confirmEdit(user.id)" class="bg-amber-600 hover:bg-amber-700 text-white px-3 py-1 rounded duration-600">Save</button>
            <button @click="cancelEdit" class="bg-red-400 hover:bg-red-500 px-3 py-1 rounded duration-600">Cancel</button>
          </div>
        </div>
      </div>
    </div>

    <p v-else class="text-gray-500 italic">No participants</p>
  </div>
</template>

<script setup>
import { ref} from 'vue'
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
const editResource = ref(null)

//для изменения пользователей
const startEdit = (u) => {
  editingId.value = u.id
  editName.value = u.name
  editRole.value = u.role
  editResource.value = u.resource
}

const confirmEdit = async (id) => {
  try {
    await updateUser(id, {
      name: editName.value.trim(),
      role: editRole.value.trim(),
      resource: editResource.value,
    })

    emit('user-updated')

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
  if (!confirm('Delete user? ')) return
  try {
    await deleteUser(id)
    emit('user-removed', id)
  } catch (err) {
    console.log('Some problems with delete users: ', err)
  }
}
</script>

<style scoped>
.form {
  background-color: #FEEFE5;
}
</style>

<style>
.scrollbar-custom {
  scrollbar-width: thin;
  scrollbar-color: #e58a62 #ffffff;
}

.scrollbar-custom::-webkit-scrollbar {
  width: 10px;
}

.scrollbar-custom::-webkit-scrollbar-thumb {
  background-color: #e58a62;
  border-radius: 8px;
}

.scrollbar-custom::-webkit-scrollbar-track {
  background-color: #ffffff;
}
</style>