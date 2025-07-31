package models

type Task struct {
	ID        int         `json:"id"`           //ID задачи
	Title     string      `json:"title"`        //Описание задачи
	Hours     int         `json:"hours"`        //Оценка задачи в часах
	UserIDs   []int       `json:"userIds"`      //Список ID пользователей
	StartDate string      `json:"startDate"`    //Дата начала
	Deadline  string      `json:"deadline"`     //Дата дедлайн
}
