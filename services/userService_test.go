// services/user_service_test.go
package services

import (
	"aviation-booking/config"
	"aviation-booking/models"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceTestSuite struct {
	suite.Suite
	db          *gorm.DB
	userService *UserService
}

func (suite *UserServiceTestSuite) SetupTest() {
	// Используем тестовую базу данных
	config.InitTestDB()
	suite.db = config.TestDB
	suite.userService = NewUserService().WithDB(suite.db)

	// Очищаем таблицы перед каждым тестом
	suite.db.Exec("DELETE FROM users")
}

func (suite *UserServiceTestSuite) TearDownTest() {
	config.CloseTestDB()
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) TestRegisterUser_Success() {
	req := RegisterUserRequest{
		Email:       "test@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "+1234567890",
		DateOfBirth: "1990-01-01",
	}

	user, err := suite.userService.RegisterUser(&req)

	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal("test@example.com", user.Email)
	suite.Equal("John", user.FirstName)
	suite.Equal("Doe", user.LastName)
	suite.Equal("+1234567890", user.PhoneNumber)
	suite.Equal("", user.PasswordHash) // Не возвращаем хеш пароля

	// Проверяем, что дата рождения правильно сохранена
	expectedDOB, _ := time.Parse("2006-01-02", "1990-01-01")
	suite.Equal(expectedDOB.Format("2006-01-02"), user.DateOfBirth.Format("2006-01-02"))

	// Проверяем, что пользователь сохранен в базе
	var dbUser models.User
	suite.db.First(&dbUser, user.ID)
	suite.Equal("test@example.com", dbUser.Email)
	suite.NotEmpty(dbUser.PasswordHash) // Пароль должен быть захэширован

	// Проверяем, что пароль правильно захэширован
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte("password123"))
	suite.NoError(err)
}

func (suite *UserServiceTestSuite) TestRegisterUser_DuplicateEmail() {
	// Создаем первого пользователя
	req1 := RegisterUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	user1, err := suite.userService.RegisterUser(&req1)
	suite.NoError(err)
	suite.NotNil(user1)

	// Пытаемся создать второго пользователя с тем же email
	req2 := RegisterUserRequest{
		Email:     "test@example.com",
		Password:  "password456",
		FirstName: "Jane",
		LastName:  "Smith",
	}
	user2, err := suite.userService.RegisterUser(&req2)

	suite.Error(err)
	suite.Nil(user2)
	suite.Equal("пользователь с таким email уже существует", err.Error())
}

func (suite *UserServiceTestSuite) TestRegisterUser_ValidationErrors() {
	tests := []struct {
		name     string
		request  RegisterUserRequest
		expected string
	}{
		{
			name: "Пустой email",
			request: RegisterUserRequest{
				Email:     "",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "email обязателен",
		},
		{
			name: "Неверный формат email",
			request: RegisterUserRequest{
				Email:     "invalid-email",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "неверный формат email",
		},
		{
			name: "Пустое имя",
			request: RegisterUserRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "",
				LastName:  "Doe",
			},
			expected: "имя обязательно",
		},
		{
			name: "Пустая фамилия",
			request: RegisterUserRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "",
			},
			expected: "фамилия обязательна",
		},
		{
			name: "Неверный формат телефона",
			request: RegisterUserRequest{
				Email:       "test@example.com",
				Password:    "password123",
				FirstName:   "John",
				LastName:    "Doe",
				PhoneNumber: "123", // Слишком короткий
			},
			expected: "неверный формат телефона",
		},
		{
			name: "Неверный формат даты рождения",
			request: RegisterUserRequest{
				Email:       "test@example.com",
				Password:    "password123",
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: "invalid-date",
			},
			expected: "неверный формат даты рождения",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			user, err := suite.userService.RegisterUser(&tt.request)
			suite.Error(err)
			suite.Nil(user)
			suite.Contains(err.Error(), tt.expected)
		})
	}
}

func (suite *UserServiceTestSuite) TestLoginUser_Success() {
	// Сначала регистрируем пользователя
	req := RegisterUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	registeredUser, err := suite.userService.RegisterUser(&req)
	suite.NoError(err)

	// Тестируем логин
	user, err := suite.userService.LoginUser("test@example.com", "password123")

	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal(registeredUser.ID, user.ID)
	suite.Equal("test@example.com", user.Email)
	suite.Equal("John", user.FirstName)
	suite.Equal("Doe", user.LastName)
	suite.Equal("", user.PasswordHash) // Не возвращаем хеш пароля
}

