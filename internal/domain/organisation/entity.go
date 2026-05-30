package organisation

import "time"

// Entity is the core business representation of an organisation.
type Entity struct {
	ID                 uint      `json:"id"`
	CompanyName        string    `json:"company_name"`
	CompanyCode        string    `json:"company_code"`
	PrimaryContactName string    `json:"primary_contact_name"`
	PrimaryEmail       string    `json:"primary_email"`
	PhoneNumber        string    `json:"phone_number"`
	OfficeAddress      string    `json:"office_address"`
	City               string    `json:"city"`
	State              string    `json:"state"`
	PostalCode         string    `json:"postal_code"`
	Country            string    `json:"country"`
	PlanID             uint      `json:"plan_id"`
	PlanStartDate      time.Time `json:"plan_start_date"`
	PlanEndDate        time.Time `json:"plan_end_date"`
	PlanStatus         string    `json:"plan_status"`
	MaxEmployees       int       `json:"max_employees"`
	CompanyLogo        int       `json:"company_logo"`
	Industry           string    `json:"industry"`
	BillingAddress     string    `json:"billing_address"`
	CompanySize        string    `json:"company_size"`
	BossID             uint      `json:"boss_id"`
}

// OwnerEmployeeRow is the employee row created alongside organisation registration.
type OwnerEmployeeRow struct {
	UserID         uint      `json:"user_id"`
	EmploymentID   int       `json:"employment_id"`
	OrganisationID uint      `json:"organisation_id"`
	Type           string    `json:"type"`
	Salary         string    `json:"salary"`
	LastLoginAt    time.Time `json:"last_login_at"`
}
