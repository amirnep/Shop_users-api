package services

import (
	"github.com/amirnep/shop/domain/users"
	crypto_utils "github.com/amirnep/shop/utils/cypto_utils"

	"github.com/amirnep/shop/utils/date_utils"
	"github.com/amirnep/shop/utils/errors"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct{}

type usersServiceInterface interface {
	GetUser(int64) (*users.User, *errors.RestErr)
	CreateUser(users.User)(*users.User,*errors.RestErr)
	UpdateUser(bool, users.Profile) (*users.User, *errors.RestErr)
	DeleteUser(int64) *errors.RestErr
	Login(string) (*users.User, *errors.RestErr)
	GetProfile(int64) (*users.User, *errors.RestErr)
}

func (s *usersService) GetUser(userId int64) (*users.User, *errors.RestErr) {
	dao := &users.User{Id: userId}
	if err := dao.Get(); err != nil {
		return nil, err
	}
	return dao, nil
}

func (s *usersService) CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.Role = "user"
	user.DateCreated = date_utils.GetNowDBFormat()
	user.Password = crypto_utils.GetMd5(user.Password)
	user.ConfirmPassword = crypto_utils.GetMd5(user.ConfirmPassword)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) UpdateUser(isPartial bool, user users.Profile) (*users.User, *errors.RestErr) {
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
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
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