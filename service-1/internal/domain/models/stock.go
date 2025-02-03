package models

import (
	"errors"
	"time"
)

type Stock struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Quantity    int       `json:"quantity"`
    Unit        string    `json:"unit"`
    MinQuantity int       `json:"min_quantity"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

func NewStock(id, name, description string, quantity int, unit string, minQuantity int) (*Stock, error) {
    if name == "" {
        return nil, errors.New("name is required")
    }
    if quantity < 0 {
        return nil, errors.New("quantity cannot be negative")
    }
    if minQuantity < 0 {
        return nil, errors.New("min quantity cannot be negative")
    }

    now := time.Now()
    return &Stock{
        ID:          id,
        Name:        name,
        Description: description,
        Quantity:    quantity,
        Unit:        unit,
        MinQuantity: minQuantity,
        CreatedAt:   now,
        UpdatedAt:   now,
    }, nil
}

func (s *Stock) UpdateQuantity(adjustment int) error {
    newQuantity := s.Quantity + adjustment
    if newQuantity < 0 {
        return errors.New("insufficient stock")
    }
    s.Quantity = newQuantity
    s.UpdatedAt = time.Now()
    return nil
}

func (s *Stock) IsLowStock() bool {
    return s.Quantity <= s.MinQuantity
}
