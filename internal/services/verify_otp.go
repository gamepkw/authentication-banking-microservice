package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gamepkw/authentication-banking-microservice/internal/utils"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func (auth *authService) VerifyOtp(c context.Context, tel string, otpUser string) bool {
	ctx, cancel := context.WithTimeout(c, auth.contextTimeout)
	defer cancel()

	secretKey, expiredAt, err := auth.getSecretKeyByUUID(ctx, tel)

	if secretKey == "" {
		fmt.Println("key error")
		return false
	}

	if err != nil {
		fmt.Println("OTP error")
		return false
	}

	if expiredAt.Before(time.Now()) {
		fmt.Println("OTP expired")
		return false
	}

	validateOpts := totp.ValidateOpts{
		Period:    180,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	}

	valid, err := totp.ValidateCustom(otpUser, secretKey, time.Now(), validateOpts)
	if err != nil {
		fmt.Println("OTP error")
		return false
	}

	return valid
}

func (auth *authService) getSecretKeyByUUID(c context.Context, tel string) (string, time.Time, error) {
	ctx, cancel := context.WithTimeout(c, auth.contextTimeout)
	defer cancel()

	utils.EncodeBase64(&tel)

	secretKey, expiredAt, _ := auth.authenticationRepo.GetOtpSecret(ctx, tel)

	return secretKey, expiredAt, nil
}
