package models

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Extension string `json:"extension"`
	Size int `json:"size"`
	UploadedAt  time.Time `json:",format:datetime"`
	TenantId uuid.UUID `json:"tenant_id"`
}