<template>
  <div class="bg-white p-4 rounded-xl shadow-md space-y-4 ml-15 mr-10 mx-auto mt-10 flex-1 w-[70%] h-93 overflow-y-auto custom-scrollbar task">
    <h2 class="text-xl font-semibold mb-4 text-gray-700">Задачи</h2>

    <div v-if="tasks.length" class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
          v-for="task in props.tasks"
          :key="task.id"
          class="bg-white rounded-xl shadow p-4 border-l-4 border-green-500 relative"
      >
        <!-- Отображение -->
        <div v-if="editingId !== task.id">
          <h3 class="text-lg font-bold text-gray-800">{{ task.title }}</h3>
          <p class="text-sm text-gray-600">
            👤 {{ getUserName(task.userId) }}<br />
            📅 Срок: {{ formatDate(task.deadline) }}
          </p>
          <div class="absolute top-2 right-2 space-x-2">
            <button @click="startEdit(task)" class="text-sm transition-transform duration-300 hover:scale-200">✏️</button>
            <button @click="remove(task.id)" class="text-sm transition-transform duration-300 hover:scale-200">🗑️</button>
          </div>
        </div>

        <!-- Редактирование -->
        <div v-else>
          <input v-model="editTitle" class="mb-2 w-full border rounded px-2 py-1" />

          <select v-model="editUserId" class="mb-2 w-full border rounded px-2 py-1">
            <option disabled value="">— Выберите участника —</option>
            <option v-for="user in props.users" :key="user.id" :value="user.id">
              {{ user.name }}
            </option>
          </select>

          <input v-model="editStartDate" type="date" class="mb-2 w-full border rounded px-2 py-1" />
          <input v-model="editDeadline" type="date" class="mb-2 w-full border rounded px-2 py-1" />

          <div class="flex gap-2">
            <button @click="confirmEdit(task.id)" class="bg-blue-500 text-white px-3 py-1 rounded">Сохранить</button>
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
    required: true,
  },
  users: {
    type: Array,
    required: true,
  }
})
const emit = defineEmits(['task-updated', 'task-removed'])

const editingId = ref(null)
const editTitle = ref('')
const editUserId = ref('')
const editStartDate = ref('')
const editDeadline = ref('')

const getUserName = (id) => {
  const user = props.users.find(u => u.id === id)
  return user ? user.name : '❌ Удалён'
}

const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('ru-RU')
}


//Для изменения задач
const startEdit = (t) => {
  editingId.value = t.id
  editTitle.value = t.title
  editUserId.value = t.userId ?? ''
  editStartDate.value = t.startDate
  editDeadline.value = t.deadline
}

const confirmEdit = async (id) => {
  try {
    await updateTask(id, {
      title: editTitle.value,
      userId: editUserId.value,
      startDate: editStartDate.value,
      deadline: editDeadline.value
    })

    //обновление локального списка
    const task = props.tasks.find(t => t.id === id)
    if (task) {
      emit('task-updated', { id, title: editTitle.value, userId: editUserId.value, startDate: editStartDate.value, deadline: editDeadline.value })}

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
