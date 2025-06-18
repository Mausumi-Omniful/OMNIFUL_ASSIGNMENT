package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Mausumi-Omniful/ims/db"
	"github.com/Mausumi-Omniful/ims/models"
	"context"
	
	"gorm.io/gorm/clause"
	"github.com/Mausumi-Omniful/ims/redisclient"
)



// CreateInventory
func CreateInventory(c *gin.Context) {
	var item models.Inventory

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	
	if err := db.DB.GetMasterDB(c.Request.Context()).Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create inventory item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory item created successfully", "item": item})
}



// GetInventories
func GetInventories(c *gin.Context) {
    sku := c.Query("sku")
    location := c.Query("location")

    var inv models.Inventory
    err := db.DB.GetMasterDB(c.Request.Context()).Where("sku = ? AND location = ?", sku, location).First(&inv).Error

    if err != nil {
        // default to quantity 0 if not found
        c.JSON(http.StatusOK, gin.H{
            "sku":      sku,
            "location": location,
            "quantity": 0,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "sku":      inv.SKU,
        "location": inv.Location,
        "quantity": inv.Quantity,
    })
}


// GetInventoryByID
func GetInventoryByID(c *gin.Context) {
	id := c.Param("id")
	var inventory models.Inventory
	err := db.DB.GetMasterDB(context.Background()).First(&inventory, id).Error
	if err != nil {
		c.JSON(404, gin.H{"error": "Inventory not found"})
		return
	}
	c.JSON(200, inventory)
}



// UpdateInventory
func UpdateInventory(c *gin.Context) {
	id := c.Param("id")

	var inventory models.Inventory
	if err := db.DB.GetMasterDB(context.Background()).First(&inventory, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		return
	}

	var updatedData models.Inventory
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inventory.ProductID = updatedData.ProductID
	inventory.SKU = updatedData.SKU
	inventory.Quantity = updatedData.Quantity
	inventory.Location = updatedData.Location

	if err := db.DB.GetMasterDB(context.Background()).Save(&inventory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventory)
}



// DeleteInventory
func DeleteInventory(c *gin.Context) {
	id := c.Param("id")

	
	result := db.DB.GetMasterDB(context.Background()).Delete(&models.Inventory{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete inventory"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory deleted successfully"})
}





//UpsertInventory
func UpsertInventory(c *gin.Context) {
	var inv models.Inventory
	ctx := c.Request.Context()

	
	if err := c.ShouldBindJSON(&inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate SKU from Redis
	skuKey := "sku_code:" + inv.SKU
	skuExists, err := redisclient.Client.Exists(ctx, skuKey)
	if err != nil || skuExists == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SKU"})
		return
	}

	// Validate Location from Redis
	hubKey := "hub_location:" + inv.Location
	hubExists, err := redisclient.Client.Exists(ctx, hubKey)
	if err != nil || hubExists == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Location"})
		return
	}

	// Atomic upsert
	err = db.DB.GetMasterDB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "sku"}, {Name: "location"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity", "updated_at"}),
	}).Create(&inv).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert inventory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory upserted", "inventory": inv})
}
