package assignmentdm

import (
	"time"

	"gorm.io/gorm"
)

type Assignment struct {
	AssignmentID   uint    `gorm:"primaryKey;autoIncrement" json:"AssignmentID"`
	AssignmentName string  `gorm:"type:varchar(64);not null" json:"AssignmentName"`
	Description    string  `gorm:"type:varchar(255)" json:"Description"`
	GradePoint     float64 `gorm:"type:integer;not null" json:"GradePoint"`
	ClassID        uint    `gorm:"type:bigint;not null" json:"ClassID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string

	DeletedAt gorm.DeletedAt `gorm:"index" json:"DeletedAt"`
}

func (a Assignment) TableName() string {
	return "assignment"
}
