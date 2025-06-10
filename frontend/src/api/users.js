import axios from 'axios'

const API_URL = 'http://localhost:3000/users'

export const getUsers = async () => await axios.get(API_URL)
export const getUser = async (id) => await axios.get(`${API_URL}/${id}`)
export const createUser = async (u) => await axios.post(API_URL, u)
export const updateUser = async (id, user) => {
    await axios.patch(`${API_URL}/${id}`, user)
    console.log(u) }
export const deleteUser = async (id) => await axios.delete(`${API_URL}/${id}`)