func (suite *UserServiceTestSuite) TestLoginUser_WrongPassword() {
	// Регистрируем пользователя
	req := RegisterUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	_, err := suite.userService.RegisterUser(&req)
	suite.NoError(err)

	// Пытаемся войти с неправильным паролем
	user, err := suite.userService.LoginUser("test@example.com", "wrongpassword")

	suite.Error(err)
	suite.Nil(user)
	suite.Equal("неверный email или пароль", err.Error())
}

func (suite *UserServiceTestSuite) TestLoginUser_UserNotFound() {
	// Пытаемся войти с несуществующим email
	user, err := suite.userService.LoginUser("nonexistent@example.com", "password123")

	suite.Error(err)
	suite.Nil(user)
	suite.Equal("неверный email или пароль", err.Error())
}

func (suite *UserServiceTestSuite) TestGetUserByID_Success() {
	// Создаем пользователя
	req := RegisterUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	registeredUser, err := suite.userService.RegisterUser(&req)
	suite.NoError(err)

	// Получаем пользователя по ID
	user, err := suite.userService.GetUserByID(registeredUser.ID)

	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal(registeredUser.ID, user.ID)
	suite.Equal("test@example.com", user.Email)
	suite.Equal("John", user.FirstName)
	suite.Equal("Doe", user.LastName)
	suite.Equal("", user.PasswordHash)
}

func (suite *UserServiceTestSuite) TestGetUserByID_NotFound() {
	// Пытаемся получить несуществующего пользователя
	user, err := suite.userService.GetUserByID(999)

	suite.Error(err)
	suite.Nil(user)
	suite.Equal("пользователь не найден", err.Error())
}

func (suite *UserServiceTestSuite) TestGetUserByEmail_Success() {
	// Создаем пользователя
	req := RegisterUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	registeredUser, err := suite.userService.RegisterUser(&req)
	suite.NoError(err)

	// Получаем пользователя по email
	user, err := suite.userService.GetUserByEmail("test@example.com")

	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal(registeredUser.ID, user.ID)
	suite.Equal("test@example.com", user.Email)
	suite.Equal("John", user.FirstName)
	suite.Equal("Doe", user.LastName)
	suite.Equal("", user.PasswordHash)
}

func (suite *UserServiceTestSuite) TestGetUserByEmail_NotFound() {
	// Пытаемся получить несуществующего пользователя
	user, err := suite.userService.GetUserByEmail("nonexistent@example.com")

	suite.Error(err)
	suite.Nil(user)
	suite.Equal("пользователь не найден", err.Error())
}

func (suite *UserServiceTestSuite) TestUpdateUser_Success() {
	// Создаем пользователя
	req := RegisterUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	registeredUser, err := suite.userService.RegisterUser(&req)
	suite.NoError(err)

	// Обновляем пользователя
	updateReq := UpdateUserRequest{
		FirstName:   "Jonathan",
		LastName:    "Smith",
		PhoneNumber: "+1234567890",
		DateOfBirth: "1985-05-15",
	}

	updatedUser, err := suite.userService.UpdateUser(registeredUser.ID, &updateReq)

	suite.NoError(err)
	suite.NotNil(updatedUser)
	suite.Equal(registeredUser.ID, updatedUser.ID)
	suite.Equal("test@example.com", updatedUser.Email) // Email не меняется
	suite.Equal("Jonathan", updatedUser.FirstName)
	suite.Equal("Smith", updatedUser.LastName)
	suite.Equal("+1234567890", updatedUser.PhoneNumber)

	// Проверяем дату рождения
	expectedDOB, _ := time.Parse("2006-01-02", "1985-05-15")
	suite.Equal(expectedDOB.Format("2006-01-02"), updatedUser.DateOfBirth.Format("2006-01-02"))
}

func (suite *UserServiceTestSuite) TestUpdateUser_PartialUpdate() {
	// Создаем пользователя
	req := RegisterUserRequest{
		Email:       "test@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "+1111111111",
		DateOfBirth: "1990-01-01",
	}
	registeredUser, err := suite.userService.RegisterUser(&req)
	suite.NoError(err)

	// Обновляем только имя
	updateReq := UpdateUserRequest{
		FirstName: "Jonathan",
	}

	updatedUser, err := suite.userService.UpdateUser(registeredUser.ID, &updateReq)

	suite.NoError(err)
	suite.NotNil(updatedUser)
	suite.Equal("Jonathan", updatedUser.FirstName)      // Новое имя
	suite.Equal("Doe", updatedUser.LastName)            // Фамилия не изменилась
	suite.Equal("+1111111111", updatedUser.PhoneNumber) // Телефон не изменился

	// Проверяем дату рождения
	expectedDOB, _ := time.Parse("2006-01-02", "1990-01-01")
	suite.Equal(expectedDOB.Format("2006-01-02"), updatedUser.DateOfBirth.Format("2006-01-02"))
}

