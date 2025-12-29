package models

import (
	"time"

	"github.com/google/uuid"
)

// Subscription represents a user subscription record
type Subscription struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	ServiceName string     `json:"service_name" db:"service_name"`
	Price       int        `json:"price" db:"price"` // price in rubles
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	StartDate   time.Time  `json:"start_date" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
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

