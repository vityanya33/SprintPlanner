<template>
  <div class="bg-gray-50 border rounded p-3 mb-3 shadow-sm">
    <div class="flex justify-between items-start gap-2">
      <div>
        <h5 class="font-medium">{{ task.title }}</h5>
        <p class="text-sm text-gray-600">hours: {{ task.hours }}</p>
        <p v-if="task.userIds?.length" class="text-xs text-gray-400">assigned: {{ task.userIds.join(', ') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  task: { type: Object, required: true }
})
</script>

<style scoped>
/* доп. стили по вкусу */
</style>
