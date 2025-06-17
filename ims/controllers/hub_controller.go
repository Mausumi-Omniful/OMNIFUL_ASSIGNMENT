package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Mausumi-Omniful/ims/models"
	"github.com/Mausumi-Omniful/ims/db"
	"context"
)

func CreateHub(c *gin.Context) {
	var hub models.Hub
	if err := c.ShouldBindJSON(&hub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.GetMasterDB(context.Background()).Create(&hub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Hub"})
		return
	}
	c.JSON(http.StatusOK, hub)
}

func GetAllHubs(c *gin.Context) {
	var hubs []models.Hub
	if err := db.DB.GetMasterDB(context.Background()).Find(&hubs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hubs"})
		return
	}
	c.JSON(http.StatusOK, hubs)
}
