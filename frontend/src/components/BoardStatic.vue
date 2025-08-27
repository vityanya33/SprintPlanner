<template>
  <div
      class="mt-5 rounded-xl shadow p-4 w-[94%] mx-auto hover:scale-102 transition-transform duration-600 ml-15 bg-amber-50">
    <h2 class="text-xl font-semibold mb-4 text-gray-700">Drag&Drop</h2>

    <!-- Свободные задачи (горизонтально) -->
    <div class="bg-gray-200 rounded-xl p-4 shadow-md mb-6">
      <h3 class="text-center text-gray-700 font-semibold mb-3">Free tasks</h3>
      <draggable
          :list="lists[0]"
          :group="{ name: 'tasks', pull: true, put: true }"
          item-key="id"
          :data-userid="0"
          @add="onAdd"
          @end="onEnd"
          @change="onChangeLog"
          class="flex gap-3 overflow-x-auto pb-2"
      >
        <template #item="{ element }">
          <div
              class="bg-gray-50 border border-gray-200 rounded-lg p-3 shadow-sm cursor-grab transition hover:bg-blue-50 hover:-translate-y-0.5 hover:shadow-md active:cursor-grabbing min-w-[140px] text-center"
              :data-id="element.id"
          >
            <h1>Task: {{ element.title }}</h1>
            <p>Evaluation: {{ element.hours }}</p>
          </div>
        </template>
      </draggable>
    </div>

    <!-- Колонки пользователей (в ряд со скроллом) -->
    <div class="flex gap-5 items-start overflow-x-auto pb-2 scrollbar-custom-DaD">
      <div
          v-for="user in users"
          :key="user.id"
          class="bg-pink-50 rounded-xl p-4 w-60 min-h-[320px] shadow-md flex-shrink-0"
      >
        <h3 class="font-semibold text-gray-700 mb-3 text-center">{{ user.name }}</h3>
        <div class="flex justify-between mb-3">
          <p class="text-green-600 ml-2">Free: {{ user.free }}</p>
          <p class="text-red-800 mr-2">Busy: {{ user.busy }}</p>
        </div>
        <draggable
            :list="lists[user.id]"
            :group="{ name: 'tasks', pull: true, put: true }"
            item-key="id"
            :data-userid="user.id"
            @add="onAdd"
            @end="onEnd"
            @change="onChangeLog"
            class="space-y-2"
        >
          <template #item="{ element }">
            <div
                class="bg-gray-50 border border-gray-200 rounded-lg p-3 shadow-sm cursor-grab transition hover:bg-blue-50 hover:-translate-y-0.5 hover:shadow-md active:cursor-grabbing min-w-[140px] text-center"
                :data-id="element.id"
            >
              <h1>Task: {{ element.title }}</h1>
              <p>Evaluation: {{ element.hours }}</p>
            </div>
          </template>
        </draggable>
      </div>
    </div>
  </div>
</template>


<script setup>
import {ref, watch, onMounted} from 'vue'
import draggable from 'vuedraggable'
import {setTaskUsers} from '../api/tasks'

// ===== props / emits =====
const props = defineProps({
  users: {type: Array, default: () => []},
  // task: { id, title, startDate, deadline, userIds: number[] }
  tasks: {type: Array, default: () => []}
})
const emit = defineEmits(['tasks-updated'])

// ===== локальное состояние списков =====
// ключ = userId (0 — Unassigned), значение = массив объектов задач (ссылки на props.tasks[*])
const lists = ref({})

// Создаём списки из входных props
function buildListsFromProps(reason = 'init') {
  console.groupCollapsed(`[BoardStatic] buildListsFromProps: ${reason}`)
  console.log('users:', props.users)
  console.log('tasks:', props.tasks)

  const next = {}
  next[0] = []
  for (const u of props.users) next[u.id] = []

  for (const t of props.tasks) {
    const primary = Array.isArray(t.userIds) && t.userIds.length ? t.userIds[0] : 0
    if (!next[primary]) next[primary] = []
    next[primary].push(t)
  }

  // печать содержимого
  Object.keys(next).forEach(uid => {
    console.log(`→ list[${uid}] length = ${next[uid].length}`)
  })

  lists.value = next
  console.groupEnd()
}

