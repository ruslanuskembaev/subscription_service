package models

import (
	"time"

	"github.com/google/uuid"
)

// Subscription represents a user subscription record
type Subscription struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ServiceName string     `gorm:"type:varchar(255);not null" json:"service_name"`
	Price       int        `gorm:"type:integer;not null;check:price > 0" json:"price"` // price in rubles
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	StartDate   time.Time  `gorm:"type:date;not null;index" json:"start_date"`
	EndDate     *time.Time `gorm:"type:date;index" json:"end_date,omitempty"`
	CreatedAt   time.Time  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// CreateSubscriptionRequest represents request body for creating subscription
type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required,min=1"`
	UserID      string `json:"user_id" binding:"required,uuid"`
	StartDate   string `json:"start_date" binding:"required"` // format: "MM-YYYY"
}

// UpdateSubscriptionRequest represents request body for updating subscription
type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty"`
	StartDate   *string `json:"start_date,omitempty"` // format: "MM-YYYY"
	EndDate     *string `json:"end_date,omitempty"`   // format: "MM-YYYY"
}

// TotalCostRequest represents request body for calculating total cost
type TotalCostRequest struct {
	StartDate   string  `json:"start_date" binding:"required"` // format: "MM-YYYY"
	EndDate     string  `json:"end_date" binding:"required"`   // format: "MM-YYYY"
	UserID      *string `json:"user_id,omitempty"`             // optional filter
	ServiceName *string `json:"service_name,omitempty"`        // optional filter
}

// TotalCostResponse represents response for total cost calculation
type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
}
