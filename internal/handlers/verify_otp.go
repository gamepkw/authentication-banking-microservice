package handler

import (
	"net/http"

	model "github.com/gamepkw/authentication-banking-microservice/internal/models"
	"github.com/labstack/echo/v4"
)

func (auth *AuthHandler) VerifyOtp(c echo.Context) (err error) {
	var set model.UpdatePassword

	if err = c.Bind(&set); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()

	if !auth.authService.VerifyOtp(ctx, set.Tel, set.Otp) {
		return c.JSON(http.StatusBadRequest, Response{Message: "Otp is invalid", Body: nil})
	}

	return c.JSON(http.StatusOK, Response{Message: "Otp is valid", Body: nil})

}
