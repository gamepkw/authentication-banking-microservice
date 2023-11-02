package service

import (
	"context"

	producer "github.com/gamepkw/authentication-banking-microservice/internal/kafka/producer"
	"github.com/spf13/viper"
)

func (auth *authService) SendOtp(c context.Context, tel string) error {
	topic := "sms"
	brokerAddress := viper.GetString("kafka.broker_address")
	ctx, cancel := context.WithTimeout(c, auth.contextTimeout)
	defer cancel()
	otp, err := auth.GenerateOtp(ctx, tel)
	if err != nil {
		return err
	}

	producer.RunKafkaProducer(brokerAddress, topic, otp)
	return nil
}
