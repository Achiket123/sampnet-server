package settings

import (
	"errors"
	"time"
)

var (
	ErrSettingsNotFound  = errors.New("settings not found")
	ErrPlanNotFound      = errors.New("subscription plan not found")
	ErrPolicyNotFound    = errors.New("policy not found")
	ErrNotAuthorized     = errors.New("not authorized to perform this action")
	ErrInvalidInput      = errors.New("invalid input data")
	ErrOrganisationNotFound = errors.New("organisation not found")
)

type OrgSettings struct {
	ID                 uint   `json:"id"`
	CompanyName        string `json:"company_name"`
	CompanyCode        string `json:"company_code"`
	PrimaryContactName string `json:"primary_contact_name"`
	PrimaryEmail       string `json:"primary_email"`
	PhoneNumber        string `json:"phone_number"`
	OfficeAddress      string `json:"office_address"`
	City               string `json:"city"`
	State              string `json:"state"`
	PostalCode         string `json:"postal_code"`
	Country            string `json:"country"`
	CompanyLogo        *int   `json:"company_logo"`
	Industry           string `json:"industry"`
	BillingAddress     string `json:"billing_address"`
	CompanySize        string `json:"company_size"`
}

type PlanInfo struct {
	PlanID        uint      `json:"plan_id"`
	PlanName      string    `json:"plan_name"`
	PlanStatus    string    `json:"plan_status"`
	MaxEmployees  int       `json:"max_employees"`
	PlanStartDate time.Time `json:"plan_start_date"`
	PlanEndDate   time.Time `json:"plan_end_date"`
}

type RolePermissionsSet struct {
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type AttendancePolicy struct {
	ID              uint   `json:"id"`
	OrganisationID  uint   `json:"organisation_id"`
	CheckInTime     string `json:"check_in_time"`
	CheckOutTime    string `json:"check_out_time"`
	GracePeriodMins int    `json:"grace_period_mins"`
}

type LeavePolicyEntry struct {
	ID             uint   `json:"id"`
	OrganisationID uint   `json:"organisation_id"`
	LeaveType      string `json:"leave_type"`
	MaxDays        int    `json:"max_days"`
	Description    string `json:"description"`
}

type TaskTypeEntry struct {
	ID             uint   `json:"id"`
	OrganisationID uint   `json:"organisation_id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
}

type UserProfile struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	DateOfBirth time.Time `json:"date_of_birth"`
	ProfilePic  string    `json:"profile_id"`
}

type NotificationPreferenceEntry struct {
	Category string `json:"category"`
	Email    bool   `json:"email"`
	Push     bool   `json:"push"`
	InApp    bool   `json:"in_app"`
}