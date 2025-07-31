<template>
  <div class="mt-10 bg-white rounded-xl shadow p-4 w-[95%] mx-auto">
    <h2 class="text-xl font-semibold mb-4 text-gray-700">График задач</h2>
    <div ref="container" class="overflow-x-auto"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { Timeline } from 'vis-timeline/standalone'
import { DataSet } from 'vis-data'
import 'vis-timeline/styles/vis-timeline-graph2d.min.css'

const props = defineProps({
  users: {
    type: Array,
    default: () => []
  },
  tasks: {
    type: Array,
    default: () => []
  }
})

// Валидация — необязательно, но можно оставить для отладки
if (!props.users || !Array.isArray(props.users) || !props.users.length) {
  console.warn('Пустой список пользователей')
}
if (!props.tasks || !Array.isArray(props.tasks) || !props.tasks.length) {
  console.warn('Пустой список задач')
}

const container = ref(null)
let timeline = null

onMounted(() => {
  renderChart()
})

// Обновляем график при изменении задач или пользователей
watch([() => props.tasks, () => props.users], renderChart, { deep: true })

function renderChart() {
  if (!container.value) return

  // Группы — пользователи
  const groups = props.users.map(user => ({
    id: user.id,
    content: user.name,
  }))

  // Элементы — задачи, по каждому пользователю
  const itemsArray = []
  props.tasks.forEach(task => {
    const color = getRandomColor()
    const userIds = Array.isArray(task.userIds) ? task.userIds : []
    userIds.forEach(userId => {
      itemsArray.push({
        id: `${task.id}-${userId}`,
        group: userId,
        content: task.title,
        start: task.startDate,
        end: task.deadline,
        style: `background-color: ${color}; color: white; border: none;`
      })
    })
  })

  const items = new DataSet(itemsArray)

  const options = {
    stack: false,
    horizontalScroll: true,
    zoomKey: 'ctrlKey',
    margin: {
      item: 20,
      axis: 40,
    },
  }

  if (timeline) {
    timeline.setItems(items)
    timeline.setGroups(groups)
  } else {
    timeline = new Timeline(container.value, items, groups, options)
  }
}

// Генерация случайного цвета
function getRandomColor() {
  const clr = Math.floor(Math.random() * 360)
  return `hsl(${clr}, 70%, 50%)`
}
</script>
