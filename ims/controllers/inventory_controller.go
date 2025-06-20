package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Mausumi-Omniful/ims/db"
	"github.com/Mausumi-Omniful/ims/models"
	"context"
	"time"
	"encoding/json"
	"gorm.io/gorm/clause"
	"github.com/Mausumi-Omniful/ims/redisclient"
)



// CreateInventory
func CreateInventory(c *gin.Context) {
	var item models.Inventory
     err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(400,gin.H{"error": err.Error()})
		return
	}

	res:= db.DB.GetMasterDB(c.Request.Context()).Create(&item)
	if res.Error!= nil {
		c.JSON(500,gin.H{"error": "Failed to create inventory item"})
		return
	}

	go func() {
		var updated []models.Inventory
		db.DB.GetMasterDB(context.Background()).Find(&updated)
		data, err:=json.Marshal(updated)
		if err==nil {
			redisclient.Client.Set(context.Background(), "All_Inventory", string(data), 1*time.Hour)
		}
	}()
	c.JSON(200, gin.H{
		"message": "Inventory item created successfully", 
		"item": item,
	})
}



// GetInventories
func GetInventories(c *gin.Context) {
    var inv []models.Inventory

	ctx:= context.Background()
	cacheKey:= "All_Inventory"

	cached,err:= redisclient.Client.Get(ctx, cacheKey)
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &inv); err == nil {
			c.JSON(200, gin.H{
				"source": "cache",
				"data": inv,
			})
			return
		}
	}


	res:=db.DB.GetMasterDB(c.Request.Context()).Find(&inv)
	if res.Error!=nil{
		c.JSON(500, gin.H{"error": res.Error.Error()})
		return
	}

	for i:=range inv{
		if inv[i].SKU=="" {
			inv[i].SKU=""
		}
		if inv[i].Location=="" {
			inv[i].Location=""
		}
		if inv[i].Quantity<0 {
			inv[i].Quantity=0
		}
	}
	
	data,err:= json.Marshal(inv)
	if err==nil{
	  _, _ =redisclient.Client.Set(ctx, cacheKey, string(data), 1*time.Hour)
	}


	c.JSON(200, gin.H{
		"source": "database",
		"data": inv,
	})
}



// UpdateInventory
func UpdateInventory(c *gin.Context) {
	id := c.Param("id")

	var inventory models.Inventory
	res:= db.DB.GetMasterDB(context.Background()).First(&inventory, id)
	if res.Error!= nil {
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

	save:= db.DB.GetMasterDB(context.Background()).Save(&inventory)
	if save.Error!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func() {
		var updated []models.Inventory
		db.DB.GetMasterDB(context.Background()).Find(&updated)
		data, err := json.Marshal(updated)
		if err == nil {
			redisclient.Client.Set(context.Background(), "All_Inventory", string(data), 1*time.Hour)
		}
	}()


	c.JSON(http.StatusOK, inventory)
}



// DeleteInventory
func DeleteInventory(c *gin.Context) {
	id := c.Param("id")

	
	result := db.DB.GetMasterDB(context.Background()).Delete(&models.Inventory{}, id)
	if result.Error!= nil {
		c.JSON(500,gin.H{"error": "Failed to delete inventory"})
		return
	}

	go func() {
		var updatedList []models.Inventory
		db.DB.GetMasterDB(context.Background()).Find(&updatedList)
		data, err := json.Marshal(updatedList)
		if err == nil {
			redisclient.Client.Set(context.Background(), "All_Inventory", string(data), 1*time.Hour)
		}
	}()

	c.JSON(200,gin.H{"message": "Inventory deleted successfully"})
}





//UpsertInventory
func UpsertInventory(c *gin.Context) {
	var inv models.Inventory
	
	err:=c.ShouldBindJSON(&inv)
	if err!=nil{
		c.JSON(400,gin.H{
			"error":err.Error(),
		})
	}

	dbconn:=db.DB.GetMasterDB(c.Request.Context())

	res:=dbconn.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name:"sku"},
			{Name:"location"},
		},
		UpdateAll: true,
	}).Create(&inv)

	if res.Error!=nil{
		c.JSON(500,gin.H{
			"error":"Failed to upsert",
		})
		return 
	}
	go func(){
		var  updated []models.Inventory

		db.DB.GetMasterDB(context.Background()).Find(&updated)
		data,err:=json.Marshal(updated)
		if err==nil{
			redisclient.Client.Set(context.Background(),"All_Inventory",string(data),1*time.Hour)
		}
	}()

	c.JSON(200,gin.H{
		"message":"Inventory upserted successfully",
		"item":inv,
	})
}
