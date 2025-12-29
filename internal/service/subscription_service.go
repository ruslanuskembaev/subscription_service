package service

import (
	"fmt"
	"time"

	"subscription_service/internal/models"
	"subscription_service/internal/repository"

	"github.com/google/uuid"
)

// SubscriptionService handles business logic for subscriptions
type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

// NewSubscriptionService creates new subscription service
func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

// CreateSubscription creates a new subscription
func (s *SubscriptionService) CreateSubscription(req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}

	now := time.Now()
	sub := &models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     nil,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(sub); err != nil {
		return nil, err
	}

	return sub, nil
}

// GetSubscription retrieves subscription by ID
func (s *SubscriptionService) GetSubscription(id uuid.UUID) (*models.Subscription, error) {
	return s.repo.GetByID(id)
}

// UpdateSubscription updates subscription
func (s *SubscriptionService) UpdateSubscription(id uuid.UUID, req *models.UpdateSubscriptionRequest) (*models.Subscription, error) {
	sub, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.ServiceName != nil {
		sub.ServiceName = *req.ServiceName
	}
	if req.Price != nil {
		sub.Price = *req.Price
	}
	if req.StartDate != nil {
		startDate, err := parseMonthYear(*req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date format: %w", err)
		}
		sub.StartDate = startDate
	}
	if req.EndDate != nil {
		endDate, err := parseMonthYear(*req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format: %w", err)
		}
		sub.EndDate = &endDate
	}

	sub.UpdatedAt = time.Now()

	if err := s.repo.Update(id, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

// DeleteSubscription deletes subscription by ID
func (s *SubscriptionService) DeleteSubscription(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// ListSubscriptions retrieves all subscriptions
func (s *SubscriptionService) ListSubscriptions(limit, offset int) ([]*models.Subscription, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(limit, offset)
}

// CalculateTotalCost calculates total cost for a period
func (s *SubscriptionService) CalculateTotalCost(req *models.TotalCostRequest) (*models.TotalCostResponse, error) {
	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}

	endDate, err := parseMonthYear(req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format: %w", err)
	}

	// Set end date to last day of the month
	endDate = time.Date(endDate.Year(), endDate.Month()+1, 0, 23, 59, 59, 999999999, endDate.Location())

	var userID *uuid.UUID
	if req.UserID != nil {
		parsed, err := uuid.Parse(*req.UserID)
		if err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		userID = &parsed
	}

	totalCost, err := s.repo.CalculateTotalCost(startDate, endDate, userID, req.ServiceName)
	if err != nil {
		return nil, err
	}

	return &models.TotalCostResponse{TotalCost: totalCost}, nil
}

// parseMonthYear parses date in "MM-YYYY" format
func parseMonthYear(dateStr string) (time.Time, error) {
	layout := "01-2006"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	// Set to first day of the month
	date = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	return date, nil
}
