package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/ducanng/no-name/internal/model/userdm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *userdm.User) (*userdm.User, error)
	CreateUserCredential(ctx context.Context, userCredential *userdm.UserCredential) (*userdm.UserCredential, error)
	GetExistingUserByID(ctx context.Context, userID uint64) (user *userdm.User, isFound bool, err error)
	GetExistingUserByEmail(ctx context.Context, email string) (*userdm.User, bool, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (u *userRepository) CreateUser(ctx context.Context, user *userdm.User) (*userdm.User, error) {
	err := u.DB.WithContext(ctx).Create(&user).Error
	return user, err
}

func (u *userRepository) CreateUserCredential(ctx context.Context, userCredential *userdm.UserCredential) (*userdm.UserCredential, error) {
	err := u.DB.WithContext(ctx).Create(&userCredential).Error
	return userCredential, err
}

func (u *userRepository) GetExistingUserByID(ctx context.Context, userID uint64) (user *userdm.User, isFound bool, err error) {
	err = u.DB.WithContext(ctx).Model(&userdm.User{}).Where("user_id = ?", userID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return user, true, nil
}

func (u *userRepository) GetExistingUserByEmail(ctx context.Context, email string) (user *userdm.User, isFound bool, err error) {
	err = u.DB.WithContext(ctx).Model(&userdm.User{}).Where("user_email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return user, true, nil
}
