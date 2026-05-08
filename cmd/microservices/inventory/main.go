package main

import (
	"reserach/pkg/database"
	"reserach/pkg/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func main() {
	db := database.Connect()
	r := gin.Default()

	r.POST("/products/:id/deduct", func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Amount int `json:"amount"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		tx := db.Begin()
		var product model.Product

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, id).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if product.Stock < req.Amount {
			tx.Rollback()
			c.JSON(409, gin.H{"error": "Out of stock"})
			return
		}

		product.Stock -= req.Amount
		tx.Save(&product)
		tx.Commit()

		c.JSON(200, product)
	})

	r.Run(":8081")
}
