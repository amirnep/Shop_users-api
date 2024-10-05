package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/amirnep/shop/src/domain/users"
	"github.com/amirnep/shop/src/jwt"
	"github.com/amirnep/shop/src/services"
	crypto_utils "github.com/amirnep/shop/src/utils/cypto_utils"
	"github.com/amirnep/shop/src/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	UsersController usersControllerInterface = &usersController{}
)

type usersController struct {}

type usersControllerInterface interface {
	getUserId(string) (int64, *errors.RestErr)
	GetUsers(c *gin.Context)
	Create(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateRole(c *gin.Context)
	ChangePassword(c *gin.Context)
}

func (u *usersController) getUserId(userIdParam string) (int64, *errors.RestErr) {
	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)
	if userErr != nil {
		return 0, errors.NewBadRequestError("user id should be a number")
	}
	return userId, nil
}

func (u *usersController) GetUsers(c *gin.Context) {
	result, getErr := services.UsersService.GetAll()
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (u *usersController) Create(c *gin.Context) {
	var user *users.User
	if err := c.ShouldBind(&user); err != nil {
		restErr := *errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	file := c.Request.MultipartForm.File["Image"][0]

	if file.Size > 3<<20 {
		restErr := *errors.NewBadRequestError("image size must less then 3mb")
		c.JSON(restErr.Status, restErr)
		return
	}

	uniqueId := uuid.New().String()
	dst := "wwwroot/" + filepath.Base(uniqueId + ".jpg")

	if err := c.SaveUploadedFile(file, dst); err != nil {
		restErr := *errors.NewBadRequestError("error in uploading and saving file")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.ImageUrl = "src/wwwroot/" + uniqueId + ".jpg"

	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		c.JSON(http.StatusBadRequest,saveErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func (u *usersController) Get(c *gin.Context) {
	userId, idErr := UsersController.getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	result, getErr := services.UsersService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func (u *usersController) Update(c *gin.Context) {
	userId, idErr := jwt.JWTUserId(c)
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	var user users.User

	if err := c.ShouldBind(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	file := c.Request.MultipartForm.File["Image"][0]

	if file.Size > 3<<20 {
		restErr := *errors.NewBadRequestError("image size must less then 3mb")
		c.JSON(restErr.Status, restErr)
		return
	}

	uniqueId := uuid.New().String()
	dst := "wwwroot/" + filepath.Base(uniqueId + ".jpg")

	if err := c.SaveUploadedFile(file, dst); err != nil {
		restErr := *errors.NewBadRequestError("error in uploading and saving file")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId
	user.ImageUrl = "src/wwwroot/" + uniqueId + ".jpg"

	isPartial := c.Request.Method == http.MethodPatch

	result, err := services.UsersService.UpdateUser(isPartial, user)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func (u *usersController) Delete(c *gin.Context) {
	userId, idErr := UsersController.getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	if err := services.UsersService.DeleteUser(userId); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (u *usersController) Login(c *gin.Context){
	input := users.LoginInput{}
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))

	if inputErr := c.ShouldBindJSON(&input); inputErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": inputErr.Error()})
		return
	}

	result, loginErr := services.UsersService.Login(input.Email)
	if loginErr!= nil {
        c.JSON(loginErr.Status, loginErr)
		return
	}

	inputPassword := crypto_utils.GetMd5(input.Password)
	if inputPassword != result.Password{
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is incorrect"})
		return
	}

	token, tokenErr := jwt.GenerateJWT(*result)

	if tokenErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": result.Id, "token_type": "Bearer", "expires_in": tokenTTL,"access_token": token})
}

func (u *usersController) GetProfile(c *gin.Context) {
	userId, idErr := jwt.JWTUserId(c)
	if idErr != nil {
		c.JSON(idErr.Status, userId)
		return
	}
	result, getErr := services.UsersService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func (u *usersController) UpdateRole(c *gin.Context) {
	userId, idErr := UsersController.getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	result := services.UsersService.EditRole(userId)
	if result != nil {
		c.JSON(http.StatusBadRequest, result)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role edited to Admin successfully"})
}

func (u *usersController) ChangePassword(c *gin.Context) {
	userId, idErr := jwt.JWTUserId(c)
	if idErr != nil {
		c.JSON(idErr.Status, userId)
		return
	}

	var user *users.Password

	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId

	result := services.UsersService.EditPassword(userId, user)
	if result != nil {
		c.JSON(result.Status, result)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}