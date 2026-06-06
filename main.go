package main

import (
	"cafe-app-backend/database"
	"cafe-app-backend/model/owner"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {

	db := database.InitializeDB()
	defer db.Close()

	// Initialize the Gin
	router := gin.Default()

	// Get the cost of an item by its code
	// Example URL: https://cafe-app-test.onrender.com/cost/C01
	router.GET("/cost/:code", func(c *gin.Context) {
		code := c.Param("code")

		cost, err := owner.GetCost(db, code)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"cost": cost,
		})
	})

	// Update the price of an item
	// JSON from the frontend {"code": "C01", "price": 400}
	router.POST("/update-cost", func(c *gin.Context) {
		var input struct {
			Code  string `json:"code" binding:"required"`
			Price int    `json:"price" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		owner.UpdateCost(db, input.Code, input.Price)

		c.JSON(http.StatusOK, gin.H{
			"message": "Price updated successfully to " + strconv.Itoa(input.Price) + " yen",
		})
	})

	// Simple health check to verify the app is live
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Getting the port that can be assigned dynamically in Render or defaulting to local
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server (This loops the main function! Needed for continuous running)
	router.Run(":" + port)
}
