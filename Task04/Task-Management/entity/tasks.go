package entity

import "time"

type Task struct {
	TaskId      int64     `json:"taskid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
	CreatedBy   string    `json:"createdby"`
	AssignedTo  []string  `json:"assignedto"`
}
