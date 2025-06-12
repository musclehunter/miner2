package models

import "time"

// WarehouseItem represents an item or ore stored in a warehouse
// It corresponds to the `warehouse_items` table
// Note: This model can store both items and ores by referencing their respective IDs.
// A more robust implementation might use a polymorphic association or separate tables,
// but for simplicity, we'll use nullable fields for ItemID and OreID.
type WarehouseItem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	WarehouseID uint      `json:"warehouse_id" gorm:"not null"`
	ItemID      *uint     `json:"item_id" gorm:"null"` // Pointer to allow null
	OreID       *uint     `json:"ore_id" gorm:"null"`  // Pointer to allow null
	Quantity    int       `json:"quantity" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Eager loading associations
	Item Item `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Ore  Ore  `json:"ore,omitempty" gorm:"foreignKey:OreID"`
}
