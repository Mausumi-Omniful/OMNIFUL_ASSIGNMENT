package controllers

import (
	"context"
	"net/http"
     "fmt"
	 "time"
	"github.com/Mausumi-Omniful/ims/db"
	"github.com/Mausumi-Omniful/ims/models"
	"github.com/gin-gonic/gin"
	"github.com/Mausumi-Omniful/ims/redisclient"
	

)

// CreateHub
func CreateHub(c *gin.Context) {
	var hub models.Hub

	
	if err := c.ShouldBindJSON(&hub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save to Postgres
	if err := db.DB.GetMasterDB(c.Request.Context()).Create(&hub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hub"})
		return
	}

	
	// Cache hub location in Redis
    cacheKey := fmt.Sprintf("hub_location:%s", hub.Location)
    _, _ = redisclient.Client.Set(c.Request.Context(), cacheKey, "true", time.Hour)




	c.JSON(http.StatusOK, gin.H{"message": "Hub created", "hub": hub})
}





// GetHubs
func GetHubs(c *gin.Context) {
	var hubs []models.Hub
	query := db.DB.GetMasterDB(c.Request.Context())

	
	if tenantID := c.Query("tenant_id"); tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if sellerID := c.Query("seller_id"); sellerID != "" {
		query = query.Where("seller_id = ?", sellerID)
	}

	if err := query.Find(&hubs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hubs"})
		return
	}

	c.JSON(http.StatusOK, hubs)
}



// GetHubByID
func GetHubByID(c *gin.Context) {
	id := c.Param("id")
	var hub models.Hub
	if err := db.DB.GetMasterDB(context.Background()).First(&hub, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hub not found"})
		return
	}
	c.JSON(http.StatusOK, hub)
}



// UpdateHub
func UpdateHub(c *gin.Context) {
	id := c.Param("id")
	var hub models.Hub
	if err := db.DB.GetMasterDB(context.Background()).First(&hub, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hub not found"})
		return
	}
	var updated models.Hub
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hub.Name = updated.Name
	hub.Location = updated.Location
	hub.TenantID = updated.TenantID
	hub.SellerID = updated.SellerID

	if err := db.DB.GetMasterDB(context.Background()).Save(&hub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hub"})
		return
	}
	c.JSON(http.StatusOK, hub)
}




// DeleteHub
func DeleteHub(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.GetMasterDB(context.Background()).Delete(&models.Hub{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete hub"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Hub deleted"})
}
