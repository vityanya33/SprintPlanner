import axios from 'axios'

const API_URL = 'http://localhost:3000/tasks'

export const getTasks = async () => await axios.get(API_URL)
export const getTask = async (id) => await axios.get(`${API_URL}/${id}`)
export const createTask = async (task) => await axios.post(API_URL, task)
export const updateTask = async (id, t) => {
    await axios.patch(`${API_URL}/${id}`, t)
    console.log(t) }
export const deleteTask = async (id) => await axios.delete(`${API_URL}/${id}`)