<template>
  <div class="max-w-md mx-auto mt-10 ml-10 mr-20">
    <!-- Кнопка для показа формы добавления вручную -->
    <button
        @click="toggleForm"
        class="flex items-center gap-2 bg-lime-600 hover:bg-lime-700 text-white font-semibold py-2 px-4 rounded-md transition duration-500"
    >
      <span v-if="!isOpen">Add a task manually</span>
      <span v-else>Hide the form</span>
      <span class="text-xl font-bold">+</span>
    </button>

    <!-- Кнопка для показа формы добавления из Jira -->
    <button
        @click="toggleJiraForm"
        class="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-md transition duration-500 mt-4"
    >
      <span v-if="!isJiraOpen">Add tasks from Jira</span>
      <span v-else>Hide Jira form</span>
      <span class="text-xl font-bold">+</span>
    </button>

    <!-- Форма добавления вручную -->
    <transition name="fade">
      <form
          v-if="isOpen"
          @submit.prevent="handleSubmit"
          class="hover:scale-102 transition-transform duration-300 form p-4 rounded-xl w-full shadow-md space-y-4 mt-4"
      >
        <h2 class="text-xl font-semibold text-gray-700">Add a task</h2>
        <div>
          <label class="block text-sm font-bold text-gray-600">Task name</label>
          <input v-model="form.title" type="text"
                 class="bg-white p-2 mt-1 block w-full text- rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
                 placeholder="Layout of the main page" required/>
        </div>
        <label>
          <input type="checkbox" v-model="showOnlyAvailable"/>
          Only available users
        </label>
        <div>
          <label class="block text-sm font-bold text-gray-600">Performers</label>
          <Multiselect
              v-model="selectedUsers"
              :options="showOnlyAvailable ? availableUsers : allUsers"
              :multiple="true"
              :close-on-select="false"
              :clear-on-select="false"
              :preserve-search="true"
              :custom-label="u => `${u.name} (${u.free ?? '-'} free)`"
              placeholder="Select performer"
              label="name"
              track-by="id"
              class="multiselect-control"
          />
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-600">Task assessment (resource)</label>
          <input
              v-model="form.hours"
              type="number"
              min="1"
              class="bg-white p-2 mt-1 block w-full text- rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
              placeholder="10"
              required
          />
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-600">Start date</label>
          <input
              v-model="form.startDate"
              type="date"
              required
              class="bg-white p-2 mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
          />
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-600">Deadline</label>
          <input
              v-model="form.deadline"
              type="date"
              required
              class="bg-white p-2 mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
          />
        </div>

        <button
            type="submit"
            class="w-full bg-green-600 hover:bg-green-700 hover:-translate-y-1.5 duration-700 text-white font-semibold py-2 px-4 rounded-md transition mt-3"
        >
          Add
        </button>
      </form>
    </transition>

    <!-- Форма добавления из Jira -->
    <transition name="fade">
      <form
          v-if="isJiraOpen"
          @submit.prevent="handleJiraSubmit"
          class="hover:scale-102 transition-transform duration-300 bg-blue-100 p-4 rounded-xl w-full shadow-md space-y-4 mt-4"
      >
        <h2 class="text-xl font-semibold text-gray-700">Import tasks from Jira</h2>
        <div>
          <label class="block text-sm font-bold text-gray-600">Jira Query</label>
          <textarea
              v-model="jiraJQL"
              class="bg-white p-2 mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
              placeholder="project = ABC AND status = 'In Progress'"
              required
          ></textarea>
        </div>
        <button
            type="submit"
            class="w-full bg-blue-600 hover:bg-blue-700 hover:-translate-y-1.5 duration-700 text-white font-semibold py-2 px-4 rounded-md transition mt-3"
            :disabled="isJiraLoading"
        >
          <span v-if="!isJiraLoading">Import</span>
          <span v-else>Loading...</span>
        </button>
      </form>
    </transition>
  </div>
</template>

<script setup>
import {reactive, ref, watch, onMounted} from 'vue'
import {createTask, getAvailableUsers, syncTasksWithJira, getTasks} from '../api/tasks.js'
import {getUsers} from '../api/users.js'

const allUsers = ref([])
const availableUsers = ref([])
const showOnlyAvailable = ref(false)
const selectedUsers = ref([])

const emit = defineEmits(['task-added', 'update-users'])

const isOpen = ref(false)
const isJiraOpen = ref(false)

const form = reactive({
  title: '',
  userIds: [],
  startDate: '',
  deadline: '',
  hours: null,
})

const jiraJQL = ref('')
const isJiraLoading = ref(false)

const props = defineProps({
  users: {
    type: Array,
    required: true,
  }
})

const toggleForm = () => {
  isOpen.value = !isOpen.value
}

const toggleJiraForm = () => {
  isJiraOpen.value = !isJiraOpen.value
}

onMounted(async () => {
  const res = await getUsers()
  allUsers.value = res.data
})

watch([showOnlyAvailable, () => form.startDate, () => form.deadline, () => form.hours], async ([show, start, end, h]) => {
  const hInt = parseInt(h)
  if (show && start && end && !isNaN(hInt) && hInt > 0) {
    try {
      const res = await getAvailableUsers(start, end, hInt)
      availableUsers.value = res.data
    } catch (err) {
      console.log('Error loading available users', err)
    }
  }
})

watch(selectedUsers, (newUsers) => {
  form.userIds = newUsers.map(u => u.id)
})

const handleSubmit = async () => {
  try {
    const res = await createTask({
      title: form.title,
      hours: form.hours,
      userIds: form.userIds,
      startDate: form.startDate,
      deadline: form.deadline,
    })
    form.title = ''
    form.hours = null
    form.userIds = []
    form.startDate = ''
    form.deadline = ''
    emit('task-added', res.data)
    emit('update-users')
  } catch (err) {
    console.log(err)
  }
}

const handleJiraSubmit = async () => {
  isJiraLoading.value = true
  try {
    await syncTasksWithJira(jiraJQL.value)
    // после синхронизации обновляем списки
    emit('task-added')
    emit('update-users')
    jiraJQL.value = ''
  } catch (err) {
    console.log('Error importing Jira tasks', err)
  } finally {
    isJiraLoading.value = false
  }
}
</script>

<style scoped>
.form {
  background-color: #B2EDC5;
}
</style>
