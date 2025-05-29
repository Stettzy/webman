package models

import (
	"time"
)

// Collection represents a group of related requests
type Collection struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Requests    []Request `json:"requests" gorm:"foreignKey:CollectionID"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Request represents a saved API request
type Request struct {
	ID           string            `json:"id" gorm:"primaryKey;type:uuid"`
	CollectionID string            `json:"collection_id" gorm:"type:uuid;not null"`
	Name         string            `json:"name" gorm:"not null"`
	Method       string            `json:"method" gorm:"not null"`
	URL          string            `json:"url" gorm:"not null"`
	Headers      map[string]string `json:"headers" gorm:"serializer:json"`
	Body         string            `json:"body"`
	BodyType     string            `json:"body_type"`
	CreatedAt    time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
}

// DefaultHeader represents a commonly used HTTP header
type DefaultHeader struct {
	ID          string `json:"id" gorm:"primaryKey;type:uuid"`
	Name        string `json:"name" gorm:"not null"`
	Value       string `json:"value"`
	Description string `json:"description"`
}
