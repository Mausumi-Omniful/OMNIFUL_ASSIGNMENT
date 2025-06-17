package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Mausumi-Omniful/ims/db"
	"github.com/Mausumi-Omniful/ims/models"
	"context"
)

// CreateInventory handles the creation of an inventory item
func CreateInventory(c *gin.Context) {
	var item models.Inventory

	// Bind JSON body to item struct
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the item to the database
	if err := db.DB.GetMasterDB(c.Request.Context()).Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create inventory item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory item created successfully", "item": item})
}

func GetInventories(c *gin.Context) {
	var inventories []models.Inventory
	err := db.DB.GetMasterDB(context.Background()).Find(&inventories).Error
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch inventories"})
		return
	}
	c.JSON(200, inventories)
}

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

// DeleteInventory deletes an inventory item by ID
func DeleteInventory(c *gin.Context) {
	id := c.Param("id")

	// Try to delete the inventory record with matching ID
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


