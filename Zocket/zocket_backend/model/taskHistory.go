package models

import "gorm.io/gorm"

type TaskHistory struct {
	gorm.Model
	TaskID    uint
	UpdatedBy uint
	OldStatus string
	NewStatus string
}
