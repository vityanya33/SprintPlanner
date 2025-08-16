import { createRouter, createWebHistory } from 'vue-router'
import Home from '../pages/Home.vue'
import Tasks from '../pages/Tasks.vue'
import Users from '../pages/Users.vue'

const routes = [
    { path: '/', name: 'Home', component: Home },
    { path: '/tasks', name: 'Tasks', component: Tasks },
    { path: '/users', name: 'Users', component: Users },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

export default router