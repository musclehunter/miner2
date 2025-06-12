package database

import (
	"database/sql"
	"strconv"

	"github.com/musclehunter/miner2/models"
)

// GetWarehouseAndItemsByUserID finds a user's warehouse and its contents.
func GetWarehouseAndItemsByUserID(db *sql.DB, userID string) (*models.Warehouse, []models.WarehouseItem, error) {
	// 1. Find the user's base and warehouse
	warehouse := &models.Warehouse{}
	queryWarehouse := `
		SELECT w.id, w.base_id, w.level, w.capacity, w.created_at, w.updated_at
		FROM warehouses w
		JOIN bases b ON w.base_id = b.id
		WHERE b.user_id = ?
		LIMIT 1
	`
	err := db.QueryRow(queryWarehouse, userID).Scan(
		&warehouse.ID, &warehouse.BaseID, &warehouse.Level, &warehouse.Capacity, &warehouse.CreatedAt, &warehouse.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []models.WarehouseItem{}, nil // No warehouse found, which is a valid case
		}
		return nil, nil, err
	}

	// 2. Find all items in that warehouse
	queryItems := `
		SELECT wi.id, wi.warehouse_id, wi.item_id, wi.ore_id, wi.quantity, wi.created_at, wi.updated_at,
		       i.name, i.rarity, i.description,
		       o.name, o.rarity, o.purity, o.processing_difficulty
		FROM warehouse_items wi
		LEFT JOIN items i ON wi.item_id = i.id
		LEFT JOIN ores o ON wi.ore_id = o.id
		WHERE wi.warehouse_id = ?
	`
	rows, err := db.Query(queryItems, warehouse.ID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var items []models.WarehouseItem
	for rows.Next() {
		item := models.WarehouseItem{}
		// Need to scan into nullable fields for item/ore details
		var itemName, itemRarity, itemDescription, oreName, orePurity, oreProcessingDifficulty sql.NullString
		var oreRarity sql.NullInt32

		err := rows.Scan(
			&item.ID, &item.WarehouseID, &item.ItemID, &item.OreID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
			&itemName, &itemRarity, &itemDescription,
			&oreName, &oreRarity, &orePurity, &oreProcessingDifficulty,
		)
		if err != nil {
			return nil, nil, err
		}

		if item.ItemID != nil {
			rarity, _ := strconv.Atoi(itemRarity.String)
			item.Item = models.Item{
				Name:        itemName.String,
				Rarity:      rarity,
				Description: itemDescription.String,
			}
		}
		if item.OreID != nil {
			purity, _ := strconv.ParseFloat(orePurity.String, 64)
			difficulty, _ := strconv.Atoi(oreProcessingDifficulty.String)
			item.Ore = models.Ore{
				Name:                 oreName.String,
				Rarity:               int(oreRarity.Int32),
				Purity:               purity,
				ProcessingDifficulty: difficulty,
			}
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	// Ensure we return an empty slice, not null
	if items == nil {
		items = []models.WarehouseItem{}
	}

	return warehouse, items, nil
}
