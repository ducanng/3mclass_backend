package userservice

import (
	"context"
	"errors"
	"fmt"
	"net/mail"

	"golang.org/x/crypto/bcrypt"

	"github.com/ducanng/3mclass_backend/config"
	"github.com/ducanng/3mclass_backend/internal/model/userdm"
	"github.com/ducanng/3mclass_backend/internal/repository"
	"github.com/ducanng/3mclass_backend/pkg/logutil"
)

type UserService interface {
	RegisterUser(ctx context.Context, req *UserRegistrationRequest) (*UserRegistrationResponse, error)
	UserLogin(ctx context.Context, login UserLogin) (*UserInfo, error)
	GetUserInfo(ctx context.Context, userID uint64) (*UserProfile, error)
	UpdateUserProfile(ctx context.Context, userID uint64, req *UpdateUserProfileRequest) error
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

func (s *userService) RegisterUser(ctx context.Context, req *UserRegistrationRequest) (resp *UserRegistrationResponse, err error) {
	var (
		logger   = logutil.GetLogger()
		userData *userdm.User
	)
	if userData, err = s.validateUserAndPassword(ctx, req.Email, req.Password, req.RePassword); err != nil {
		logger.Errorf("error while validating user and password, err: %s", err.Error())
		return nil, err
	}
	var encryptedPassword []byte
	if encryptedPassword, err = bcrypt.GenerateFromPassword([]byte(req.Password), 11); err != nil {
		logger.Errorf("Generate password failed, err: %s", err.Error())
		return nil, err
	}

	userData, err = s.repo.CreateUser(ctx, &userdm.User{
		UserEmail:    req.Email,
		DisplayName:  fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		FirstName:    req.LastName,
		LastName:     req.FirstName,
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

func (s *userService) validateUserAndPassword(ctx context.Context, email string, password string, rePassword string) (userData *userdm.User, err error) {
	var (
		logger  = logutil.GetLogger()
		isFound bool
	)
	userData, isFound, err = s.repo.GetExistingUserByEmail(ctx, email)
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

func (s *userService) UserLogin(ctx context.Context, loginReq UserLogin) (user *UserInfo, err error) {
	var (
		logger   = logutil.GetLogger()
		userData *userdm.User
		isFound  bool
	)
	userData, isFound, err = s.repo.GetExistingUserByEmail(ctx, loginReq.Email)
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

func (s *userService) GetUserInfo(ctx context.Context, userID uint64) (user *UserProfile, err error) {
	var (
		logger   = logutil.GetLogger()
		userData *userdm.User
		isFound  bool
	)
	userData, isFound, err = s.repo.GetExistingUserByID(ctx, userID)
	if err != nil {
		logger.Errorf("error while getting user by id, err: %s", err.Error())
		return nil, err
	}
	if !isFound {
		logger.Errorf("user not found with id: %d", userID)
		return nil, errors.New("user not found")
	}
	user = &UserProfile{
		UserID:      userData.UserID,
		FirstName:   userData.FirstName,
		LastName:    userData.LastName,
		Email:       userData.UserEmail,
		PhoneNumber: userData.PhoneNumber,
		DisplayName: userData.DisplayName,
	}
	return user, nil
}

func (s *userService) UpdateUserProfile(ctx context.Context, userID uint64, req *UpdateUserProfileRequest) (err error) {
	var (
		logger   = logutil.GetLogger()
		userData *userdm.User
		isFound  bool
	)
	userData, isFound, err = s.repo.GetExistingUserByID(ctx, userID)
	if err != nil {
		logger.Errorf("error while getting user by id, err: %s", err.Error())
		return err
	}
	if !isFound {
		logger.Errorf("user not found with id: %d", userID)
		return errors.New("user not found")
	}
	if req.Email != "" {
		if _, err = mail.ParseAddress(req.Email); err != nil {
			logger.Errorf("Invalid input email %s, err:%s ", req.Email, err)
			return err
		}
		userData.UserEmail = req.Email
	}
	if req.PhoneNumber != "" {
		userData.PhoneNumber = req.PhoneNumber
	}
	if req.FirstName != "" {
		userData.FirstName = req.FirstName
	}
	if req.LastName != "" {
		userData.LastName = req.LastName
	}
	userData.DisplayName = fmt.Sprintf("%s %s", userData.FirstName, userData.LastName)
	if err = s.repo.UpdateUser(ctx, userData); err != nil {
		logger.Errorf("error while updating user, err: %s", err.Error())
		return err
	}
	return nil
}
