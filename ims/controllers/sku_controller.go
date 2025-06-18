package controllers

import (
	"time"
	"context"
	"net/http"
	"fmt"
    
	"github.com/Mausumi-Omniful/ims/db"
	"github.com/Mausumi-Omniful/ims/models"
	"github.com/gin-gonic/gin"
	"github.com/Mausumi-Omniful/ims/redisclient"
)



// CreateSKU
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

	// cache in Redis
	cacheKey := fmt.Sprintf("sku_code:%s", sku.Code)
    _, _ = redisclient.Client.Set(c.Request.Context(), cacheKey, "yes", time.Hour)



	c.JSON(http.StatusOK, gin.H{"message": "SKU created", "sku": sku})
}





// GetSKUs
func GetSKUs(c *gin.Context) {
	var skus []models.SKU
	query := db.DB.GetMasterDB(c.Request.Context())

	if tenantID := c.Query("tenant_id"); tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if sellerID := c.Query("seller_id"); sellerID != "" {
		query = query.Where("seller_id = ?", sellerID)
	}
	if skuCode := c.Query("sku_code"); skuCode != "" {
		query = query.Where("sku_code = ?", skuCode)
	}

	if err := query.Find(&skus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SKUs"})
		return
	}

	c.JSON(http.StatusOK, skus)
}




// UpdateSKU
func UpdateSKU(c *gin.Context) {
	id := c.Param("id")
	var sku models.SKU
	if err := db.DB.GetMasterDB(context.Background()).First(&sku, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SKU not found"})
		return
	}
	var updated models.SKU
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sku.Code = updated.Code
	sku.Name = updated.Name
	sku.SKUCode = updated.SKUCode
	sku.Description = updated.Description
	sku.TenantID = updated.TenantID
	sku.SellerID = updated.SellerID

	if err := db.DB.GetMasterDB(context.Background()).Save(&sku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update SKU"})
		return
	}
	c.JSON(http.StatusOK, sku)
}



// DeleteSKU
func DeleteSKU(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.GetMasterDB(context.Background()).Delete(&models.SKU{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SKU"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "SKU deleted"})
}



