// controllers/user_controller.go
package controllers

import (
	"aviation-booking/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

// NewUserController создает новый контроллер пользователей
func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// Register - регистрация нового пользователя
func (uc *UserController) Register(c *gin.Context) {
	var req services.RegisterUserRequest

	// Привязываем JSON из тела запроса к структуре
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Регистрируем пользователя
	user, err := uc.userService.RegisterUser(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User registered successfully",
		"data":    user,
	})
}

// Login - аутентификация пользователя
func (uc *UserController) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Привязываем JSON из тела запроса к структуре
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Аутентифицируем пользователя
	user, err := uc.userService.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login successful",
		"data":    user,
	})
}

// GetUserByID - получение пользователя по ID
func (uc *UserController) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	// Конвертируем string в uint
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Получаем пользователя
	user, err := uc.userService.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}

// UpdateUser - обновление данных пользователя
func (uc *UserController) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	// Конвертируем string в uint
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Обновляем пользователя
	user, err := uc.userService.UpdateUser(uint(userID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User updated successfully",
		"data":    user,
	})
}

// UpdatePassword - обновление пароля пользователя
func (uc *UserController) UpdatePassword(c *gin.Context) {
	userIDStr := c.Param("id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	// Конвертируем string в uint
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Обновляем пароль
	err = uc.userService.UpdatePassword(uint(userID), req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password updated successfully",
	})
}

// DeleteUser - удаление пользователя
func (uc *UserController) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	// Конвертируем string в uint
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Удаляем пользователя
	err = uc.userService.DeleteUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User deleted successfully",
	})
}

// GetAllUsers - получение всех пользователей
func (uc *UserController) GetAllUsers(c *gin.Context) {
	// Получаем всех пользователей
	users, err := uc.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"users": users,
			"total": len(users),
		},
	})
}

// SearchUsers - поиск пользователей
func (uc *UserController) SearchUsers(c *gin.Context) {
	// Парсим параметры запроса
	limitStr := c.DefaultQuery("limit", "50")
	pageStr := c.DefaultQuery("page", "1")

	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)

	params := services.UserSearchParams{
		Email:       c.Query("email"),
		FirstName:   c.Query("first_name"),
		LastName:    c.Query("last_name"),
		PhoneNumber: c.Query("phone_number"),
		SortBy:      c.Query("sort_by"),
		Limit:       limit,
		Offset:      (page - 1) * limit,
	}

	// Ищем пользователей
	users, err := uc.userService.SearchUsers(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"users": users,
			"total": len(users),
			"page":  page,
			"limit": limit,
		},
	})
}