func (suite *UserServiceTestSuite) TestUpdateUser_NotFound() {
	updateReq := UpdateUserRequest{
		FirstName: "John",
	}

	user, err := suite.userService.UpdateUser(999, &updateReq)

	suite.Error(err)
	suite.Nil(user)
	suite.Equal("пользователь не найден", err.Error())
}

func (suite *UserServiceTestSuite) TestUpdateUser_ValidationError() {
	// Создаем пользователя
	req := RegisterUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	registeredUser, err := suite.userService.RegisterUser(&req)
	suite.NoError(err)

	// Пытаемся обновить с неверным телефоном
	updateReq := UpdateUserRequest{
		PhoneNumber: "123", // Неверный формат
	}

	user, err := suite.userService.UpdateUser(registeredUser.ID, &updateReq)

	suite.Error(err)
	suite.Nil(user)
	suite.Equal("неверный формат телефона", err.Error())
}

func (suite *UserServiceTestSuite) TestUpdatePassword_Success() {
	// Создаем пользователя
	req := RegisterUserRequest{
		Email:     "test@example.com",
		Password:  "oldpassword",
		FirstName: "John",
		LastName:  "Doe",
	}
	registeredUser, err := suite.userService.RegisterUser(&req)
	suite.NoError(err)

	// Обновляем пароль
	err = suite.userService.UpdatePassword(registeredUser.ID, "oldpassword", "newpassword123")
	suite.NoError(err)

	// Проверяем, что можно войти с новым паролем
	user, err := suite.userService.LoginUser("test@example.com", "newpassword123")
	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal(registeredUser.ID, user.ID)
}

func (suite *UserServiceTestSuite) TestUpdatePassword_WrongOldPassword() {
	// Создаем пользователя
	req := RegisterUserRequest{
		Email:     "test@example.com",
		Password:  "oldpassword",
		FirstName: "John",
		LastName:  "Doe",
	}
	registeredUser, err := suite.userService.RegisterUser(&req)
	suite.NoError(err)

	// Пытаемся обновить пароль с неправильным старым паролем
	err = suite.userService.UpdatePassword(registeredUser.ID, "wrongpassword", "newpassword123")

	suite.Error(err)
	suite.Equal("неверный текущий пароль", err.Error())
}

func (suite *UserServiceTestSuite) TestUpdatePassword_UserNotFound() {
	// Пытаемся обновить пароль несуществующего пользователя
	err := suite.userService.UpdatePassword(999, "oldpassword", "newpassword123")

	suite.Error(err)
	suite.Equal("пользователь не найден", err.Error())
}

func (suite *UserServiceTestSuite) TestDeleteUser_Success() {
	// Создаем пользователя
	req := RegisterUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	registeredUser, err := suite.userService.RegisterUser(&req)
	suite.NoError(err)

	// Удаляем пользователя
	err = suite.userService.DeleteUser(registeredUser.ID)
	suite.NoError(err)

	// Проверяем, что пользователь удален
	user, err := suite.userService.GetUserByID(registeredUser.ID)
	suite.Error(err)
	suite.Nil(user)
}

func (suite *UserServiceTestSuite) TestDeleteUser_NotFound() {
	// Пытаемся удалить несуществующего пользователя
	err := suite.userService.DeleteUser(999)

	suite.Error(err)
	suite.Equal("пользователь не найден", err.Error())
}

func (suite *UserServiceTestSuite) TestGetAllUsers() {
	// Создаем нескольких пользователей
	usersToCreate := []RegisterUserRequest{
		{
			Email:     "user1@example.com",
			Password:  "password1",
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			Email:     "user2@example.com",
			Password:  "password2",
			FirstName: "Jane",
			LastName:  "Smith",
		},
		{
			Email:     "user3@example.com",
			Password:  "password3",
			FirstName: "Bob",
			LastName:  "Johnson",
		},
	}

	for _, req := range usersToCreate {
		_, err := suite.userService.RegisterUser(&req)
		suite.NoError(err)
	}

	// Получаем всех пользователей
	users, err := suite.userService.GetAllUsers()

	suite.NoError(err)
	suite.Len(users, 3)

	// Проверяем, что хеши паролей не возвращаются
	for _, user := range users {
		suite.Equal("", user.PasswordHash)
	}

	// Проверяем, что все пользователи присутствуют
	emails := make(map[string]bool)
	for _, user := range users {
		emails[user.Email] = true
	}
	suite.True(emails["user1@example.com"])
	suite.True(emails["user2@example.com"])
	suite.True(emails["user3@example.com"])
}

