package handler

import (
	"net/http"

	model "github.com/gamepkw/authentication-banking-microservice/internal/models"
	"github.com/labstack/echo/v4"
)

func (auth *AuthHandler) SendOtp(c echo.Context) (err error) {

	var set model.UpdatePassword
	if err = c.Bind(&set); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()

	if err = auth.authService.SendOtp(ctx, set.Tel); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Response{Message: "Send otp successfully", Body: nil})

}
