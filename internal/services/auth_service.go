package service

import (
	"context"
	"time"

	authRepo "github.com/gamepkw/authentication-banking-microservice/internal/repositories"
)

type authService struct {
	authenticationRepo authRepo.AuthRepository
	contextTimeout     time.Duration
}

func NewauthService(auth authRepo.AuthRepository, timeout time.Duration) *authService {
	return &authService{
		authenticationRepo: auth,
		contextTimeout:     timeout,
	}
}

type AuthService interface {
	GenerateOtp(c context.Context, tel string) (string, error)
	SendOtp(c context.Context, tel string) error
	VerifyOtp(c context.Context, userOTP string, secretKey string) bool
}
