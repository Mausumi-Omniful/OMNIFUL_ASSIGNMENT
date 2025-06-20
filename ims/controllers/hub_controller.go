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

// CreateHub
func CreateHub(c *gin.Context) {
	var hub models.Hub

	
	err := c.ShouldBindJSON(&hub)
	if err!=nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result:=db.DB.GetMasterDB(c.Request.Context()).Create(&hub) 
	if result.Error !=nil{
		c.JSON(500, gin.H{"error": "Failed to create hub"})
		return
	}
	go func() {
		var updated []models.Hub
		db.DB.GetMasterDB(context.Background()).Find(&updated)
		data, err:= json.Marshal(updated)
		if err== nil {
			redisclient.Client.Set(context.Background(),"All_hubs",string(data),1*time.Hour)
		}
	}()
	c.JSON(200, gin.H{"message":"Hub created","hub": hub})
}





// GetHubs
func GetHubs(c *gin.Context) {
	var hubs []models.Hub
	ctx:= context.Background()
	cacheKey:="All_hubs"

	cached,err:= redisclient.Client.Get(ctx, cacheKey)
	if err==nil {
		if err := json.Unmarshal([]byte(cached), &hubs); err == nil {
			c.JSON(200, gin.H{
				"source": "cache",
				"data": hubs,
			})
			return
		}
	}

	query := db.DB.GetMasterDB(c.Request.Context())

	tenantID:= c.Query("tenant_id")
	if tenantID!= "" {
		query= query.Where("tenant_id = ?", tenantID)
	}
	sellerID := c.Query("seller_id")
	if sellerID !="" {
		query= query.Where("seller_id = ?", sellerID)
	}

     res:=query.Find(&hubs)
	 if res.Error!= nil {
		c.JSON(500,gin.H{"error":"Failed to fetch hubs"})
		return
	}

	
	data,err:= json.Marshal(hubs)
	if err==nil{
	  _, _ =redisclient.Client.Set(ctx, cacheKey, string(data), 1*time.Hour)
	}


	c.JSON(200, gin.H{
		"source": "database",
		"data": hubs,
	})
}




// UpdateHub
func UpdateHub(c *gin.Context) {
	id := c.Param("id")
	var hub models.Hub
	find:= db.DB.GetMasterDB(context.Background()).First(&hub, id)
	if find.Error!= nil {
		c.JSON(400,gin.H{"error": "Hub not found"})
		return
	}
	var updated models.Hub
	update:= c.ShouldBindJSON(&updated) 
	if update!= nil {
		c.JSON(400,gin.H{"error": update.Error()})
		return
	}

	hub.Name = updated.Name
	hub.Location = updated.Location
	hub.TenantID = updated.TenantID
	hub.SellerID = updated.SellerID

	res:= db.DB.GetMasterDB(context.Background()).Save(&hub) 
	if res.Error!= nil {
		c.JSON(500, gin.H{"error": "Failed to update hub"})
		return
	}

	go func() {
		var updated []models.Hub
		db.DB.GetMasterDB(context.Background()).Find(&updated)
		data, err:= json.Marshal(updated)
		if err== nil {
			redisclient.Client.Set(context.Background(),"All_hubs",string(data),1*time.Hour)
		}
	}()
	c.JSON(200,hub)
}




// DeleteHub
func DeleteHub(c *gin.Context) {
	id := c.Param("id")
	res:= db.DB.GetMasterDB(context.Background()).Delete(&models.Hub{}, id)
	if res.Error!= nil {
		c.JSON(500,gin.H{"error": "Failed to delete hub"})
		return
	}

	go func() {
		var updated []models.Hub
		db.DB.GetMasterDB(context.Background()).Find(&updated)
		data, err:= json.Marshal(updated)
		if err== nil {
			redisclient.Client.Set(context.Background(),"All_hubs",string(data),1*time.Hour)
		}
	}()
	c.JSON(200,gin.H{"message": "Hub deleted"})
}
