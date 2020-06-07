package entities

import (
	"github.com/google/uuid"
	"time"
)

type Restaurant struct {
	Interface
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name string
	CreatedAt time.Time
	UpdateAt time.Time
	Orders []Order `gorm:"foreignkey:RestaurantID;association_foreignkey:ID"`
}

func (r *Restaurant) TableName() string {
	return "restaurants"
}

type Order struct {
	Interface
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
	Description string
	RestaurantID uuid.UUID `sql:"type:uuid REFERENCES restaurants(id) ON DELETE CASCADE"`
	CreatedAt time.Time
	UpdateAt time.Time
}

func (o *Order) TableName() string {
	return "orders"
}