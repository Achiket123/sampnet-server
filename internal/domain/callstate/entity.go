package callstate

import "time"

type State struct {
	ID               uint       `json:"id"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Email            string     `json:"email"`
	OrganisationID   uint       `json:"organisation_id"`
	InCall           bool       `json:"in_call"`
	LastCall         *time.Time `json:"last_call"`
	CallingID        *string    `json:"calling_id"`
	CallingFirstName *string    `json:"calling_first_name"`
	CallingLastName  *string    `json:"calling_last_name"`
	Offer            *string    `json:"offer"`
}
