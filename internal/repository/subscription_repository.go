package repository

import (
	"errors"
	"fmt"
	"time"

	"subscription_service/internal/database"
	"subscription_service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubscriptionRepository handles database operations for subscriptions
type SubscriptionRepository struct {
	db *database.DB
}

// NewSubscriptionRepository creates new subscription repository
func NewSubscriptionRepository(db *database.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create creates a new subscription
func (r *SubscriptionRepository) Create(sub *models.Subscription) error {
	if err := r.db.Create(sub).Error; err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	return nil
}

// GetByID retrieves subscription by ID
func (r *SubscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	sub := &models.Subscription{}
	if err := r.db.Where("id = ?", id).First(sub).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	return sub, nil
}

// Update updates subscription
func (r *SubscriptionRepository) Update(id uuid.UUID, sub *models.Subscription) error {
	result := r.db.Model(&models.Subscription{}).Where("id = ?", id).Updates(sub)
	if result.Error != nil {
		return fmt.Errorf("failed to update subscription: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}
	return nil
}

// Delete deletes subscription by ID
func (r *SubscriptionRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&models.Subscription{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete subscription: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}
	return nil
}

// List retrieves all subscriptions with optional pagination
func (r *SubscriptionRepository) List(limit, offset int) ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription
	if err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&subscriptions).Error; err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	return subscriptions, nil
}

// FindActiveSubscriptions finds subscriptions active in the given period
func (r *SubscriptionRepository) FindActiveSubscriptions(startDate, endDate time.Time, userID *uuid.UUID, serviceName *string) ([]*models.Subscription, error) {
	query := r.db.Model(&models.Subscription{}).
		Where("start_date <= ?", endDate).
		Where("end_date IS NULL OR end_date >= ?", startDate)

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if serviceName != nil {
		query = query.Where("service_name = ?", *serviceName)
	}

	var subscriptions []*models.Subscription
	if err := query.Find(&subscriptions).Error; err != nil {
		return nil, fmt.Errorf("failed to find active subscriptions: %w", err)
	}

	return subscriptions, nil
}

// CalculateTotalCost calculates total cost for a period with optional filters
// It correctly calculates cost by multiplying price by number of months in the intersection
func (r *SubscriptionRepository) CalculateTotalCost(startDate, endDate time.Time, userID *uuid.UUID, serviceName *string) (int, error) {
	// Find all subscriptions that are active in the period
	subscriptions, err := r.FindActiveSubscriptions(startDate, endDate, userID, serviceName)
	if err != nil {
		return 0, err
	}

	totalCost := 0

	// For each subscription, calculate the intersection with the requested period
	for _, sub := range subscriptions {
		// Find intersection start: max of subscription start and requested start
		// Normalize to first day of month
		subStart := time.Date(sub.StartDate.Year(), sub.StartDate.Month(), 1, 0, 0, 0, 0, sub.StartDate.Location())
		reqStart := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())

		intersectionStart := subStart
		if reqStart.After(intersectionStart) {
			intersectionStart = reqStart
		}

		// Find intersection end: min of subscription end (or endDate if nil) and requested end
		intersectionEnd := endDate
		if sub.EndDate != nil {
			// Get last day of subscription end month
			subEndLastDay := time.Date(sub.EndDate.Year(), sub.EndDate.Month()+1, 0, 23, 59, 59, 999999999, sub.EndDate.Location())
			if subEndLastDay.Before(endDate) {
				intersectionEnd = subEndLastDay
			}
		}

		// Normalize intersectionEnd to first day of month for calculation
		intersectionEndFirstDay := time.Date(intersectionEnd.Year(), intersectionEnd.Month(), 1, 0, 0, 0, 0, intersectionEnd.Location())

		// Calculate number of months in intersection
		months := countMonths(intersectionStart, intersectionEndFirstDay)

		// Multiply price by number of months
		totalCost += sub.Price * months
	}

	return totalCost, nil
}

// countMonths calculates the number of months between two dates (inclusive)
// Both dates should be normalized to first day of month
func countMonths(start, end time.Time) int {
	if end.Before(start) {
		return 0
	}

	// Calculate difference in months
	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	totalMonths := years*12 + months + 1 // +1 because both months are inclusive

	return totalMonths
}
