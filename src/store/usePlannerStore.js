import { defineStore } from 'pinia'
import { ref } from 'vue'

export const usePlannerStore = defineStore('planner', () => {
    const users = ref([])
    const tasks = ref([])

    function addUser(user) {
        users.value.push({
            id: Date.now() % 1000,
            name: user.name,
            role: user.role || '',
        })
    }

    function addTask(task) {
        tasks.value.push({
            id: Date.now() % 1000,
            title: task.title,
            userId: task.userId,
            startDate: task.startDate,
            deadline: task.deadline,
        })
    }

    function updateUser(id, data) {
        const user = users.value.find(u => u.id === id)
        if (user) {
            user.name = data.name
            user.role = data.role
        }
    }

    function updateTask(id, data) {
        const task = tasks.value.find(t => t.id === id)
        if (task) {
            task.title = data.title
            task.userId = data.userId
            task.startDate = data.startDate
            task.deadline = data.deadline
        }
    }

    function removeUser(id) {
        console.log('Удаляем пользователя с id:', id)
        console.log('До удаления задач:', tasks.value.map(t => t.userId))
        users.value = users.value.filter(u => u.id !== id)
        tasks.value = tasks.value.filter(t => t.userId !== id)
        console.log('После удаления задач:', tasks.value.map(t => t.userId))
    }

    function removeTask(id) {
        tasks.value = tasks.value.filter(t => t.id !== id)
    }

    return {
        users,
        tasks,
        addUser,
        addTask,
        updateUser,
        updateTask,
        removeUser,
        removeTask
    }
})