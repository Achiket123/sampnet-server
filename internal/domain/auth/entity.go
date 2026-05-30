package auth

import "time"

// User represents a user entity in the domain.
type User struct {
	ID             uint      `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	PhoneNumber    string    `json:"phone_number"`
	IsVerified     bool      `json:"is_verified"`
	HashedPassword string    `json:"-"`
	ProfilePic     string    `json:"profile_pic"`
	City           string    `json:"city"`
	Country        string    `json:"country"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	LastLoginAt    time.Time `json:"last_login_at"`
}
