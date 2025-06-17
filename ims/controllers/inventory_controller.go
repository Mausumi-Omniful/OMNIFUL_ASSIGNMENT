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


