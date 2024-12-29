package models

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	TicketID        uint           `gorm:"primaryKey;autoIncrement" json:"ticket_id"`
	Title           string         `gorm:"type:varchar(100);not null" json:"title"`
	InitialQuantity int            `gorm:"not null" json:"initial_quantity"`
	CurrentQuantity int            `gorm:"not null" json:"current_quantity"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
