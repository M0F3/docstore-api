package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Email        string    `json:"email"`
	TenantId     uuid.UUID    `json:"tenantId"`
	PasswordHash string
}

type UserDTO struct {
	ID        uuid.UUID     `json:"id"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Email     string        `json:"email"`
	Tenant    Tenant `json:"tenant"`
}

type ListUserDTO struct {
	ID        uuid.UUID     `json:"id"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Email     string        `json:"email"`
}

type RegisterUserPayload struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	TenantName string `json:"tenantName"`
	Password string `json:"password"`
}

type RegisteredUserDTO struct {
	ID        uuid.UUID     `json:"id"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Email     string        `json:"email"`
	Tenant    Tenant `json:"tenant"`
}


type LoginPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type SuccessfulLogin struct {
	Token string
	User UserDTO
}