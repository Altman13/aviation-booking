package controllers

import (
	"aviation-booking/config"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse структура для health check ответа
type HealthResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Data    HealthData `json:"data"`
}

type HealthData struct {
	Service   string         `json:"service"`
	Version   string         `json:"version"`
	Database  DatabaseStatus `json:"database"`
	Timestamp TimeStamps     `json:"timestamp"`
}

type DatabaseStatus struct {
	Status      string `json:"status"`
	Driver      string `json:"driver,omitempty"`
	Version     string `json:"version,omitempty"`
	Error       string `json:"error,omitempty"`
	TablesCount int    `json:"tables_count,omitempty"`
}

type TimeStamps struct {
	Server string `json:"server"`
	Unix   int64  `json:"unix"`
}

// DatabaseInfo содержит информацию о БД
type DatabaseInfo struct {
	Version     string
	TablesCount int
}

// HealthCheck проверяет статус сервера и базы данных с GORM
func HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:  "success",
		Message: "Service is healthy",
		Data: HealthData{
			Service: "Aviation Booking API",
			Version: "1.0.0",
			Timestamp: TimeStamps{
				Server: time.Now().Format("2006-01-02T15:04:05Z07:00"),
				Unix:   time.Now().Unix(),
			},
		},
	}

	// Проверяем базу данных с GORM
	if config.DB == nil {
		response.Status = "error"
		response.Message = "Database not initialized"
		response.Data.Database = DatabaseStatus{
			Status: "disconnected",
			Error:  "GORM DB instance is nil",
		}
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Получаем underlying sql.DB для проверки соединения
	sqlDB, err := config.DB.DB()
	if err != nil {
		response.Status = "error"
		response.Message = "Failed to get database connection"
		response.Data.Database = DatabaseStatus{
			Status: "disconnected",
			Error:  err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Пинг базы данных через sql.DB
	err = sqlDB.Ping()
	if err != nil {
		response.Status = "error"
		response.Message = "Database connection failed"
		response.Data.Database = DatabaseStatus{
			Status: "disconnected",
			Error:  err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Получаем информацию о базе данных через GORM
	dbInfo := getDatabaseInfo()

	response.Data.Database = DatabaseStatus{
		Status:      "connected",
		Driver:      "MySQL",
		Version:     dbInfo.Version,
		TablesCount: dbInfo.TablesCount,
	}

	c.JSON(http.StatusOK, response)
}

// getDatabaseInfo получает информацию о базе данных
func getDatabaseInfo() DatabaseInfo {
	info := DatabaseInfo{
		Version:     "unknown",
		TablesCount: 0,
	}

	// Проверяем DB
	if config.DB == nil {
		return info
	}

	// Получаем версию MySQL
	var version string
	sqlDB, err := config.DB.DB()
	if err == nil && sqlDB != nil {
		// Используем контекст с таймаутом
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err = sqlDB.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version)
		if err == nil {
			info.Version = version
		}
	}

	// Считаем количество таблиц
	var count int64
	config.DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE()").Scan(&count)
	info.TablesCount = int(count)

	return info
}

// SimpleHealth упрощенная версия health check
func SimpleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "aviation-booking",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// DatabaseStats возвращает статистику базы данных с GORM
func DatabaseStats(c *gin.Context) {
	// ПРОВЕРКА: Добавлена проверка на nil
	if config.DB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Database not initialized",
			"data":    nil,
		})
		return
	}

	// Проверка соединения с БД
	sqlDB, err := config.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Failed to get database connection: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Пинг базы данных
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Database ping failed: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Создаем контекст с таймаутом для запросов
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	stats := gin.H{
		"database": gin.H{
			"status": "connected",
		},
		"tables": gin.H{},
	}

	// Получаем список таблиц и их статистику через GORM
	tables := []string{"users", "airports", "airlines", "searches", "flights", "price_history", "user_favorites"}

	for _, table := range tables {
		var count int64
		// Используем контекст с таймаутом
		result := config.DB.WithContext(ctx).Table(table).Count(&count)
		if result.Error != nil {
			stats["tables"].(gin.H)[table] = gin.H{
				"count": 0,
				"error": result.Error.Error(),
			}
		} else {
			stats["tables"].(gin.H)[table] = gin.H{
				"count": count,
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   stats,
	})
}

// ReadyCheck для Kubernetes readiness probe
func ReadyCheck(c *gin.Context) {
	// Проверяем критичные зависимости
	if config.DB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not_ready",
			"message": "Database not connected",
		})
		return
	}

	// Получаем underlying sql.DB для проверки
	sqlDB, err := config.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not_ready",
			"message": "Failed to get database connection: " + err.Error(),
		})
		return
	}

	// Пинг базы данных
	err = sqlDB.Ping()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not_ready",
			"message": "Database ping failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"message": "All dependencies are healthy",
	})
}

// LiveCheck для Kubernetes liveness probe
func LiveCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"message":   "Service is running",
		"timestamp": time.Now().Unix(),
	})
}

// checkDatabaseConnection общая функция для проверки состояния подключения к БД
func checkDatabaseConnection() (bool, string) {
	if config.DB == nil {
		return false, "Database not initialized"
	}

	sqlDB, err := config.DB.DB()
	if err != nil {
		return false, "Failed to get database connection: " + err.Error()
	}

	if err := sqlDB.Ping(); err != nil {
		return false, "Database ping failed: " + err.Error()
	}

	return true, ""
}
