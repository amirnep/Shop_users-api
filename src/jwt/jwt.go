package jwt

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/amirnep/shop/src/domain/users"
	"github.com/amirnep/shop/src/utils/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

func GenerateJWT(user users.User) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.Id,
		"role": user.Role,
		"iat":  time.Now().Unix(),
		"eat":  time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(privateKey)
}

func ValidateJWT(context *gin.Context) *errors.RestErr {
	token, err := getToken(context)
	if err != nil {
		return errors.NewBadRequestError("invalid token provided")
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.NewBadRequestError("invalid token provided")
}

func getToken(context *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

func getTokenFromRequest(context *gin.Context) string {
	bearerToken := context.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

func ValidateAdminRoleJWT(context *gin.Context) *errors.RestErr {
	token, err := getToken(context)
	if err != nil {
		return errors.NewBadRequestError("invalid admin token provided")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := string(claims["role"].(string))
	if ok && token.Valid && userRole == "admin" {
		return nil
	}
	return errors.NewBadRequestError("invalid admin token provided")
}

func ValidateCustomerRoleJWT(context *gin.Context) *errors.RestErr {
	token, err := getToken(context)
	if err != nil {
		return errors.NewBadRequestError("invalid author token provided")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := string(claims["role"].(string))
	if ok && token.Valid && userRole == "user" || userRole == "admin" {
		return nil
	}
	return errors.NewBadRequestError("invalid author token provided")
}

func JWTUserId(context *gin.Context) (int64, *errors.RestErr) {
	token, err := getToken(context)
	if err != nil {
		return 0, errors.NewBadRequestError("invalid token provided")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userId := int64(claims["id"].(float64))
	if ok && token.Valid {
		return userId, nil
	}
	return 0, errors.NewBadRequestError("Only registered Customers are allowed to perform this action")
}