// Перестраиваем, когда реально меняются данные сверху
watch(
    () => [props.users, props.tasks],
    () => buildListsFromProps('props changed'),
    {deep: true, immediate: true}
)

onMounted(() => {
  console.info('[BoardStatic] mounted')
  buildListsFromProps('mounted')
})

// ===== D&D handlers =====

// Только логируем «change», чтобы видеть added/removed/moved
function onChangeLog(evt) {
  console.groupCollapsed('[BoardStatic] @change')
  try {
    if (evt?.added) {
      console.log('event.added:', evt.added)
    }
    if (evt?.removed) {
      console.log('event.removed:', evt.removed)
    }
    if (evt?.moved) {
      console.log('event.moved (inside one column):', evt.moved)
    }
  } finally {
    console.groupEnd()
  }
}

// Основная логика — перенос между колонками
async function onAdd(evt) {
  console.groupCollapsed('[BoardStatic] @add')
  try {
    // DOM-узлы контейнеров
    const toEl = evt?.to
    const fromEl = evt?.from

    const newUserId = Number(toEl?.dataset?.userid ?? 0)
    const oldUserId = Number(fromEl?.dataset?.userid ?? 0)

    // В какую позицию в целевом списке вставили
    const newIndex = evt?.newIndex ?? -1

    console.table({
      oldUserId,
      newUserId,
      oldIndex: evt?.oldIndex,
      newIndex
    })

    // Берём задачу из целевого массива (она уже вставлена `vuedraggable`-ом)
    const targetList = lists.value[newUserId] || []
    const task = targetList[newIndex]
    console.log('targetList length:', targetList.length)
    console.log('task at newIndex:', task)

    if (!task) {
      console.warn('❗ Didn\'t find the task in the target list. Interrupting.')
      return
    }

    const currentPrimary = Array.isArray(task.userIds) && task.userIds.length ? task.userIds[0] : 0
    console.log('currentPrimary:', currentPrimary, ' -> newUserId:', newUserId)

    // Если «перенесли» в ту же колонку — ничего не делаем (только порядок поменяли)
    if (currentPrimary === newUserId) {
      console.log('Transferring within the same column - skips DB update.')
      return
    }

    // 1) локально обновляем userIds (график и доска увидят изменение сразу)
    task.userIds = newUserId === 0 ? [] : [newUserId]
    console.log('Updated task.userIds locally:', task.userIds)

    // 2) сохраняем в БД
    try {
      await setTaskUsers(task.id, task.userIds)
      console.log(`✅ DB sync ok: task #${task.id} -> user_ids = [${task.userIds.join(',')}]`)
    } catch (err) {
      console.error('❌ DB sync failed, I visually roll back to the previous list', err)

      // Откат (возвращаем задачу в старую колонку)
      // Удаляем из нового места
      targetList.splice(newIndex, 1)
      // Вставляем в старый список на конец (или oldIndex, если хочется точно)
      if (!lists.value[oldUserId]) lists.value[oldUserId] = []
      lists.value[oldUserId].push(task)

      // И восстанавливаем userIds локально
      task.userIds = oldUserId === 0 ? [] : [oldUserId]
      return
    }

    // уведомим родителя (если он захочет перезагрузиться)
    emit('tasks-updated', [...props.tasks])
  } finally {
    console.groupEnd()
  }
}

// Просто лог в конце dnd
function onEnd(evt) {
  console.groupCollapsed('[BoardStatic] @end')
  console.table({
    fromUserId: evt?.from?.dataset?.userid,
    toUserId: evt?.to?.dataset?.userid,
    oldIndex: evt?.oldIndex,
    newIndex: evt?.newIndex
  })
  console.groupEnd()
}
</script>

<style>
.scrollbar-custom-DaD {
  scrollbar-width: thin;
  scrollbar-color: #bcb2cb #ece4e4;
}

</style>
