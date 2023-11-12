package userservice

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/ducanng/no-name/config"
	"github.com/ducanng/no-name/internal/model/userdm"
	"github.com/ducanng/no-name/internal/repository"
	"github.com/ducanng/no-name/pkg/logutil"
)

type UserService interface {
	RegisterUser(ctx context.Context, req *UserRegistrationRequest) (*UserRegistrationResponse, error)
	UserLogin(ctx context.Context, login UserLogin) (*UserInfo, error)
}

type userService struct {
	cfg  *config.Config
	repo repository.UserRepository
}

func NewUserService(cfg *config.Config, repo repository.UserRepository) UserService {
	return &userService{
		cfg:  cfg,
		repo: repo,
	}
}

func (u *userService) RegisterUser(ctx context.Context, req *UserRegistrationRequest) (resp *UserRegistrationResponse, err error) {
	var (
		logger   = logutil.GetLogger()
		userData *userdm.User
	)
	if userData, err = u.validateUserAndPassword(ctx, req.Email, req.Password, req.RePassword); err != nil {
		logger.Errorf("error while validating user and password, err: %s", err.Error())
		return nil, err
	}
	var encryptedPassword []byte
	if encryptedPassword, err = bcrypt.GenerateFromPassword([]byte(req.Password), 11); err != nil {
		logger.Errorf("Generate password failed, err: %s", err.Error())
		return nil, err
	}

	userData, err = u.repo.CreateUser(ctx, &userdm.User{
		UserEmail:    req.Email,
		DisplayName:  req.DisplayName,
		PhoneNumber:  req.PhoneNumber,
		Password:     string(encryptedPassword),
		ActiveStatus: true,
		CreatedBy:    req.Email,
		UpdatedBy:    req.Email,
	})
	if err != nil {
		logger.Errorf("Failed to create user, err=%s", err.Error())
		return nil, err
	}
	return &UserRegistrationResponse{
		NextAction:  "",
		Session:     "",
		PhoneNumber: userData.PhoneNumber,
		Email:       userData.UserEmail,
		Message:     "",
	}, err
}

func (u *userService) validateUserAndPassword(ctx context.Context, email string, password string, rePassword string) (userData *userdm.User, err error) {
	var (
		logger  = logutil.GetLogger()
		isFound bool
	)
	userData, isFound, err = u.repo.GetExistingUserByEmail(ctx, email)
	if err != nil {
		errMsg := fmt.Sprintf("Error while getting user by email, err: %s", err.Error())
		logger.Errorf(errMsg)
		return userData, errors.New(errMsg)
	}
	if isFound {
		errMsg := fmt.Sprintf("User existed with email: %s", email)
		logger.Error(errMsg)
		return userData, errors.New(errMsg)
	}
	if password != rePassword {
		errMsg := fmt.Sprint("password and repeated password are not matched")
		logger.Error(errMsg)
		return userData, errors.New(errMsg)
	}
	return userData, nil
}

func (u *userService) UserLogin(ctx context.Context, loginReq UserLogin) (user *UserInfo, err error) {
	var (
		logger   = logutil.GetLogger()
		userData *userdm.User
		isFound  bool
	)
	userData, isFound, err = u.repo.GetExistingUserByEmail(ctx, loginReq.Email)
	if err != nil {
		errMsg := fmt.Sprintf("Error while getting user by email, err: %s", err.Error())
		logger.Errorf(errMsg)
		return user, errors.New(errMsg)
	}
	if !isFound {
		errMsg := fmt.Sprintf("User not existed with email: %s", loginReq.Email)
		logger.Error(errMsg)
		return user, errors.New(errMsg)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(loginReq.Password)); err != nil {
		logger.Errorf("password not matched, err: %s", err.Error())
		return nil, errors.New("password not matched")
	}
	user = &UserInfo{
		UserName: userData.UserEmail,
		UserID:   userData.UserID,
	}
	return user, nil
}
