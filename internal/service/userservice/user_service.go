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
		logger = logutil.GetLogger()
	)
	if userData, err = u.repo.GetExistingUserByEmail(ctx, email); err != nil {
		logger.Errorf("Error while finding user, err: %s", err.Error())
		return userData, errors.New("user existed")
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
	)
	if userData, err = u.repo.GetExistingUserByEmail(ctx, loginReq.Email); err != nil {
		logger.Errorf("Not found account, err: %s", err.Error())
		return nil, errors.New("user not existed")
	}
	if loginReq.Password != userData.Password {
		logger.Errorf("password not matched")
		return nil, errors.New("password not matched")
	}
	user = &UserInfo{
		UserName: userData.UserEmail,
		UserID:   userData.UserID,
	}
	return user, nil
}