func (suite *UserServiceTestSuite) TestSearchUsers() {
	// Создаем тестовых пользователей
	usersToCreate := []RegisterUserRequest{
		{
			Email:       "john.doe@example.com",
			Password:    "password1",
			FirstName:   "John",
			LastName:    "Doe",
			PhoneNumber: "+1234567890",
		},
		{
			Email:       "jane.doe@example.com",
			Password:    "password2",
			FirstName:   "Jane",
			LastName:    "Doe",
			PhoneNumber: "+0987654321",
		},
		{
			Email:       "bob.smith@example.com",
			Password:    "password3",
			FirstName:   "Bob",
			LastName:    "Smith",
			PhoneNumber: "+1111111111",
		},
		{
			Email:       "alice.johnson@example.com",
			Password:    "password4",
			FirstName:   "Alice",
			LastName:    "Johnson",
			PhoneNumber: "+2222222222",
		},
	}

	for _, req := range usersToCreate {
		_, err := suite.userService.RegisterUser(&req)
		suite.NoError(err)
	}

	// Тест 1: Поиск по фамилии
	params := UserSearchParams{
		LastName: "Doe",
	}
	users, err := suite.userService.SearchUsers(params)
	suite.NoError(err)
	suite.Len(users, 2) // John Doe и Jane Doe

	// Тест 2: Поиск по имени
	params = UserSearchParams{
		FirstName: "Bob",
	}
	users, err = suite.userService.SearchUsers(params)
	suite.NoError(err)
	suite.Len(users, 1)
	suite.Equal("Bob", users[0].FirstName)

	// Тест 3: Поиск по email (частичный)
	params = UserSearchParams{
		Email: "doe",
	}
	users, err = suite.userService.SearchUsers(params)
	suite.NoError(err)
	suite.Len(users, 2) // Оба пользователя с фамилией Doe

	// Тест 4: Поиск с пагинацией
	params = UserSearchParams{
		Limit:  2,
		Offset: 0,
		SortBy: "first_name ASC",
	}
	users, err = suite.userService.SearchUsers(params)
	suite.NoError(err)
	suite.Len(users, 2)
	// Проверяем сортировку по имени
	if len(users) >= 2 {
		// Alice должна быть первой при сортировке ASC
		suite.Equal("Alice", users[0].FirstName)
		suite.Equal("Bob", users[1].FirstName)
	}

	// Тест 5: Поиск по телефону
	params = UserSearchParams{
		PhoneNumber: "1111111111",
	}
	users, err = suite.userService.SearchUsers(params)
	suite.NoError(err)
	suite.Len(users, 1)
	suite.Equal("Bob", users[0].FirstName)

	// Тест 6: Пустой поиск (все пользователи)
	params = UserSearchParams{}
	users, err = suite.userService.SearchUsers(params)
	suite.NoError(err)
	suite.Len(users, 4)
}

func (suite *UserServiceTestSuite) TestValidateUser() {
	tests := []struct {
		name     string
		user     models.User
		expected string
	}{
		{
			name: "Валидный пользователь",
			user: models.User{
				Email:       "test@example.com",
				FirstName:   "John",
				LastName:    "Doe",
				PhoneNumber: "+1234567890",
			},
			expected: "",
		},
		{
			name: "Пустой email",
			user: models.User{
				Email:     "",
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "email обязателен",
		},
		{
			name: "Неверный формат email",
			user: models.User{
				Email:     "invalid-email",
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "неверный формат email",
		},
		{
			name: "Пустое имя",
			user: models.User{
				Email:     "test@example.com",
				FirstName: "",
				LastName:  "Doe",
			},
			expected: "имя обязательно",
		},
		{
			name: "Пустая фамилия",
			user: models.User{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "",
			},
			expected: "фамилия обязательна",
		},
		{
			name: "Неверный формат телефона",
			user: models.User{
				Email:       "test@example.com",
				FirstName:   "John",
				LastName:    "Doe",
				PhoneNumber: "123",
			},
			expected: "неверный формат телефона",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := suite.userService.validateUser(&tt.user)
			if tt.expected == "" {
				suite.NoError(err)
			} else {
				suite.Error(err)
				suite.Contains(err.Error(), tt.expected)
			}
		})
	}
}
