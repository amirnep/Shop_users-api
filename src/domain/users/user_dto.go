package users

import (
	"mime/multipart"
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
	ImageUrl		string `json:"image_url"`
	Image			*multipart.FileHeader `form:"file"`
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

type Password struct {
	Id              int64  `json:"id"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}