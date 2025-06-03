package models

type User struct {
	ID   int    `json:"id"`   //Уникальный ID
	Name string `json:"name"` //Имя
	Role string `json:"role"` //Роль
}
