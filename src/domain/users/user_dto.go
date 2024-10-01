package users

import (
	"strings"

	"github.com/amirnep/shop/src/utils/errors"
	"github.com/amirnep/shop/src/validation"
)

type User struct {
	Id              int64  `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email" binding:"required"`
	Role            string `json:"role"`
	DateCreated 	string `json:"date_created"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type Users []User

type LoginInput struct {
	Email 			string `json:"email" binding:"required"`
	Password 		string `json:"password" binding:"required"`
}

type Profile struct {
	Id              int64  `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
}

func (user *User) Validate() *errors.RestErr {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" || !validation.EmailValidation(user.Email){
		return errors.NewBadRequestError("invalid email address.")
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" || !validation.PasswordValidation(user.Password){
		return errors.NewBadRequestError("Password must have upperLetter, lowerLetter, number, specialChar, and longer than 8.")
	}

	user.ConfirmPassword = strings.TrimSpace(user.ConfirmPassword)
	if user.ConfirmPassword == ""{
		return errors.NewBadRequestError("invalid confirm password.")
	}

	if user.Password != user.ConfirmPassword {
		return errors.NewBadRequestError("passwords does not match.")
	}

	return nil
}