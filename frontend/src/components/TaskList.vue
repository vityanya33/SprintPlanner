<template>
  <div class="form p-4 rounded-xl hover:scale-102 transition-transform duration-600 shadow-md space-y-4 ml-15 mr-10 mx-auto mt-10 flex-1 w-[70%] h-93 overflow-y-auto custom-scrollbar task">
    <h2 class="text-xl font-semibold mb-4 text-gray-700">Задачи</h2>

    <div v-if="tasks && tasks.length" class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
          v-for="task in props.tasks"
          :key="task.id"
          class="bg-white rounded-xl shadow p-4 border-l-4 border-green-500 relative"
      >
        <!-- Отображение -->
        <div v-if="editingId !== task.id">
          <h3 class="text-lg font-bold text-gray-800">{{ task.title }}</h3>
          <p class="text-sm text-gray-600">
            👤 Участники: {{ getUserNames(task.userIds) }}<br />
            ⏱️ Количество часов на выполнение: {{ task.hours }}<br />
            📅 Срок: {{ formatDate(task.startDate) }} - {{ formatDate(task.deadline) }}
          </p>
          <div class="absolute top-2 right-2 space-x-2">
            <button @click="startEdit(task)" class="text-sm transition-transform duration-300 hover:scale-200">✏️</button>
            <button @click="remove(task.id)" class="text-sm transition-transform duration-300 hover:scale-200">🗑️</button>
          </div>
        </div>

        <!-- Редактирование -->
        <div v-else>
          <input v-model="editTitle" class="mb-2 w-full border rounded px-2 py-1" />

          <input v-model="editHours" type="number" class="mb-2 w-full border rounded px-2 py-1" />

          <select multiple v-model="editUserIds" class="mb-2 w-full border rounded px-2 py-1">
            <option v-for="user in props.users" :key="user.id" :value="user.id">
              {{ user.name }}
            </option>
          </select>

          <input v-model="editStartDate" type="date" class="mb-2 w-full border rounded px-2 py-1" />
          <input v-model="editDeadline" type="date" class="mb-2 w-full border rounded px-2 py-1" />

          <div class="flex gap-2">
            <button @click="confirmEdit" class="bg-blue-500 text-white px-3 py-1 rounded">Сохранить</button>
            <button @click="cancelEdit" class="bg-gray-300 px-3 py-1 rounded">Отмена</button>
          </div>
        </div>
      </div>
    </div>

    <p v-else class="text-gray-500 italic">Нет задач</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { updateTask, deleteTask } from '../api/tasks.js'

// import { usePlannerStore } from '../store/usePlannerStore'
// import { storeToRefs } from 'pinia'
//
// const store = usePlannerStore()
// const { tasks, users } = storeToRefs(store)

const props = defineProps({
  tasks: {
    type: Array,
    default: () => []
  },
  users: {
    type: Array,
    default: () => []
  }
})
const emit = defineEmits(['task-updated', 'task-removed', 'update-users'])

const editingId = ref(null)
const editTitle = ref('')
const editHours = ref(null)
const editUserIds = ref([])
const editStartDate = ref('')
const editDeadline = ref('')

const getUserNames = (ids) => {
  if (!Array.isArray(ids)) return '—'
  return ids.map(id => {
        const user = props.users.find(u => u.id === id)
        return user ? user.name : '❌Удален'
      })
      .join(', ')
}

const formatDate = (dateStr) => {
  if (!dateStr || typeof dateStr !== 'string') return '—'

  const parts = dateStr.split('-')
  if (parts.length !== 3) return '—'

  const [year, month, day] = parts.map(Number)
  const date = new Date(year, month - 1, day)

  if (isNaN(date.getTime())) return '—'

  return date.toLocaleDateString('ru-RU', { day: '2-digit', month: 'long', year: 'numeric' })
}


//Для изменения задач
const startEdit = (task) => {
  editingId.value = task.id
  editTitle.value = task.title
  editHours.value = task.hours
  editUserIds.value = (task.userIds ?? []).map(Number)
  editStartDate.value = task.startDate
  editDeadline.value = task.deadline
}

const confirmEdit = async () => {
  const task = props.tasks.find(t => t.id === editingId.value)
  if (!task) {
    console.error('Не найдена задача с таким ID:', editingId.value)
    return
  }

  const updatedTask = {
    id: editingId.value,
    title: editTitle.value,
    hours: editHours.value,
    userIds: editUserIds.value,
    startDate: editStartDate.value,
    deadline: editDeadline.value
  }

  try {
    await updateTask(editingId.value, updatedTask)

    emit('task-updated', updatedTask)
    emit('update-users')

    editingId.value = null
  } catch (err) {
    console.log('Some problems with editing tasks: ', err)
  }
}

//для отмены редактирования
const cancelEdit = () => {
  editingId.value = null
}

//для удаления задачи
const remove = async(id) => {
  if (!confirm('Удалить задачу?')) return
  try {
    await deleteTask(id)
    emit('task-removed', id)
  } catch (err) {
    console.log('Some problems with delete tasks: ', err)
  }
}
</script>
<style scoped>
.form {
  background-color: #D6F6DD;
}
</style>
