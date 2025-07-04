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
  users: Array,
  tasks: Array,
})

const container = ref(null)
let timeline = null

onMounted(() => {
  renderChart()
})

watch(() => props.tasks, renderChart, { deep: true })

function renderChart() {
  if (!container.value) return

  // Группы это пользователям
  const groups = props.users.map(user => ({
    id: user.id,
    content: user.name,
  }))

  // Элементы это задачи
  const items = new DataSet(
      props.tasks.map(task => {
        const color = getRandomColor()
        return {
          id: task.id,
          group: task.userId,
          content: task.title,
          start: task.startDate,
          end: task.deadline,
          style: `background-color: ${color}; color: white; border: none;`
        }
      })
  )

  const options = {
    stack: false,
    horizontalScroll: true,
    zoomKey: 'ctrlKey',
    min: new Date(),
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

//Функция создания рандомных цветов для колонок
function getRandomColor() {
  const clr = Math.floor(Math.random() * 360)
  return `hsl(${clr}, 70%, 50%)`
}

</script>
