package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"main/model"
)

// ResponseError represent the response error struct

// UserHandler  represent the httphandler for user
type AuthenticationHandler struct {
	AuthService model.AuthenticationService
}

// NewUserHandler will initialize the users/ resources endpoint
func NewAuthenticationHandler(e *echo.Echo, auths model.AuthenticationService) {
	handler := &AuthenticationHandler{
		AuthService: auths,
	}

	e.POST("/users/send-otp", handler.SendOtp)
	e.POST("/users/verify-otp", handler.ValidateOtp)

}

func (auth *AuthenticationHandler) SendOtp(c echo.Context) (err error) {

	var set model.UpdatePassword
	if err = c.Bind(&set); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()

	if err = auth.AuthService.SendOtp(ctx, set.Tel); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Response{Message: "Send otp successfully", Body: nil})

}

func (auth *AuthenticationHandler) ValidateOtp(c echo.Context) (err error) {
	var set model.UpdatePassword

	if err = c.Bind(&set); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()

	if !auth.AuthService.ValidateOtp(ctx, set.Tel, set.Otp) {
		return c.JSON(http.StatusBadRequest, Response{Message: "Otp is invalid", Body: nil})
	}

	return c.JSON(http.StatusOK, Response{Message: "Otp is valid", Body: nil})

}
