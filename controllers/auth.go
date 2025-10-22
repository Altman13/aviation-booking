package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRequest - структура для запроса регистрации
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginRequest - структура для запроса входа
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterUser - регистрация нового пользователя
func RegisterUser(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid registration data",
			"details": err.Error(),
		})
		return
	}

	// Здесь будет логика регистрации
	// Пока возвращаем заглушку

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User registered successfully",
		"data": gin.H{
			"user_id": "user-" + req.Email, // Временный ID
			"email":   req.Email,
			"name":    req.FirstName + " " + req.LastName,
		},
	})
}

// LoginUser - аутентификация пользователя
func LoginUser(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid login data",
			"details": err.Error(),
		})
		return
	}

	// Здесь будет логика аутентификации
	// Пока возвращаем заглушку

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"token": "jwt-token-example",
			"user": gin.H{
				"id":    "user-" + req.Email,
				"email": req.Email,
				"name":  "John Doe", // Временные данные
			},
		},
	})
}
