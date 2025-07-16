import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL + '/users'

export const getUsers = async () => await axios.get(API_URL)
export const getUser = async (id) => await axios.get(`${API_URL}/${id}`)
export const createUser = async (user) => await axios.post(API_URL, user)
export const updateUser = async (id, u) => {
    await axios.patch(`${API_URL}/${id}`, u)
    console.log(u) }
export const deleteUser = async (id) => await axios.delete(`${API_URL}/${id}`)