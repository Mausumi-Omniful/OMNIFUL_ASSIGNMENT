package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mausumi-Omniful/ims/db"
	"github.com/Mausumi-Omniful/ims/models"
	"github.com/Mausumi-Omniful/ims/redisclient"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func invalidateInventoryCache(ctx context.Context) {
	
	cacheKeys := []string{
		"inventory:all",
		"inventory::",
		"inventory:s003:mausmi",
		"inventory:s003:",
		"inventory::mausmi",
	}

	for _, key := range cacheKeys {
		redisclient.Client.Del(ctx, key)
	}

	
	redisclient.Client.Set(ctx, "inventory:clear", "", 1*time.Second)
}

// CreateInventory
func CreateInventory(c *gin.Context) {
	var item models.Inventory
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res := db.DB.GetMasterDB(c.Request.Context()).Create(&item)
	if res.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create inventory item"})
		return
	}

	go invalidateInventoryCache(context.Background())

	c.JSON(200, gin.H{
		"message": "Inventory item created successfully",
		"item":    item,
	})
}

// GetInventories
func GetInventories(c *gin.Context) {
	var inv []models.Inventory

	ctx := context.Background()
    sku := c.Query("sku")
	location := c.Query("location")

	
	var cacheKey string
	if sku == "" && location == "" {
		cacheKey = "inventory:all"
	} else {
		cacheKey = fmt.Sprintf("inventory:%s:%s", sku, location)
	}
	cached, err := redisclient.Client.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		if err := json.Unmarshal([]byte(cached), &inv); err == nil {
			c.JSON(200, gin.H{
				"source": "cache",
				"data":   inv,
			})
			return
		}
	}

	query := db.DB.GetMasterDB(c.Request.Context())
	if sku != "" {
		query = query.Where("sku = ?", sku)
	}
	if location != "" {
		query = query.Where("location = ?", location)
	}
	query = query.Order("id ASC")

	res := query.Find(&inv)
	if res.Error != nil {
		c.JSON(500, gin.H{"error": res.Error.Error()})
		return
	}
	for i := range inv {
		if inv[i].Quantity < 0 {
			inv[i].Quantity = 0
		}
		if inv[i].SKU == "" {
			inv[i].SKU = ""
		}
		if inv[i].Location == "" {
			inv[i].Location = ""
		}
	}
	data, err := json.Marshal(inv)
	if err == nil {
		_, _ = redisclient.Client.Set(ctx, cacheKey, string(data), 1*time.Hour)
	}

	c.JSON(200, gin.H{
		"source": "database",
		"data":   inv,
	})
}






// UpdateInventory
func UpdateInventory(c *gin.Context) {
	id := c.Param("id")

	var inventory models.Inventory
	res := db.DB.GetMasterDB(context.Background()).First(&inventory, id)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		return
	}

	var updatedData models.Inventory
	err := c.ShouldBindJSON(&updatedData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inventory.ProductID = updatedData.ProductID
	inventory.SKU = updatedData.SKU
	inventory.Quantity = updatedData.Quantity
	inventory.Location = updatedData.Location

	save := db.DB.GetMasterDB(context.Background()).Save(&inventory)
	if save.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go invalidateInventoryCache(context.Background())

	c.JSON(http.StatusOK, inventory)
}




// DeleteInventory
func DeleteInventory(c *gin.Context) {
	id := c.Param("id")

	result := db.DB.GetMasterDB(context.Background()).Delete(&models.Inventory{}, id)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to delete inventory"})
		return
	}

	go invalidateInventoryCache(context.Background())

	c.JSON(200, gin.H{"message": "Inventory deleted successfully"})
}




// UpsertInventory
func UpsertInventory(c *gin.Context) {
	var inv models.Inventory

	err := c.ShouldBindJSON(&inv)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	dbconn := db.DB.GetMasterDB(c.Request.Context())

	res := dbconn.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "sku"},
			{Name: "location"},
		},
		UpdateAll: true,
	}).Create(&inv)

	if res.Error != nil {
		c.JSON(500, gin.H{
			"error": "Failed to upsert",
		})
		return
	}
	go invalidateInventoryCache(context.Background())

	c.JSON(200, gin.H{
		"message": "Inventory upserted successfully",
		"item":    inv,
	})
}






// ReduceInventory
func ReduceInventory(c *gin.Context) {
	var request struct {
		SKU      string `json:"sku" binding:"required"`
		Location string `json:"location" binding:"required"`
		Quantity int    `json:"quantity" binding:"required,gt=0"`
	}

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx := db.DB.GetMasterDB(c.Request.Context()).Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var inventory models.Inventory
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("sku = ? AND location = ?", request.SKU, request.Location).
		First(&inventory)

	if result.Error != nil {
		tx.Rollback()
		if result.Error.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Inventory not found for SKU: " + request.SKU + " at location: " + request.Location,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}
	if inventory.Quantity < request.Quantity {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Insufficient inventory",
			"details": gin.H{
				"requested": request.Quantity,
				"available": inventory.Quantity,
				"sku":       request.SKU,
				"location":  request.Location,
			},
		})
		return
	}
	newQuantity:= inventory.Quantity-request.Quantity
	updateResult := tx.Model(&inventory).Update("quantity", newQuantity)
	if updateResult.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	inventory.Quantity = newQuantity
	go invalidateInventoryCache(context.Background())

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory reduced successfully",
		"details": gin.H{
			"sku":               request.SKU,
			"location":          request.Location,
			"reduced_by":        request.Quantity,
			"new_quantity":      newQuantity,
			"previous_quantity": inventory.Quantity + request.Quantity,
		},
		"inventory": inventory,
	})
}
