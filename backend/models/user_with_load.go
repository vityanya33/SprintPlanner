package models
//Структура для отображения пользователей с их занятостью
type UserWithLoad struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Resource int    `json:"resource"`
	Busy     int    `json:"busy"`
	Free     int    `json:"free"`
}