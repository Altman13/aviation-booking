package main

import (
	"aviation-booking/config"
	"aviation-booking/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	config.InitDB()
	defer config.CloseDB()

	// Создание маршрутов
	r := gin.Default()

	// Создаем контроллеры
	flightController := controllers.NewFlightController()
	userController := controllers.NewUserController()
	bookingController := controllers.NewBookingController()

	// Health check endpoints
	r.GET("/health", controllers.HealthCheck)
	r.GET("/health/simple", controllers.SimpleHealth)
	r.GET("/ready", controllers.ReadyCheck)
	r.GET("/live", controllers.LiveCheck)
	r.GET("/stats/db", controllers.DatabaseStats)

	// Публичные роуты
	public := r.Group("/api")
	{
		// Аутентификация
		public.POST("/auth/register", userController.Register)
		public.POST("/auth/login", userController.Login)

		// Поиск рейсов (публичный)
		public.GET("/flights/search", flightController.SearchFlights)
		public.GET("/flights/upcoming", flightController.GetUpcomingFlights)
		public.GET("/flights/airline/:airline", flightController.GetFlightsByAirline)
		public.GET("/flights/:id", flightController.GetFlightByID)

		// Проверка доступности мест
		public.POST("/flights/:id/check-availability", flightController.CheckFlightAvailability)
	}

	// Защищенные роуты (требуют аутентификации)
	protected := r.Group("/api")
	// protected.Use(AuthMiddleware()) // Добавить middleware позже
	{
		// Бронирования
		protected.POST("/bookings", bookingController.BookFlight)
		protected.GET("/bookings/user/:user_id", bookingController.GetUserBookings)
		protected.GET("/bookings/search", bookingController.SearchBookings)
		protected.GET("/bookings/:booking_id", bookingController.GetBookingByID)
		protected.GET("/bookings/number/:booking_number", bookingController.GetBookingByNumber)
		protected.PUT("/bookings/:booking_id", bookingController.UpdateBooking)
		protected.POST("/bookings/:booking_id/cancel", bookingController.CancelBooking)
		protected.POST("/bookings/:booking_id/confirm-payment", bookingController.ConfirmPayment)

		// Профиль пользователя
		protected.GET("/users/me", userController.GetUserByID) // TODO: Нужно будет доработать
		protected.PUT("/users/me", userController.UpdateUser)
		protected.PATCH("/users/me/password", userController.UpdatePassword)
	}

	// Админские роуты (дополнительная проверка прав)
	admin := r.Group("/api/admin")
	// admin.Use(AuthMiddleware(), AdminMiddleware())
	{
		// Управление рейсами
		admin.GET("/flights", flightController.GetAllFlights)
		admin.POST("/flights", flightController.CreateFlight)
		admin.PUT("/flights/:id", flightController.UpdateFlight)
		admin.DELETE("/flights/:id", flightController.DeleteFlight)
		admin.PATCH("/flights/:id/seats", flightController.UpdateFlightSeats)

		// Управление пользователями
		admin.GET("/users", userController.GetAllUsers)
		admin.GET("/users/search", userController.SearchUsers)
		admin.GET("/users/:id", userController.GetUserByID)
		admin.PUT("/users/:id", userController.UpdateUser)
		admin.DELETE("/users/:id", userController.DeleteUser)

		// Управление бронированиями
		admin.GET("/bookings", bookingController.SearchBookings)
	}

	// Для обратной совместимости со старыми роутами
	r.POST("/book", bookingController.BookFlight)
	r.GET("/bookings/:user_id", bookingController.GetUserBookings)
	r.DELETE("/bookings/:booking_id", bookingController.CancelBooking)

	// Запуск сервера
	r.Run(":8080")
}
