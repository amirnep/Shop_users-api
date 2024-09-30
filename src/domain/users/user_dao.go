package users

import (
	"github.com/amirnep/shop/datasources/mysql/users_db"
	"github.com/amirnep/shop/logger"
	"github.com/amirnep/shop/utils/errors"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, password, confirm_password) VALUES (?,?,?,?,?);"

	queryGetUser = "SELECT id, first_name, last_name, email, role, date_created FROM users WHERE id = ?;"

	queryUpdateUser = "UPDATE users SET first_name=?, last_name=? WHERE id = ?;"

	queryDeleteUser = "DELETE FROM users WHERE id = ?;"

	queryGetLoginInfo = "SELECT id, email, role, password FROM users WHERE email = ?;"
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
		logger.Error("error when trying to get user by id", err)
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
		logger.Error("error when trying to save user", err)
		return errors.NewInternalServerError("database error")
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
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

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Id)
	if err != nil {
		logger.Error("error when trying to update user", err)
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

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user", err)
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
		logger.Error("error when trying to get user by email", err)
		return errors.NewInternalServerError("database error")
	}
	return nil
}