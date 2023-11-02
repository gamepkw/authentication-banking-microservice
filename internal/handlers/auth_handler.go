package handler

import (
	"github.com/labstack/echo/v4"

	model "github.com/gamepkw/authentication-banking-microservice/internal/models"
	authService "github.com/gamepkw/authentication-banking-microservice/internal/services"
)

type Response struct {
	Message string              `json:"message"`
	Body    *model.UserResponse `json:"body,omitempty"`
}

type AuthHandler struct {
	authService authService.AuthService
}

func NewAuthHandler(e *echo.Echo, auths authService.AuthService) {
	handler := &AuthHandler{
		authService: auths,
	}
	e.POST("/users/send-otp", handler.SendOtp)
	e.POST("/users/verify-otp", handler.VerifyOtp)
}
