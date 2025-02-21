package controllers

import (
	"fmt"
	"net/http"
	"taskmanagement/config"
	models "taskmanagement/model"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User

	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	userResponses := make([]gin.H, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, gin.H{
			"id":   user.ID,
			"name": user.Username,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": userResponses})
}

func GetUser(c *gin.Context) {
	var user models.User
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	fmt.Println("User ID", userId)

	if err := config.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var cutomerResponse = gin.H{
		"id":    user.ID,
		"name":  user.Username,
		"email": user.Email,
	}

	c.JSON(http.StatusOK, gin.H{"user": cutomerResponse})
}
