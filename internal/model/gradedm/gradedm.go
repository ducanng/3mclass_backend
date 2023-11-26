package gradedm

import (
	"time"

	"gorm.io/gorm"
)

type GradeStructure struct {
	GradeStructureID uint64 `gorm:"primaryKey;autoIncrement" json:"GradeStructureID"`
	Type             string `gorm:"type:varchar(64);not null" json:"Type"`
	ClassID          uint64 `gorm:"not null" json:"ClassID"`
	GradeName        string `gorm:"type:varchar(64);not null" json:"GradeName"`
	GradeScale       string `gorm:"type:varchar(64);not null" json:"GradeScale"`
	Position         string `gorm:"type:varchar(64);not null" json:"Position"`

	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string

	DeletedAt gorm.DeletedAt `gorm:"index" json:"DeletedAt"`
}

type Grade struct {
	GradeID          uint64  `gorm:"primaryKey;autoIncrement" json:"GradeID"`
	AssignmentID     uint64  `gorm:"type:bigint;not null" json:"AssignmentID"`
	GradeStructureID uint64  `gorm:"type:bigint;not null" json:"GradeStructureID"`
	StudentID        uint64  `gorm:"type:bigint;not null" json:"StudentID"`
	Mark             float64 `gorm:"type:integer;not null" json:"Mark"`
	Remark           string  `gorm:"type:varchar(255)" json:"Remark"`

	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string

	DeletedAt gorm.DeletedAt `gorm:"index" json:"DeletedAt"`
}

func (g Grade) TableName() string {
	return "grade"
}
func (g GradeStructure) TableName() string {
	return "gradeStructure"
}
