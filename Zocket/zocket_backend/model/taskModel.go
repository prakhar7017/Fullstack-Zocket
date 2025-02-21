package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	Pending    TaskStatus = "PENDING"
	InProgress TaskStatus = "IN_PROGRESS"
	Completed  TaskStatus = "COMPLETED"
)

type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	AssigneeID  uint       `json:"assignee_id"`
	Status      TaskStatus `json:"status" gorm:"type:varchar(20)"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Deadline    time.Time  `json:"deadline"`
	Importance  int        `json:"importance" gorm:"type:int"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	t.Status = Pending
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

func (t *TaskStatus) FromString(s string) error {
	switch s {
	case "pending":
		*t = Pending
	case "in_progress":
		*t = InProgress
	case "completed":
		*t = Completed
	default:
		return fmt.Errorf("invalid task status: %s", s)
	}
	return nil
}
