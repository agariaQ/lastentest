package main

import (
	"reserach/pkg/database"
	"reserach/pkg/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func main() {

	db := database.Connect()
	database.SeedData(db)

	r := gin.Default()

	r.POST("/orders", func(c *gin.Context) {
		var order model.Order
		if err := c.BindJSON(&order); err != nil {
			c.JSON(400, gin.H{"error": "Ungültiges JSON Format: " + err.Error()})
			return
		}

		tx := db.Begin()

		var product model.Product

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, order.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if product.Stock < order.Amount {
			tx.Rollback()
			c.JSON(409, gin.H{"error": "Out of stock"})
			return
		}

		product.Stock -= order.Amount
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		tx.Commit()
		c.JSON(201, order)
	})

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
