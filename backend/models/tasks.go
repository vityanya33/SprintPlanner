package models

type Task struct {
	ID        int         `json:"id"`           //ID задачи
	Title     string      `json:"title"`        //Описание задачи
	UserID    int         `json:"userId"`       //ID исполнителя
	StartDate string      `json:"startDate"`    //Дата начала
	Deadline  string      `json:"deadline"`     //Дата дедлайн
}
