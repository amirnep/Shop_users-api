package users

import (
	"encoding/json"
)

type PublicUser struct {
	Id          int64  `json:"id"`
	Role 		string `json:"role"`
	DateCreated string `json:"date_created"`
}

type PrivateUser struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Role 		string `json:"role"`
	DateCreated string `json:"date_created"`
	ImageUrl	string `json:"image_url"`
}

func (users Users) Marshall(isPublic bool) []interface {} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.Marshall(isPublic)
	}
	return result
}

func (user *User) Marshall(isPublic bool) interface{} {
	if isPublic {
		return PublicUser{
			Id:          user.Id,
			Role: 		 user.Role,
			DateCreated: user.DateCreated,
		}
	}
	userJson, _ := json.Marshal(user)
	var privateUser PrivateUser
	json.Unmarshal(userJson, &privateUser)
	return privateUser
}