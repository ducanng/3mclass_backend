package classdm

import (
	"time"

	"gorm.io/gorm"
)

type (
	Class struct {
		ClassID          uint64 `gorm:"primaryKey;autoIncrement" json:"ClassID"`
		ClassCode        string `gorm:"type:varchar(64);not null;unique" json:"ClassCode"`
		ClassName        string `gorm:"type:varchar(64);" json:"ClassName"`
		NumberOfStudents uint64 `gorm:"type:integer" json:"NumberOfStudents"`
		NumberOfTeachers uint64 `gorm:"type:integer" json:"NumberOfTeachers"`
		Description      string `gorm:"type:text" json:"Description"`
		Credits          uint64 `gorm:"type:integer" json:"Credits"`
		OwnerID          uint64 `gorm:"type:bigint;not null" json:"OwnerID"`

		CreatedAt time.Time
		UpdatedAt time.Time
		CreatedBy string
		UpdatedBy string

		DeletedAt gorm.DeletedAt `gorm:"index" json:"DeletedAt"`
	}
	ClassMember struct {
		ClassMemberID uint64 `gorm:"primaryKey;autoIncrement" json:"ClassMemberID"`
		ClassID       uint64 `gorm:"type:bigint;not null" json:"ClassID"`
		UserID        uint64 `gorm:"type:bigint;not null" json:"UserID"`
	}
	ClassInvitation struct {
		ClassInvitationID uint64 `gorm:"primaryKey;autoIncrement" json:"ClassInvitationID"`
		ClassID           uint64 `json:"ClassID"`
		InvitationCode    string `gorm:"type:varchar(64);not null;unique" json:"InvitationCode"`
		Email             string `gorm:"type:varchar(64);not null;" json:"Email"`
	}
)

func (c Class) TableName() string {
	return "class"
}
func (c ClassMember) TableName() string {
	return "classMember"
}
func (c ClassInvitation) TableName() string {
	return "classInvitation"
}
