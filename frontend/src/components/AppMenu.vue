<template>
  <div>
    <!-- Кнопка-бургер -->
    <button
        @click="openMenu"
        class="z-50 fixed top-5 left-4 w-10 h-10 flex flex-col justify-center items-center gap-1 group bg-lime-200 rounded-xl shadow-md"
    >
      <span class="w-6 h-1 bg-black rounded transition-transform duration-300 group-hover:scale-x-125"></span>
      <span class="w-6 h-1 bg-black rounded transition-transform duration-300 group-hover:scale-x-110"></span>
      <span class="w-6 h-1 bg-black rounded transition-transform duration-300 group-hover:scale-x-125"></span>
    </button>

    <!-- Затемнённый фон -->
    <transition name="fade">
      <div
          v-show="isOpen"
          @click.self="closeMenu"
          class="fixed inset-0 bg-black/30 z-40 backdrop-blur-sm"
      ></div>
    </transition>

    <!-- Сайдбар -->
    <transition name="slide">
      <aside
          v-if="isOpen"
          class="fixed top-0 left-0 w-72 h-full bg-gradient-to-br from-white to-gray-100 shadow-2xl z-50 p-6 flex flex-col gap-6 rounded-tr-2xl rounded-br-2xl"
      >
        <h2 class="text-2xl font-bold mb-4 text-gray-800">📋 Menu</h2>
        <router-link
            to="/"
            @click="closeMenu"
            class="flex items-center gap-2 px-4 py-2 rounded-lg transition-all hover:bg-lime-100 hover:text-lime-800 font-medium text-gray-700"
        >
          🏠 Main
        </router-link>
        <router-link
            to="/users"
            @click="closeMenu"
            class="flex items-center gap-2 px-4 py-2 rounded-lg transition-all hover:bg-lime-100 hover:text-lime-800 font-medium text-gray-700"
        >
          👥 Users
        </router-link>
        <router-link
            to="/tasks"
            @click="closeMenu"
            class="flex items-center gap-2 px-4 py-2 rounded-lg transition-all hover:bg-lime-100 hover:text-lime-800 font-medium text-gray-700"
        >
          📜 Tasks
        </router-link>
        <div class="mt-auto mb-0 text-sm text-gray-400"> © SprintPlanner 2025 created by Vitek</div>
      </aside>
    </transition>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const isOpen = ref(false)

const openMenu = () => {
  isOpen.value = true
}

const closeMenu = () => {
  isOpen.value = false
}
</script>

<style scoped>
.slide-enter-active {
  transition: transform 0.5s ease;
}
.slide-leave-active {
  transition: transform 0.5s ease;
}
.slide-enter-from {
  transform: translateX(-100%);
}
.slide-leave-to {
  transform: translateX(-100%);
}
</style>
<style>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.5s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
