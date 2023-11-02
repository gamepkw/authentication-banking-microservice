package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"time"

	"github.com/gamepkw/authentication-banking-microservice/internal/utils"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func (auth *authService) GenerateOtp(c context.Context, tel string) (string, error) {
	ctx, cancel := context.WithTimeout(c, auth.contextTimeout)
	defer cancel()
	secretKey, err := generateRandomSecretKey()
	if err != nil {
		return "", err
	}

	validateOpts := totp.ValidateOpts{
		Period:    180,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	}

	otp, err := totp.GenerateCodeCustom(secretKey, time.Now(), validateOpts)
	if err != nil {
		return "", err
	}

	auth.saveOtpSecret(ctx, tel, secretKey)

	return otp, nil
}

func (auth *authService) saveOtpSecret(c context.Context, uuid string, secretKey string) (err error) {
	ctx, cancel := context.WithTimeout(c, auth.contextTimeout)
	defer cancel()

	if err = utils.EncodeBase64(&uuid); err != nil {
		return
	}

	if err = auth.authenticationRepo.SaveOtpSecret(ctx, uuid, secretKey); err != nil {
		return err
	}

	return nil
}

func generateRandomSecretKey() (string, error) {
	key := make([]byte, 16) // Generate a 16-byte random key
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(key), nil
}
