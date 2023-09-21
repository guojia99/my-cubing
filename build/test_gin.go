package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Address struct {
	City string `json:"city" binding:"max=16,required"`
}

type User struct {
	Username string  `json:"username" binding:"len=8,required"`
	Email    string  `json:"email" binding:"required,email"`
	Address  Address `json:"address" binding:"required"`
}

func main() {
	router := gin.Default()

	router.POST("/", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
	})

	router.Run(":8080")
}
