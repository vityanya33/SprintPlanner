package models

type Task struct {
	ID        int    `json:"id"`         //ID задачи
	Title     string `json:"title"`      //Описание задачи
	UserID    int    `json:"user_id"`    //ID исполнителя
	StartDate string `json:"start_date"` //Дата начала
	Deadline  string `json:"deadline"`   //Дата дедлайн
}
