package models

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	UserID         uint         `gorm:"column:user_id;not null;<-:create" json:"user_id"`
	User           UserModel    `gorm:"foreignKey:UserID;references:ID" json:"-"` // Omitting from JSON response
	Date           time.Time    `gorm:"column:date;type:DATE;default:CURRENT_DATE;not null;<-:create" json:"date"`
	CheckInTime    *time.Time   `gorm:"column:check_in_time;<-:create" json:"check_in_time"` // Pointer to allow nulls
	CheckOutTime   *time.Time   `gorm:"column:check_out_time" json:"check_out_time"`         // Updatable, no restriction
	OrganisationID uint         `gorm:"column:organisation_id;not null;<-:create" json:"organisation_id"`
	Organisation   Organisation `gorm:"foreignKey:OrganisationID;references:ID" json:"-"` // Omitting from JSON response
	CheckInPhoto   int          `gorm:"column:check_in_photo;<-:create" json:"check_in_photo"`
	CheckOutPhoto  int          `gorm:"column:check_out_photo" json:"check_out_photo"` // Updatable, no restriction
}
