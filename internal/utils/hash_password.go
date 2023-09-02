package utils

import (
	model "github.com/gamepkw/users-banking-microservice/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func HashPasswordBcrypt(input *string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(*input), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*input = string(hashedBytes)
	return nil
}

func ComparePasswords(hashedPassword, Password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(Password))
	if err != nil {
		return model.ErrWrongPassword
	}
	return err
}
