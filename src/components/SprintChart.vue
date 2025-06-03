<template>
  <div class="mt-10 bg-white rounded-xl shadow p-4 w-[95%] mx-auto overflow-x-auto">
    <h2 class="text-xl font-semibold mb-4 text-gray-700">График задач</h2>

    <div v-if="users.length && tasks.length">
      <div class="grid" :style="`grid-template-columns: 150px repeat(${dayCount}, 1fr);`">
        <!-- Заголовки дней -->
        <div class="font-semibold text-sm">Участник</div>
        <div
            v-for="(day, index) in days"
            :key="index"
            class="text-xs text-center text-gray-500 border-b border-gray-300 pb-1"
        >
          {{ formatDay(day) }}
        </div>

        <!-- Строки по каждому участнику -->
        <template v-for="user in users" :key="user.id">
          <div class="font-medium text-sm text-gray-800 border-t pt-2">{{ user.name }}</div>
          <div
              v-for="(day, index) in days"
              :key="index"
              class="h-6 border-t border-gray-200 relative"
          >
            <!-- Показываем задачу если на эту дату она у юзера -->
            <div
                v-for="task in getTasksForDay(user.id, day)"
                :key="task.id"
                class="absolute top-0 left-0 h-full w-full opacity-70 text-[10px] text-white flex items-center justify-center truncate rounded"
                :style="{
        //HSL метод генерации цвета, для упращения и максимальной уникальности цветов
        backgroundColor: `hsl(${getTaskColor(task.id)}, 70%, 40%)`
      }"
            >
              {{ task.title }}
            </div>
          </div>
        </template>
      </div>
    </div>

    <p v-else class="text-gray-500 italic">Добавьте участников и задачи для отображения графика</p>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { usePlannerStore } from '../store/usePlannerStore'
import { storeToRefs } from 'pinia'

const store = usePlannerStore()
const { users, tasks } = storeToRefs(store)

const startDate = new Date() // по умолчанию — сегодня
const sprintLength = 10 //определяет количество дней на графике

const days = computed(() => {
  const arr = []
  for (let i = 0; i < sprintLength; i++) {
    const d = new Date(startDate)
    d.setDate(d.getDate() + i)
    arr.push(d)
  }
  return arr
})

const dayCount = computed(() => days.value.length)

const formatDay = (date) => {
  return date.toLocaleDateString('ru-RU', { day: '2-digit', month: 'short' })
}

const getTasksForDay = (userId, day) => {
  return tasks.value.filter(task => {
    if (task.userId !== userId) return false

    const start = new Date(task.startDate)
    const end = new Date(task.deadline)
    //Нормализую дату для более правильного создания интервала для задачи
    const normalize = (d) => new Date(d.getFullYear(), d.getMonth(), d.getDate())

    return normalize(day) >= normalize(start) && normalize(day) <= normalize(end)

  })
}

  // Генерация цвета на основе ID пользователя
const getTaskColor = (userId) => {
  //угол распределения цвета диаграмм
  return (parseInt(userId) * 137.508) % 360
}
</script>
