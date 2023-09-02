package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gamepkw/users-banking-microservice/internal/utils"

	model "github.com/gamepkw/users-banking-microservice/internal/models"

	"github.com/gamepkw/users-banking-microservice/internal/middleware"
)

type userService struct {
	userRepo       model.UserRepository
	contextTimeout time.Duration
}

// NewUserService will create new an userService object representation of model.UserService interface
func NewUserService(ur model.UserRepository, timeout time.Duration) model.UserService {
	return &userService{
		userRepo:       ur,
		contextTimeout: timeout,
	}
}

func (a *userService) RegisterUser(c context.Context, u *model.User) (res model.UserResponse, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	if err := utils.EncodeBase64(&u.Tel); err != nil {
		return res, err
	}

	if utils.ValidatePassword(u.Password) {

		if err := utils.HashPasswordBcrypt(&u.Password); err != nil {
			return res, err
		}
		res, err = a.userRepo.RegisterUser(ctx, u)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func (a *userService) Login(c context.Context, u *model.User) (token string, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	userResponse, err := a.getHashedPasswordByUUID(ctx, u.Tel)
	if err != nil {
		return "", err
	}

	if err = utils.ComparePasswords(userResponse.HashedPassword, u.Password); err != nil {
		return token, err
	}

	if err := utils.EncodeBase64(&u.Tel); err != nil {
		return "", err
	}

	token, err = middleware.GenerateJWTToken(u.Tel, 1*time.Hour)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *userService) SetNewPassword(c context.Context, u *model.UpdatePassword) (res model.UserResponse, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	userResponse, err := a.getHashedPasswordByUUID(ctx, u.Tel)
	if err != nil {
		return res, err
	}

	if err = utils.ComparePasswords(userResponse.HashedPassword, u.Password); err != nil {
		return res, err
	}

	if err = utils.EncodeBase64(&u.Tel); err != nil {
		return
	}

	if err = utils.HashPasswordBcrypt(&u.NewPassword); err != nil {
		return
	}

	fmt.Println(u.NewPassword)

	res, err = a.userRepo.SetPassword(ctx, u)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (a *userService) ResetPassword(c context.Context, u *model.UpdatePassword) (res model.UserResponse, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	if err = utils.EncodeBase64(&u.Tel); err != nil {
		return
	}

	if err = utils.HashPasswordBcrypt(&u.Password); err != nil {
		return
	}

	if utils.ValidatePassword(u.NewPassword) {

		if err := utils.HashPasswordBcrypt(&u.NewPassword); err != nil {
			return res, err
		}
		res, err = a.userRepo.SetPassword(ctx, u)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func (a *userService) SetUpPin(c context.Context, u *model.Pin) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	// if err := utils.EncodeBase64(&u.Tel); err != nil {
	// 	return err
	// }

	if err := utils.HashPinBcrypt(&u.Pin); err != nil {
		return err
	}

	if err = a.userRepo.SetUpPin(ctx, u); err != nil {
		return err
	}
	return nil
}

func (a *userService) SetNewPin(c context.Context, u *model.SetNewPin) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	// if err := utils.EncodeBase64(&u.Tel); err != nil {
	// 	return err
	// }

	pinDB, err := a.getHashedPinByUUID(ctx, u.Tel)
	if err != nil {
		return
	}

	if err = utils.ComparePins(pinDB.Pin, u.Pin); err != nil {
		return model.ErrWrongPassword
	}

	if err := utils.HashPinBcrypt(&u.NewPin); err != nil {
		return err
	}

	if err = a.userRepo.SetNewPin(ctx, u); err != nil {
		return err
	}

	return nil
}

func (a *userService) getHashedPasswordByUUID(c context.Context, tel string) (res *model.UserResponse, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	if err := utils.EncodeBase64(&tel); err != nil {
		return nil, err
	}

	if res, err = a.userRepo.GetHashedPasswordByUUID(ctx, tel); err != nil {
		return nil, err
	}

	return res, nil
}

func (a *userService) getHashedPinByUUID(c context.Context, uuid string) (res *model.Pin, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	if res, err = a.userRepo.GetHashedPinByUUID(ctx, uuid); err != nil {
		return nil, err
	}

	return res, nil
}

// var issuer = "MyApp"
// var accountName = "user@example.com"

// func (a *userService) GenerateOtp(c context.Context, tel string) (string, error) {
// 	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
// 	defer cancel()
// 	secretKey, err := generateRandomSecretKey()
// 	if err != nil {
// 		return "", err
// 	}

// 	validateOpts := totp.ValidateOpts{
// 		Period:    180,
// 		Skew:      1,
// 		Digits:    otp.DigitsSix,
// 		Algorithm: otp.AlgorithmSHA1,
// 	}

// 	otp, err := totp.GenerateCodeCustom(secretKey, time.Now(), validateOpts)
// 	if err != nil {
// 		return "", err
// 	}

// 	a.saveOtpSecret(ctx, tel, secretKey)

// 	return otp, nil
// }

// func (a *userService) SendOtp(c context.Context, tel string) error {
// 	topic := "sms"
// 	brokerAddress := viper.GetString("kafka.broker_address")
// 	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
// 	defer cancel()
// 	otp, err := a.GenerateOtp(ctx, tel)
// 	if err != nil {
// 		return err
// 	}

// 	producer.RunKafkaProducer(brokerAddress, topic, otp)
// 	return nil
// }

// func (a *userService) ValidateOtp(c context.Context, tel string, otpUser string) bool {
// 	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
// 	defer cancel()

// 	secretKey, expiredAt, err := a.getSecretKeyByUUID(ctx, tel)

// 	if secretKey == "" {
// 		fmt.Println("key error")
// 		return false
// 	}

// 	if err != nil {
// 		fmt.Println("OTP error")
// 		return false
// 	}

// 	if expiredAt.Before(time.Now()) {
// 		fmt.Println("OTP expired")
// 		return false
// 	}

// 	validateOpts := totp.ValidateOpts{
// 		Period:    180,
// 		Skew:      1,
// 		Digits:    otp.DigitsSix,
// 		Algorithm: otp.AlgorithmSHA1,
// 	}

// 	// valid := totp.Validate(otp, secretKey)
// 	valid, err := totp.ValidateCustom(otpUser, secretKey, time.Now(), validateOpts)
// 	if err != nil {
// 		fmt.Println("OTP error")
// 		return false
// 	}

// 	return valid
// }

func (a *userService) ValidatePin(c context.Context, uuid string, pin string) bool {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	// valid := totp.Validate(otp, secretKey)
	pinDB, err := a.getHashedPinByUUID(ctx, uuid)
	if err != nil {
		return false
	}

	if err = utils.ComparePins(pinDB.Pin, pin); err != nil {
		return false
	}

	return true
}

// func (a *userService) getSecretKeyByUUID(c context.Context, tel string) (string, time.Time, error) {
// 	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
// 	defer cancel()

// 	utils.EncodeBase64(&tel)

// 	secretKey, expiredAt, _ := a.userRepo.GetOtpSecret(ctx, tel)

// 	return secretKey, expiredAt, nil
// }

// func (a *userService) saveOtpSecret(c context.Context, uuid string, secretKey string) (err error) {
// 	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
// 	defer cancel()

// 	if err = utils.EncodeBase64(&uuid); err != nil {
// 		return
// 	}

// 	if err = a.userRepo.SaveOtpSecret(ctx, uuid, secretKey); err != nil {
// 		return err
// 	}

// 	return nil
// }
