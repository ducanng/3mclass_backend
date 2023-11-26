package userdm

import (
	"time"

	"gorm.io/gorm"
)

type (
	User struct {
		UserID       uint64 `gorm:"primaryKey;autoIncrement" json:"UserID"`
		UserEmail    string `gorm:"type:varchar(64);not null;unique" json:"UserEmail"`
		FirstName    string `gorm:"type:varchar(64)" json:"FirstName"`
		LastName     string `gorm:"type:varchar(64)" json:"LastName"`
		Sex          string `gorm:"type:varchar(64)" json:"Sex"`
		PhoneNumber  string `gorm:"type:varchar(64);not null;unique" json:"PhoneNumber"`
		Password     string `gorm:"varchar(100);not null" json:"Password"`
		ActiveStatus bool   `gorm:"type:bool;default:true" json:"ActiveStatus"`
		LoginType    string `gorm:"varchar(64);not null" json:"LoginType"`
		Role         string `gorm:"varchar(64);not null" json:"Role"`
		CreatedBy    string `gorm:"type:varchar(64)" json:"CreatedBy"`
		UpdatedBy    string `gorm:"type:varchar(64)" json:"UpdatedBy"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
		DeletedAt    gorm.DeletedAt `gorm:"index"`
	}
	Student struct {
		StudentID   uint64 `gorm:"primaryKey;autoIncrement" json:"StudentID"`
		StudentCode string `gorm:"type:varchar(64);not null;unique" json:"StudentCode"`
		UserID      uint64 `gorm:"type:bigint;not null" json:"UserID"`
	}
)

func (u User) TableName() string {
	return "user"
}
