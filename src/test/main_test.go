package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amirnep/shop/src/controllers"
	"github.com/amirnep/shop/src/domain/users"
	"github.com/amirnep/shop/src/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
    router := gin.Default()
    return router
}

func TestRegister(t *testing.T) {
    r := SetUpRouter()
    r.POST("/Register", controllers.UsersController.Create)

    user := users.User{
        FirstName: "amir",
        LastName: "nep",
        Email: "test2@test.com",
        Password : "T@1est12459",
        ConfirmPassword:"T@1est12459",
    }

    jsonValue, _ := json.Marshal(user)
    req, _ := http.NewRequest("POST", "/Register", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "multipart/form-data")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLogin(t *testing.T) {
    r := SetUpRouter()
    r.POST("/Login", controllers.UsersController.Login)

    user := users.LoginInput{
        Email: "test2@test.com",
        Password : "T@1est12459",
    }

    jsonValue, _ := json.Marshal(user)
    req, _ := http.NewRequest("POST", "/Login", bytes.NewBuffer(jsonValue))

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetProfile(t *testing.T) {
    r := SetUpRouter()
    r.Use(middlewares.JWTAuthCustomerMiddleware())
    r.GET("/api/GetProfile", controllers.UsersController.GetProfile)

    req, _ := http.NewRequest("GET", "/api/GetProfile", nil)
    req.Header.Set("Authorization" , "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlYXQiOjE3MjgyNDEyNjEsImlhdCI6MTcyODE1NDg2MSwiaWQiOjIxLCJyb2xlIjoiYWRtaW4ifQ.UGPP5Hh0R4c6QvUuPMZ2d6pcRw-yGjaZfxkNDQ2O1us")
    req.Header.Set("X-Public", "true")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    var user []users.User
    json.Unmarshal(w.Body.Bytes(), &user)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.NotEmpty(t, w.Body.String())
}

func TestGetUser(t *testing.T) {
    r := SetUpRouter()
    r.Use(middlewares.JWTAuthMiddleware())
    r.GET("/api/admin/GetUser/:user_id", controllers.UsersController.Get)

    req, _ := http.NewRequest("GET", "/api/admin/GetUser/14", nil)
    req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlYXQiOjE3MjgyNDEyNjEsImlhdCI6MTcyODE1NDg2MSwiaWQiOjIxLCJyb2xlIjoiYWRtaW4ifQ.UGPP5Hh0R4c6QvUuPMZ2d6pcRw-yGjaZfxkNDQ2O1us")
    req.Header.Set("X-Public", "true")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    var user []users.User
    json.Unmarshal(w.Body.Bytes(), &user)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.NotEmpty(t, w.Body.String())
}

func TestDeleteUser(t *testing.T) {
    r := SetUpRouter()
    r.Use(middlewares.JWTAuthMiddleware())
    r.DELETE("/api/admin/DeleteUser/:user_id", controllers.UsersController.Delete)

    req, _ := http.NewRequest("DELETE", "/api/admin/DeleteUser/13", nil)
    req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlYXQiOjE3MjgyNDEyNjEsImlhdCI6MTcyODE1NDg2MSwiaWQiOjIxLCJyb2xlIjoiYWRtaW4ifQ.UGPP5Hh0R4c6QvUuPMZ2d6pcRw-yGjaZfxkNDQ2O1us")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.NotEmpty(t, w.Body.String())
}

func TestGetUsers(t *testing.T) {
    r := SetUpRouter()
    r.Use(middlewares.JWTAuthMiddleware())
    r.GET("/api/admin/GetUsers", controllers.UsersController.GetUsers)

    req, _ := http.NewRequest("GET", "/api/admin/GetUsers", nil)
    req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlYXQiOjE3MjgyNDEyNjEsImlhdCI6MTcyODE1NDg2MSwiaWQiOjIxLCJyb2xlIjoiYWRtaW4ifQ.UGPP5Hh0R4c6QvUuPMZ2d6pcRw-yGjaZfxkNDQ2O1us")
    
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    var user []users.User
    json.Unmarshal(w.Body.Bytes(), &user)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.NotEmpty(t, w.Body.String())
}

func TestUpdateRole(t *testing.T) {
    r := SetUpRouter()
    r.Use(middlewares.JWTAuthMiddleware())
    r.PUT("api/admin/EditRole/:user_id", controllers.UsersController.UpdateRole)

    req, _ := http.NewRequest("PUT", "/api/admin/EditRole/28", nil)
    req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlYXQiOjE3MjgyNDEyNjEsImlhdCI6MTcyODE1NDg2MSwiaWQiOjIxLCJyb2xlIjoiYWRtaW4ifQ.UGPP5Hh0R4c6QvUuPMZ2d6pcRw-yGjaZfxkNDQ2O1us")
    
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}

func TestChangePassword(t *testing.T) {
    r := SetUpRouter()
    r.Use(middlewares.JWTAuthCustomerMiddleware())
    r.PUT("api/ChangePassword", controllers.UsersController.ChangePassword)

    user := users.Password{
        Password : "T@1est12459",
        ConfirmPassword: "T@1est12459",
    }

    jsonValue, _ := json.Marshal(user)
    req, _ := http.NewRequest("PUT", "/api/ChangePassword", bytes.NewBuffer(jsonValue))
    req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlYXQiOjE3MjgyNDEyNjEsImlhdCI6MTcyODE1NDg2MSwiaWQiOjIxLCJyb2xlIjoiYWRtaW4ifQ.UGPP5Hh0R4c6QvUuPMZ2d6pcRw-yGjaZfxkNDQ2O1us")
    
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}