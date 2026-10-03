package services

import (
	"errors"
	"strings"

	"github.com/Demiladeolorunsola/devboard-1.0/internal/models"
	"github.com/Demiladeolorunsola/devboard-1.0/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

type LoginResult struct {
	User  *models.User
	Token string
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		DB: db,
	}
}

func (s *AuthService) Register(
	name string,
	email string,
	password string,
) (*models.User, error) {

	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" {
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	var existingUser models.User

	result := s.DB.
		Where("email = ?", email).
		First(&existingUser)

	if result.Error == nil {
		return nil, errors.New("email already registered")
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Role:     "member",
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) Login(
	email string,
	password string,
	jwtSecret string,
) (*LoginResult, error) {

	email = strings.ToLower(strings.TrimSpace(email))

	var user models.User

	result := s.DB.
		Where("email = ?", email).
		First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("invalid email or password")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	if err := utils.CheckPassword(password, user.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, jwtSecret)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:  &user,
		Token: token,
	}, nil
}
