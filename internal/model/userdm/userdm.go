package userdm

import (
	"time"

	"gorm.io/gorm"
)

type (
	User struct {
		UserID       uint64 `gorm:"primaryKey;autoIncrement" json:"user_id"`
		UserEmail    string `gorm:"type:varchar(64);not null;unique" json:"user_email"`
		DisplayName  string `gorm:"type:varchar(64)" json:"display_name"`
		PhoneNumber  string `gorm:"type:varchar(64);not null;unique" json:"phone_number"`
		Password     string `gorm:"varchar(100);not null" json:"password"`
		ActiveStatus bool   `gorm:"type:bool;default:true" json:"active_status"`

		CreatedBy string `gorm:"type:varchar(64)" json:"created_by"`
		UpdatedBy string `gorm:"type:varchar(64)" json:"updated_by"`
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}

	UserCredential struct {
		CredentialID uint64 `gorm:"primaryKey;autoIncrement" json:"credential_id"`
		UserID       uint64 `gorm:"bigint;not null" json:"user_id"`
		Password     string `gorm:"varchar(100);not null" json:"password"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}
)

func (u User) TableName() string {
	return "tb_user"
}
