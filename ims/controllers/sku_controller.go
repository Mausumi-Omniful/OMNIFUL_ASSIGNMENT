package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Mausumi-Omniful/ims/models"
	"github.com/Mausumi-Omniful/ims/db"
	"context"
)

func CreateSKU(c *gin.Context) {
	var sku models.SKU
	if err := c.ShouldBindJSON(&sku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.GetMasterDB(context.Background()).Create(&sku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SKU"})
		return
	}
	c.JSON(http.StatusOK, sku)
}

func GetAllSKUs(c *gin.Context) {
	var skus []models.SKU
	if err := db.DB.GetMasterDB(context.Background()).Find(&skus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SKUs"})
		return
	}
	c.JSON(http.StatusOK, skus)
}
