package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	OrderID   uint           `gorm:"primaryKey;autoIncrement" json:"order_id"`
	TicketID  int            `gorm:"not null" json:"ticket_id"`
	OrderedBy string         `gorm:"type:varchar(100);not null" json:"ordered_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
