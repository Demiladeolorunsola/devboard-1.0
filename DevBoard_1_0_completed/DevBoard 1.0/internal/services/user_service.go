package services

import (
	"errors"
	"strings"

	"github.com/Demiladeolorunsola/devboard-1.0/internal/models"
	"github.com/Demiladeolorunsola/devboard-1.0/internal/utils"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	var user models.User

	result := s.DB.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (s *UserService) Update(
	id uint,
	name string,
	email string,
) (*models.User, error) {

	user, err := s.GetByID(id)

	if err != nil {
		return nil, err
	}

	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name != "" {
		user.Name = name
	}

	if email != "" && email != user.Email {
		var existing models.User

		result := s.DB.
			Where("email = ? AND id != ?", email, id).
			First(&existing)

		if result.Error == nil {
			return nil, errors.New("email already in use")
		}

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}

		user.Email = email
	}

	if err := s.DB.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(id uint) error {
	result := s.DB.Delete(&models.User{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (s *UserService) UpdatePassword(
	id uint,
	currentPassword string,
	newPassword string,
) error {

	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}

	user, err := s.GetByID(id)

	if err != nil {
		return err
	}

	if err := utils.CheckPassword(
		currentPassword,
		user.Password,
	); err != nil {
		return errors.New("current password is incorrect")
	}

	hash, err := utils.HashPassword(newPassword)

	if err != nil {
		return err
	}

	user.Password = hash

	return s.DB.Save(user).Error
}
