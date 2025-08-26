<template>
  <div
      class="mt-5 bg-white rounded-xl shadow p-4 w-[94%] mx-auto hover:scale-102 transition-transform duration-600 ml-15"
  >
    <h2 class="text-xl font-semibold mb-4 text-gray-700">График задач</h2>
    <div ref="container" class="overflow-x-auto"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import { Timeline } from "vis-timeline/standalone";
import { DataSet } from "vis-data";
import "vis-timeline/styles/vis-timeline-graph2d.min.css";

const props = defineProps({
  users: {
    type: Array,
    default: () => [],
  },
  tasks: {
    type: Array,
    default: () => [],
  },
});

const container = ref(null);
let timeline = null;

onMounted(() => {
  renderChart();
});

watch([() => props.tasks, () => props.users], renderChart, { deep: true });

function renderChart() {
  if (!container.value) return;

  // Подсчет задач у пользователя
  const taskCount = {};
  props.tasks.forEach((task) => {
    const userIds = Array.isArray(task.userIds) ? task.userIds : [];
    userIds.forEach((userId) => {
      taskCount[userId] = (taskCount[userId] || 0) + 1;
    });
  });

  // Группы (пользователи)
  const groups = props.users.map((user) => ({
    id: user.id,
    name: user.name,
    count: taskCount[user.id] || 0,
  }));

  // Элементы (таски)
  const itemsArray = [];
  props.tasks.forEach((task) => {
    const color = getRandomColor();
    const userIds = Array.isArray(task.userIds) ? task.userIds : [];
    userIds.forEach((userId) => {
      itemsArray.push({
        id: `${task.id}-${userId}`,
        group: userId,
        content: task.title,
        start: task.startDate,
        end: task.deadline,
        style: `background-color: ${color}; color: black; border: groove;`,
      });
    });
  });

  const items = new DataSet(itemsArray);

  const options = {
    stack: false,
    horizontalScroll: true,
    zoomKey: "ctrlKey",
    margin: {
      item: 20,
      axis: 40,
    },
    groupTemplate: function (group) {
      const el = document.createElement("span");
      el.textContent = group.name;
      if (group.count > 1) {
        el.style.color = "darkred";
        el.style.textDecoration = "underline";
        el.style.fontWeight = "bold";
      }
      return el;
    },
  };

  if (timeline) {
    timeline.setItems(items);
    timeline.setGroups(groups);
  } else {
    timeline = new Timeline(container.value, items, groups, options);
  }
}

// Шаблон цветов
const colors = [
  "#FFCCCC", "#FFCC99", "#FFFF99", "#CCFFCC", "#CCFFFF", "#CCCCFF", "#FFCCFF", "#FFE5CC",
  "#FFE5B4", "#FFFFCC", "#E5FFCC", "#CCFFE5", "#E5CCFF", "#FFCCE5", "#FFD9B3", "#FFFFE0",
  "#E0FFCC", "#CCFFE0", "#E0CCFF", "#FFCCE0", "#FFE6B3", "#FFFFCC", "#E6FFCC", "#CCFFF5",
  "#E6CCFF", "#FFCCF2", "#FFD9CC", "#FFFFD9", "#D9FFCC", "#CCFFF0", "#D9CCFF", "#FFCCF5",
  "#FFE6CC", "#FFFFE6", "#E6FFD9", "#CCFFF2", "#E6CCF2", "#FFCCF0", "#FFDACC", "#FFFFF0",
  "#DAFFCC", "#CCFFF7", "#DACCF2", "#FFCCFA", "#FFEACC", "#FFFFDA", "#EAFFCC", "#CCFFFA",
  "#EACCFF", "#FFCCFF", "#FFDCCC", "#FFFFF2", "#DCFFCC", "#CCFFFD", "#DCCCFF", "#FFCCFB",
  "#FFECCC", "#FFFFDC", "#ECFFCC", "#CCFFFF", "#ECCCFF", "#FFCCFC", "#FFDECC", "#FFFFEC",
  "#DEFFCC", "#CCF5FF", "#DECCFF", "#FFCCFD", "#FFECCC", "#FFFFF5", "#F0FFCC", "#CCF0FF",
  "#F0CCFF", "#FFCCFE", "#FFF2CC", "#FFFFFA", "#F2FFCC", "#CCE5FF", "#F2CCFF", "#FFCCFF"
];

// Копия шаблона для доступных цветов
let availableColors = [...colors];

// Функция выбора случайного цвета
function getRandomColor() {
  const randomIdx = Math.floor(Math.random() * availableColors.length);
  const color = availableColors[randomIdx];

  // Убираем выбранный цвет из доступных
  availableColors.splice(randomIdx, 1);

  return color;
}
</script>
