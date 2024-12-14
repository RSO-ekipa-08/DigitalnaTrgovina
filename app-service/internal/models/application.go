package models

import (
	"time"
)

type Application struct {
	ID             string    `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Description    string    `json:"description" db:"description"`
	DeveloperID    string    `json:"developer_id" db:"developer_id"`
	Category       string    `json:"category" db:"category"`
	Price          float64   `json:"price" db:"price"`
	Size           int64     `json:"size" db:"size"`
	MinAndroidVer  string    `json:"min_android_version" db:"min_android_version"`
	CurrentVersion string    `json:"current_version" db:"current_version"`
	Tags           []string  `json:"tags" db:"tags"`
	Screenshots    []string  `json:"screenshots" db:"screenshots"`
	Rating         float64   `json:"rating" db:"rating"`
	Downloads      int       `json:"downloads" db:"downloads"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	StorageURL     string    `json:"storage_url" db:"storage_url"`
}

type Download struct {
	ID            string    `json:"id" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	ApplicationID string    `json:"application_id" db:"application_id"`
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
	IPAddress     string    `json:"ip_address" db:"ip_address"`
	Success       bool      `json:"success" db:"success"`
	ErrorMessage  string    `json:"error_message,omitempty" db:"error_message"`
}
