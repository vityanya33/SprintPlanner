package models

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID      string `json:"id"`      //ID задачи
	Title   string `json:"title"`   //Описание задачи
	Hours   int    `json:"hours"`   //Оценка задачи в часах
	UserIDs []int  `json:"userIds"` //Список ID пользователей

	StartDate time.Time `json:"-"` //Дата начала (исключаем из JSON)
	Deadline  time.Time `json:"-"` //Дата дедлайн (исключаем из JSON)
}

// MarshalJSON кастомная сериализация в JSON
func (t Task) MarshalJSON() ([]byte, error) {
	type Alias Task
	return json.Marshal(&struct {
		Alias
		StartDate string `json:"startDate"`
		Deadline  string `json:"deadline"`
	}{
		Alias:     Alias(t),
		StartDate: t.StartDate.Format("2006-01-02"),
		Deadline:  t.Deadline.Format("2006-01-02"),
	})
}

// UnmarshalJSON кастомная десериализация из JSON
func (t *Task) UnmarshalJSON(data []byte) error {
	type Alias Task
	aux := &struct {
		StartDate string `json:"startDate"`
		Deadline  string `json:"deadline"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Парсим строки в time.Time
	layout := "2006-01-02"
	if aux.StartDate != "" {
		if startDate, err := time.Parse(layout, aux.StartDate); err == nil {
			t.StartDate = startDate
		}
	}
	if aux.Deadline != "" {
		if deadline, err := time.Parse(layout, aux.Deadline); err == nil {
			t.Deadline = deadline
		}
	}

	return nil
}
