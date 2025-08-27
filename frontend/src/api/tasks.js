import axios from 'axios'

// Использую переменную окружения, которую читает Vite
const API_URL = import.meta.env.VITE_API_URL + '/tasks'

//CRUD методы
export const getTasks = async () => await axios.get(API_URL)
export const getTask = async (id) => await axios.get(`${API_URL}/${id}`)
export const createTask = async (task) => await axios.post(API_URL, task)
export const updateTask = async (id, data) => {
    return await axios.patch(`${API_URL}/${id}`, data)
}
export const deleteTask = async (id) => await axios.delete(`${API_URL}/${id}`)

//функции для загрузки доступных пользователей
export const getAvailableUsers = async (startDate, deadline, hours) => {
    try {
        return await axios.get(`${API_URL}/available`, {
            params: {
                hours: hours
            }
        })
    } catch (err) {
        console.error('Error getting available users:', err)
        throw err
    }
}

export const setTaskUsers = async (taskId, userIds) => {
    console.log('[API] setTaskUsers', { taskId, userIds })
    return axios.patch(`${API_URL}/${taskId}/users`, { user_ids: userIds })
}
