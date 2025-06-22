package controllers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Mausumi-Omniful/ims/db"
	"github.com/Mausumi-Omniful/ims/models"
	"github.com/Mausumi-Omniful/ims/redisclient"
	"github.com/gin-gonic/gin"
)

// CreateSKU
func CreateSKU(c *gin.Context) {
	var sku models.SKU
	err := c.ShouldBindJSON(&sku)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	res := db.DB.GetMasterDB(context.Background()).Create(&sku)
	if res.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create SKU"})
		return
	}

	go func() {
		var updated []models.SKU
		db.DB.GetMasterDB(context.Background()).Find(&updated)
		data, err := json.Marshal(updated)
		if err == nil {
			redisclient.Client.Set(context.Background(), "All_skus", string(data), 1*time.Hour)
		}
	}()

	c.JSON(200, gin.H{"message": "SKU created", "sku": sku})
}

// GetSKUs
func GetSKUs(c *gin.Context) {
	var skus []models.SKU
	ctx := context.Background()
	cacheKey := "All_skus"
	cached, err := redisclient.Client.Get(ctx, cacheKey)
	if err == nil {
		err := json.Unmarshal([]byte(cached), &skus)
		if err == nil {
			c.JSON(200, gin.H{
				"source": "cache",
				"data":   skus,
			})
			return
		}
	}
	query := db.DB.GetMasterDB(c.Request.Context())

	tenantID := c.Query("tenant_id")
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	sellerID := c.Query("seller_id")
	if sellerID != "" {
		query = query.Where("seller_id = ?", sellerID)
	}
	skuCode := c.Query("sku_code")
	if skuCode != "" {
		query = query.Where("sku_code = ?", skuCode)
	}

	res := query.Find(&skus)
	if res.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch SKUs"})
		return
	}

	data, err := json.Marshal(skus)
	if err == nil {
		_, _ = redisclient.Client.Set(ctx, cacheKey, string(data), 1*time.Hour)
	}

	c.JSON(200, gin.H{
		"source": "database",
		"data":   skus,
	})
}

// UpdateSKU
func UpdateSKU(c *gin.Context) {
	id := c.Param("id")
	var sku models.SKU
	if err := db.DB.GetMasterDB(context.Background()).First(&sku, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "SKU not found"})
		return
	}
	var updated models.SKU
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	sku.Code = updated.Code
	sku.Name = updated.Name
	sku.Description = updated.Description
	sku.TenantID = updated.TenantID
	sku.SellerID = updated.SellerID

	if err := db.DB.GetMasterDB(context.Background()).Save(&sku).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update SKU"})
		return
	}

	go func() {
		var updated []models.SKU
		db.DB.GetMasterDB(context.Background()).Find(&updated)
		data, err := json.Marshal(updated)
		if err == nil {
			redisclient.Client.Set(context.Background(), "All_skus", string(data), 1*time.Hour)
		}
	}()
	c.JSON(200, sku)
}

// DeleteSKU
func DeleteSKU(c *gin.Context) {
	id := c.Param("id")
	res := db.DB.GetMasterDB(context.Background()).Delete(&models.SKU{}, id)
	if res.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to delete SKU"})
		return
	}

	go func() {
		var updated []models.SKU
		db.DB.GetMasterDB(context.Background()).Find(&updated)
		data, err := json.Marshal(updated)
		if err == nil {
			redisclient.Client.Set(context.Background(), "All_skus", string(data), 1*time.Hour)
		}
	}()
	c.JSON(200, gin.H{"message": "SKU deleted"})
}
