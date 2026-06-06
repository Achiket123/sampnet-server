package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, a)
}

type IntArray []int

func (a IntArray) Value() (driver.Value, error) {
	if a == nil {
		return json.Marshal([]int{})
	}
	return json.Marshal(a)
}

func (a *IntArray) Scan(value interface{}) error {
	if value == nil {
		*a = []int{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, a)
}

type JSONMap map[string]interface{}

func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return json.Marshal(map[string]interface{}{})
	}
	return json.Marshal(m)
}

func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = map[string]interface{}{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, m)
}

type PeopleContact struct {
	gorm.Model
	OrganisationID uint        `gorm:"index;not null"`
	FirstName      string      `gorm:"size:100;not null"`
	LastName       string      `gorm:"size:100;not null"`
	Email          *string     `gorm:"size:255;index"`
	Phone          *string     `gorm:"size:50"`
	Company        *string     `gorm:"size:255"`
	JobTitle       *string     `gorm:"size:255"`
	AvatarUrl      *string     `gorm:"type:text"`
	Type           string      `gorm:"size:20;not null"`
	Status         string      `gorm:"size:20;default:'active'"`
	Stage          string      `gorm:"size:100;default:'new'"`
	Tags           StringArray `gorm:"type:jsonb;default:'[]'"`
	ListIDs        IntArray    `gorm:"type:jsonb;default:'[]'"`
	AssignedToID   *uint       `gorm:"index"`
	AssignedTo     *Employee   `gorm:"foreignKey:AssignedToID"`
	Source         *string     `gorm:"size:50"`
	DealValue      *float64    `gorm:"type:numeric(15,2)"`
	Currency       string      `gorm:"size:10;default:'USD'"`
	Notes          *string     `gorm:"type:text"`
	Address        *string     `gorm:"type:text"`
	City           *string     `gorm:"size:100"`
	Country        *string     `gorm:"size:100"`
	LinkedinUrl    *string     `gorm:"type:text"`
	WebsiteUrl     *string     `gorm:"type:text"`
	CustomFields   JSONMap     `gorm:"type:jsonb;default:'{}'"`
	LastContactedAt *time.Time
}

type PeopleInteraction struct {
	gorm.Model
	ContactID    uint           `gorm:"index;not null"`
	Contact      PeopleContact  `gorm:"foreignKey:ContactID;constraint:OnDelete:CASCADE;"`
	CreatedByID  uint           `gorm:"index;not null"`
	CreatedBy    Employee       `gorm:"foreignKey:CreatedByID"`
	Type         string         `gorm:"size:30;not null"`
	Content      string         `gorm:"type:text;not null"`
	Outcome      *string        `gorm:"size:20"`
	LinkedTaskID *uint          `gorm:"index"`
	OccurredAt   time.Time      `gorm:"default:current_timestamp"`
}

type PeoplePipelineStage struct {
	gorm.Model
	OrganisationID uint   `gorm:"index;not null"`
	Key            string `gorm:"size:100;not null"`
	Label          string `gorm:"size:100;not null"`
	Color          string `gorm:"size:20;not null;default:'#6B7280'"`
	Order          int    `gorm:"not null"`
	IsDefault      bool   `gorm:"default:false"`
	IsWon          bool   `gorm:"default:false"`
	IsLost         bool   `gorm:"default:false"`
}

type PeopleList struct {
	gorm.Model
	OrganisationID uint   `gorm:"index;not null"`
	Name           string `gorm:"size:200;not null"`
	Description    string `gorm:"type:text;default:''"`
	Color          string `gorm:"size:20;default:'#6B7280'"`
}

type PeopleListContact struct {
	ListID    uint       `gorm:"primaryKey;autoIncrement:false"`
	ContactID uint       `gorm:"primaryKey;autoIncrement:false"`
	AddedAt   time.Time  `gorm:"default:current_timestamp"`
	List      PeopleList `gorm:"foreignKey:ListID;constraint:OnDelete:CASCADE;"`
	Contact   PeopleContact `gorm:"foreignKey:ContactID;constraint:OnDelete:CASCADE;"`
}
