<template>
  <div class="mt-5 rounded-xl shadow p-4 w-[94%] mx-auto hover:scale-102 transition-transform duration-600 ml-15 bg-amber-50">
    <h2 class="text-xl font-semibold mb-4 text-gray-700">Drag&Drop</h2>
    <div class="flex gap-5 items-start">
      <!-- Свободные задачи -->
      <div class="bg-gray-200 rounded-xl p-4 w-60 min-h-[320px] shadow-md">
        <h3 class="text-center text-gray-700 font-semibold mb-3">Свободные задачи</h3>
        <draggable
            :list="lists[0]"
            :group="{ name: 'tasks', pull: true, put: true }"
            item-key="id"
            :data-userid="0"
            @add="onAdd"
            @end="onEnd"
            @change="onChangeLog"
            class="space-y-2"
        >
          <template #item="{ element }">
            <div
                class="bg-gray-50 border border-gray-200 rounded-lg p-3 shadow-sm cursor-grab transition hover:bg-blue-50 hover:-translate-y-0.5 hover:shadow-md active:cursor-grabbing"
                :data-id="element.id"
            >
              {{ element.title }}
            </div>
          </template>
        </draggable>
      </div>

      <!-- Колонки пользователей -->
      <div
          v-for="user in users"
          :key="user.id"
          class="bg-pink-50 rounded-xl p-4 w-60 min-h-[320px] shadow-md"
      >
        <h3 class="text-center text-gray-700 font-semibold mb-3">{{ user.name }}</h3>
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
                class="bg-gray-50 border border-gray-200 rounded-lg p-3 shadow-sm cursor-grab transition hover:bg-blue-50 hover:-translate-y-0.5 hover:shadow-md active:cursor-grabbing"
                :data-id="element.id"
            >
              {{ element.title }}
            </div>
          </template>
        </draggable>
      </div>
    </div>
  </div>
</template>


<script setup>
import { ref, watch, onMounted } from 'vue'
import draggable from 'vuedraggable'
import { setTaskUsers } from '../api/tasks'

// ===== props / emits =====
const props = defineProps({
  users: { type: Array, default: () => [] },
  // task: { id, title, startDate, deadline, userIds: number[] }
  tasks: { type: Array, default: () => [] }
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
    { deep: true, immediate: true }
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
      console.log('event.moved (внутри одной колонки):', evt.moved)
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
      console.warn('❗ Не нашли задачу в целевом списке — странно. Прерываю.')
      return
    }

    const currentPrimary = Array.isArray(task.userIds) && task.userIds.length ? task.userIds[0] : 0
    console.log('currentPrimary:', currentPrimary, ' -> newUserId:', newUserId)

    // Если «перенесли» в ту же колонку — ничего не делаем (только порядок поменяли)
    if (currentPrimary === newUserId) {
      console.log('Перенос внутри той же колонки — пропускаю обновление БД.')
      return
    }

    // 1) локально обновляем userIds (график и доска увидят изменение сразу)
    task.userIds = newUserId === 0 ? [] : [newUserId]
    console.log('Обновил task.userIds локально:', task.userIds)

    // 2) сохраняем в БД
    try {
      await setTaskUsers(task.id, task.userIds)
      console.log(`✅ DB sync ok: task #${task.id} -> user_ids = [${task.userIds.join(',')}]`)
    } catch (err) {
      console.error('❌ DB sync failed, откатываю визуально в прежний список', err)

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

    // 3) уведомим родителя (если он захочет перезагрузиться)
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


<style scoped>
.board-form {
  display: flex;
  gap: 20px;
  padding: 20px;
  margin-left: 15px;
  background: linear-gradient(135deg, #f9fafb, #eef2f7);
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
  align-items: flex-start;
  transition: all 0.3s ease;
}

.column {
  background: #ffffff;
  border-radius: 12px;
  padding: 14px;
  width: 240px;
  min-height: 320px;
  box-shadow: 0 2px 6px rgba(0,0,0,0.08);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.column:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 10px rgba(0,0,0,0.12);
}

.column-title {
  font-size: 16px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 12px;
  text-align: center;
}

.task-card {
  background: #fdfdfd;
  padding: 10px;
  margin-bottom: 10px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  box-shadow: 0 1px 3px rgba(0,0,0,0.07);
  cursor: grab;
  transition: all 0.2s ease;
}

.task-card:hover {
  transform: translateY(-2px);
  background: #f0f9ff;
  border-color: #bae6fd;
  box-shadow: 0 3px 6px rgba(0,0,0,0.1);
}

.task-card:active { cursor: grabbing; }
</style>
