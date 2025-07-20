package models

type Task struct {
	ID        int         `json:"id"`           //ID задачи
	Title     string      `json:"title"`        //Описание задачи
	UserIDs   []int       `json:"userIds"`      //Список ID пользователей
	StartDate string      `json:"startDate"`    //Дата начала
	Deadline  string      `json:"deadline"`     //Дата дедлайн
}
