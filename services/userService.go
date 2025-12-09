// services/user_service.go
package services

import (
	"aviation-booking/config"
	"aviation-booking/models"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

// NewUserService создает новый экземпляр сервиса пользователей
func NewUserService() *UserService {
	return &UserService{db: config.DB}
}

// WithDB позволяет передать кастомное соединение с БД
func (s *UserService) WithDB(db *gorm.DB) *UserService {
	s.db = db
	return s
}

// RegisterUser регистрирует нового пользователя
func (s *UserService) RegisterUser(userData *RegisterUserRequest) (*models.User, error) {
	// Проверяем существование пользователя с таким email
	var existingUser models.User
	if err := s.db.Where("email = ?", userData.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("ошибка при хешировании пароля")
	}

	// Парсим дату рождения
	var dateOfBirth time.Time
	if userData.DateOfBirth != "" {
		dateOfBirth, err = time.Parse("2006-01-02", userData.DateOfBirth)
		if err != nil {
			return nil, errors.New("неверный формат даты рождения")
		}
	}

	// Создаем пользователя
	user := &models.User{
		Email:        userData.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    userData.FirstName,
		LastName:     userData.LastName,
		PhoneNumber:  userData.PhoneNumber,
		DateOfBirth:  dateOfBirth,
	}

	// Валидация
	if err := s.validateUser(user); err != nil {
		return nil, err
	}

	// Сохраняем в БД
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	// Не возвращаем хеш пароля
	user.PasswordHash = ""
	return user, nil
}

// LoginUser аутентифицирует пользователя
func (s *UserService) LoginUser(email, password string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("неверный email или пароль")
		}
		return nil, err
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	// Не возвращаем хеш пароля
	user.PasswordHash = ""
	return &user, nil
}

// GetUserByID получает пользователя по ID
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}

	// Не возвращаем хеш пароля
	user.PasswordHash = ""
	return &user, nil
}

// GetUserByEmail получает пользователя по email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}

	// Не возвращаем хеш пароля
	user.PasswordHash = ""
	return &user, nil
}

// UpdateUser обновляет данные пользователя
func (s *UserService) UpdateUser(userID uint, userData *UpdateUserRequest) (*models.User, error) {
	// Получаем существующего пользователя
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Обновляем поля
	if userData.FirstName != "" {
		user.FirstName = userData.FirstName
	}
	if userData.LastName != "" {
		user.LastName = userData.LastName
	}
	if userData.PhoneNumber != "" {
		user.PhoneNumber = userData.PhoneNumber
	}
	if userData.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", userData.DateOfBirth)
		if err != nil {
			return nil, errors.New("неверный формат даты рождения")
		}
		user.DateOfBirth = dob
	}

	// Валидация
	if err := s.validateUser(user); err != nil {
		return nil, err
	}

	// Сохраняем изменения
	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// UpdatePassword обновляет пароль пользователя
func (s *UserService) UpdatePassword(userID uint, oldPassword, newPassword string) error {
	// Получаем пользователя с паролем
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("пользователь не найден")
		}
		return err
	}

	// Проверяем старый пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("неверный текущий пароль")
	}

	// Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("ошибка при хешировании пароля")
	}

	// Обновляем пароль
	user.PasswordHash = string(hashedPassword)
	return s.db.Save(&user).Error
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(userID uint) error {
	// Проверяем существование пользователя
	if _, err := s.GetUserByID(userID); err != nil {
		return err
	}

	// Удаляем пользователя
	result := s.db.Delete(&models.User{}, userID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("пользователь не найден")
	}
	return nil
}

// GetAllUsers получает всех пользователей
func (s *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := s.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	// Не возвращаем хеши паролей
	for i := range users {
		users[i].PasswordHash = ""
	}
	return users, nil
}

// SearchUsers ищет пользователей по параметрам
func (s *UserService) SearchUsers(params UserSearchParams) ([]models.User, error) {
	var users []models.User

	query := s.db.Model(&models.User{})

	// Применяем фильтры
	if params.Email != "" {
		query = query.Where("email LIKE ?", "%"+params.Email+"%")
	}
	if params.FirstName != "" {
		query = query.Where("first_name LIKE ?", "%"+params.FirstName+"%")
	}
	if params.LastName != "" {
		query = query.Where("last_name LIKE ?", "%"+params.LastName+"%")
	}
	if params.PhoneNumber != "" {
		query = query.Where("phone_number LIKE ?", "%"+params.PhoneNumber+"%")
	}

	// Пагинация
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	// Сортировка
	if params.SortBy != "" {
		query = query.Order(params.SortBy)
	} else {
		query = query.Order("created_at DESC")
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	// Не возвращаем хеши паролей
	for i := range users {
		users[i].PasswordHash = ""
	}
	return users, nil
}

// validateUser - валидация данных пользователя
func (s *UserService) validateUser(user *models.User) error {
	if user.Email == "" {
		return errors.New("email обязателен")
	}
	// Простая проверка формата email
	if !isValidEmail(user.Email) {
		return errors.New("неверный формат email")
	}
	if user.FirstName == "" {
		return errors.New("имя обязательно")
	}
	if user.LastName == "" {
		return errors.New("фамилия обязательна")
	}
	if user.PhoneNumber != "" && !isValidPhone(user.PhoneNumber) {
		return errors.New("неверный формат телефона")
	}
	return nil
}

// isValidEmail - проверка формата email
func isValidEmail(email string) bool {
	// Простая проверка, можно заменить на регулярное выражение
	return len(email) > 3 && len(email) < 255 &&
		contains(email, "@") && contains(email, ".")
}

// isValidPhone - проверка формата телефона
func isValidPhone(phone string) bool {
	// Простая проверка, можно заменить на более сложную логику
	return len(phone) >= 10 && len(phone) <= 15
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || contains(s[1:], substr)))
}

// Структуры запросов

// RegisterUserRequest - запрос на регистрацию пользователя
type RegisterUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number"`
	DateOfBirth string `json:"date_of_birth"` // Формат: "2006-01-02"
}

// UpdateUserRequest - запрос на обновление пользователя
type UpdateUserRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	DateOfBirth string `json:"date_of_birth"` // Формат: "2006-01-02"
}

// UserSearchParams - параметры поиска пользователей
type UserSearchParams struct {
	Email       string
	FirstName   string
	LastName    string
	PhoneNumber string
	SortBy      string // например: "created_at DESC", "first_name ASC"
	Limit       int
	Offset      int
}
