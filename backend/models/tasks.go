package models

import "github.com/jackc/pgx/v5/pgtype"

type Task struct {
	ID        int         `json:"id"`         //ID задачи
	Title     string      `json:"title"`      //Описание задачи
	UserID    int         `json:"user_id"`    //ID исполнителя
	StartDate pgtype.Date `json:"start_date"` //Дата начала
	Deadline  pgtype.Date `json:"deadline"`   //Дата дедлайн
}
