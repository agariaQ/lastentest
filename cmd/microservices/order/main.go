package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reserach/pkg/database"
	"reserach/pkg/model"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.Connect()
	r := gin.Default()

	client := &http.Client{}

	r.POST("/orders", func(c *gin.Context) {
		var order model.Order
		if err := c.BindJSON(&order); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		inventoryPayload := map[string]int{"amount": order.Amount}
		jsonBody, _ := json.Marshal(inventoryPayload)

		inventoryURL := fmt.Sprintf("http://localhost:8081/products/%d/deduct", order.ProductID)
		resp, err := client.Post(inventoryURL, "application/json", bytes.NewBuffer(jsonBody))

		if err != nil {
			c.JSON(503, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			c.JSON(409, gin.H{"error": "Out of stock (rejected by inventory)"})
			return
		}

		if err := db.Create(&order).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, order)
	})
	r.Run(":8082")
}
