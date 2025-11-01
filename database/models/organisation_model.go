package models

import (
	"time"

	"gorm.io/gorm"
)

// Company represents an organization in the system.
type Organisation struct {
	gorm.Model                    // Unique identifier
	CompanyName        string     `gorm:"type:varchar(255);not null" json:"company_name"`         // Organization name
	CompanyCode        string     `gorm:"type:varchar(50);unique;not null" json:"company_code"`   // Unique company code (e.g., for subdomain)
	PrimaryContactName string     `gorm:"type:varchar(100)" json:"primary_contact_name"`          // Primary point of contact
	PrimaryEmail       string     `gorm:"type:varchar(255);unique;not null" json:"primary_email"` // Contact email
	PhoneNumber        string     `gorm:"type:varchar(20)" json:"phone_number"`                   // Contact phone number
	OfficeAddress      string     `gorm:"type:text" json:"office_address"`                        // Address information
	City               string     `gorm:"type:varchar(100)" json:"city"`                          // City
	State              string     `gorm:"type:varchar(100)" json:"state"`                         // State
	PostalCode         string     `gorm:"type:varchar(20)" json:"postal_code"`                    // Postal code
	Country            string     `gorm:"type:varchar(100)" json:"country"`                       // Country
	PlanID             uint       `gorm:"index" json:"plan_id"`                                   // Foreign key to Subscription Plan table
	PlanStartDate      time.Time `gorm:"type:date" json:"plan_start_date"`                       // Subscription start date
	PlanEndDate        time.Time `gorm:"type:date" json:"plan_end_date"`                         // Subscription end date
	PlanStatus         string     `gorm:"type:varchar(20);default:'active'" json:"plan_status"`   // Subscription status
	MaxEmployees       int        `gorm:"type:int;default:0" json:"max_employees"`                // Max employees allowed
	CompanyLogo        int        `gorm:"type:int" json:"company_logo"`                           // URL for company logo
	Industry           string     `gorm:"type:varchar(100)" json:"industry"`                      // Industry type
	BillingAddress     string     `gorm:"type:text" json:"billing_address"`                       // Billing address
	CompanySize        string     `gorm:"type:varchar(50)" json:"company_size"`
	BossID             uint       `gorm:"column:boss_id;index;foreignKey:UserID;references:ID" json:"boss_id"`
}
