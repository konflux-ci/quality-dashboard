package v1alpha1

import (
	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	UserEmail string    `json:"user_email"`
	Config    string    `json:"user_config"`
}

type Users []User
