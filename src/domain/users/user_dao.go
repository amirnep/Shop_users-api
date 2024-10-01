package users

import (
	"github.com/amirnep/shop/src/datasources/mysql/users_db"
	"github.com/amirnep/shop/src/logger"
	"github.com/amirnep/shop/src/utils/errors"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, password, confirm_password) VALUES (?,?,?,?,?);"

	queryGetUser = "SELECT id, first_name, last_name, email, role, date_created FROM users WHERE id = ?;"

	queryUpdateUser = "UPDATE users SET first_name=?, last_name=? WHERE id = ?;"

	queryDeleteUser = "DELETE FROM users WHERE id = ?;"

	queryGetLoginInfo = "SELECT id, email, role, password FROM users WHERE email = ?;"

	queryGetUsers = "SELECT id, first_name, last_name, email, role, date_created FROM users;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.DateCreated); getErr != nil {
		logger.Error("error when trying to get user by id", getErr)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Password, user.ConfirmPassword)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return errors.NewInternalServerError("database error")
	}

	userId, insertErr := insertResult.LastInsertId()
	if insertErr != nil {
		logger.Error("error when trying to get last insert id after creating a new user", insertErr)
		return errors.NewInternalServerError("database error")
	}

	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(user.FirstName, user.LastName, user.Id)
	if updateErr != nil {
		logger.Error("error when trying to update user", updateErr)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	if _, deleteErr := stmt.Exec(user.Id); deleteErr != nil {
		logger.Error("error when trying to delete user", deleteErr)
		return errors.NewInternalServerError("database error")
	}
		
	return nil
}

func (user *User) Login() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetLoginInfo)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email)
	if getErr := result.Scan(&user.Id, &user.Email, &user.Role, &user.Password); getErr != nil {
		logger.Error("error when trying to get user by email", getErr)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) GetAll() ([]User ,*errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryGetUsers)
	if err != nil {
		logger.Error("error when trying to get users statement", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result, queryErr := stmt.Query()
	if queryErr != nil {
		return nil, errors.NewInternalServerError("query error")
	}

	var users []User
	
	for result.Next() {
		if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.DateCreated); getErr != nil {
			logger.Error("error when trying to get users.", getErr)
			return users , errors.NewInternalServerError("database error")
		}

		users = append(users, *user)
	}

	if resultErr := result.Err(); resultErr != nil {
		return nil, errors.NewInternalServerError(resultErr.Error())
	}

	return users, nil
}