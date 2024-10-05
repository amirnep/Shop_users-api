package services

import (
	"github.com/amirnep/shop/src/domain/users"
	crypto_utils "github.com/amirnep/shop/src/utils/cypto_utils"
	"github.com/amirnep/shop/src/validation"

	"github.com/amirnep/shop/src/utils/date_utils"
	"github.com/amirnep/shop/src/utils/errors"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct{}

type usersServiceInterface interface {
	GetUser(int64) (*users.User, *errors.RestErr)
	GetAll() ([]users.User, *errors.RestErr)
	CreateUser(*users.User)(*users.User,*errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, *errors.RestErr)
	DeleteUser(int64) *errors.RestErr
	Login(string) (*users.User, *errors.RestErr)
	GetProfile(int64) (*users.User, *errors.RestErr)
	EditRole(int64) *errors.RestErr
	EditPassword(int64, *users.Password) *errors.RestErr
}

func (s *usersService) GetUser(userId int64) (*users.User, *errors.RestErr) {
	dao := &users.User{Id: userId}
	if err := dao.Get(); err != nil {
		return nil, err
	}
	return dao, nil
}

func (s *usersService) GetAll() ([]users.User, *errors.RestErr) {
	dao := &users.User{}
	res, err := dao.GetAll(); if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *usersService) CreateUser(user *users.User) (*users.User, *errors.RestErr) {
	if err := validation.Validate(user); err != nil {
		return nil, err
	}

	user.Role = "user"
	user.DateCreated = date_utils.GetNowDBFormat()
	user.Password = crypto_utils.GetMd5(user.Password)
	user.ConfirmPassword = crypto_utils.GetMd5(user.ConfirmPassword)

	if err := user.Save(); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *usersService) UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr) {
	current := &users.User{Id: user.Id}
	if err := current.Get(); err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}

		if user.LastName != "" {
			current.LastName = user.LastName
		}

		if user.ImageUrl != "" {
			current.ImageUrl = user.ImageUrl
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.ImageUrl = user.ImageUrl
	}

	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

func (s *usersService) DeleteUser(userId int64) *errors.RestErr {
	current := &users.User{Id: userId}
	if err := current.Get(); err != nil {
		return errors.NewBadRequestError("user does not exist")
	}
	return current.Delete()
}

func (s *usersService) Login(email string) (*users.User, *errors.RestErr) {
	dao := &users.User{Email: email}
	if err := dao.Login(); err != nil {
		return nil, err
	}
	return dao, nil
}

func (s *usersService) GetProfile(userId int64) (*users.User, *errors.RestErr) {
	dao := &users.User{Id: userId}
	if err := dao.Get(); err != nil {
		return nil, err
	}
	return dao, nil
}

func (s *usersService) EditRole(userId int64) (*errors.RestErr) {
	current := &users.User{Id: userId}
	if err := current.Get(); err != nil {
		return err
	}

	current.Role = "admin"

	if err := current.EditRole(); err != nil {
		return err
	}
	return nil
}

func (s *usersService) EditPassword(userId int64, user *users.Password) (*errors.RestErr) {
	current := &users.User{Id: userId}
	if err := current.Get(); err != nil {
		return err
	}

	if validation := validation.ChangePasswordValidation(user); validation != nil {
		return errors.NewBadRequestError("password not valid")
	}

	current.Password = crypto_utils.GetMd5(user.Password)
	current.ConfirmPassword = crypto_utils.GetMd5(user.ConfirmPassword)

	if err := current.EditPassword(); err != nil {
		return err
	}
	return nil
}