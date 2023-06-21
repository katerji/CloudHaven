package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/UserAuthKit/input"
	"github.com/katerji/UserAuthKit/service"
)

const RegisterPath = "/register"

const (
	errorMessageEmailAlreadyExists = "Email already exists."
	userRegisteredSuccessfully     = "User registered successfully."
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func RegisterHandler(c *gin.Context) {
	request := &RegisterRequest{}
	err := c.BindJSON(request)
	if err != nil {
		sendBadRequest(c)
		return
	}
	registerUserInput := input.AuthInput{
		Email:    request.Email,
		Password: request.Password,
		Name:     request.Name,
	}
	userService := service.GetAuthService()
	_, err = userService.Register(registerUserInput)
	if err != nil {
		sendErrorMessage(c, errorMessageEmailAlreadyExists)
		return
	}
	sendResponseMessage(c, userRegisteredSuccessfully)
	return
}
