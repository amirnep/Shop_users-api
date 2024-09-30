package users

import (
	"net/http"
	"os"
	"strconv"

	"github.com/amirnep/shop/src/domain/users"
	"github.com/amirnep/shop/src/jwt"
	"github.com/amirnep/shop/src/services"
	crypto_utils "github.com/amirnep/shop/src/utils/cypto_utils"
	"github.com/amirnep/shop/src/utils/errors"
	"github.com/gin-gonic/gin"
)

var (
	UsersController usersControllerInterface = &usersController{}
)

type usersController struct {}

type usersControllerInterface interface {
	getUserId(string) (int64, *errors.RestErr)
	Create(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
}

type LoginInput struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *usersController) getUserId(userIdParam string) (int64, *errors.RestErr) {
	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)
	if userErr != nil {
		return 0, errors.NewBadRequestError("user id should be a number")
	}
	return userId, nil
}

func (u *usersController) Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := *errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

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

	var user users.Profile
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId

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
	var input LoginInput
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