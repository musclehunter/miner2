package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musclehunter/miner2/database"
	"github.com/musclehunter/miner2/models"
)



// GetMyInventoryHandler retrieves the player's warehouse and its contents.
// It replaces the old inventory system with the new warehouse-based one.
func GetMyInventoryHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}



	// Get the user's warehouse and the items within it
	warehouse, items, err := database.GetWarehouseAndItemsByUserID(db, userID.(string))
	if err != nil {
		log.Printf("Failed to get warehouse and items for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve warehouse inventory"})
		return
	}

	// If warehouse is nil, it means the user has no base yet.
	if warehouse == nil {
		c.JSON(http.StatusOK, gin.H{
			"warehouse": nil,
			"items":     []models.WarehouseItem{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"warehouse": warehouse,
		"items":     items,
	})
}
