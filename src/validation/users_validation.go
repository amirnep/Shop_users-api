package validation

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/amirnep/shop/src/domain/users"
	"github.com/amirnep/shop/src/utils/errors"
)

func EmailValidation(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func PasswordValidation(password string) bool {
	var (
        hasMinLen  = false
        hasUpper   = false
        hasLower   = false
        hasNumber  = false
        hasSpecial = false
    )
    if len(password) >= 8 {
        hasMinLen = true
    }
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func Validate(user *users.User) *errors.RestErr {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" || !EmailValidation(user.Email){
		return errors.NewBadRequestError("invalid email address.")
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" || !PasswordValidation(user.Password){
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

func ChangePasswordValidation(password *users.Password) *errors.RestErr {
	password.Password = strings.TrimSpace(password.Password)
	if password.Password == "" || !PasswordValidation(password.Password){
		return errors.NewBadRequestError("Password must have upperLetter, lowerLetter, number, specialChar, and longer than 8.")
	}

	password.ConfirmPassword = strings.TrimSpace(password.ConfirmPassword)
	if password.ConfirmPassword == ""{
		return errors.NewBadRequestError("invalid confirm password.")
	}

	if password.Password != password.ConfirmPassword {
		return errors.NewBadRequestError("passwords does not match.")
	}

	return nil
}