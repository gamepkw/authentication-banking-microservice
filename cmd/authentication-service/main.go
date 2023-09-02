package main

import (
	"log"
	"net/http"
	"time"

	_accountRepo "github.com/atm5_microservices/accounts_service/internal/repositories"

	_accountService "github.com/atm5_microservices/accounts_service/internal/services"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {

	e := echo.New()
	middL := _accountHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	e.Use(middL.RateLimitMiddleware)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3001", "http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	accountrepo := _accountRepo.NewMysqlAccountRepository(dbConn, redis)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	accountService := _accountService.NewAccountService(accountrepo, transactionRepo, redis, timeoutContext)

	_accountHttpDelivery.NewAccountHandler(e, accountService)

	log.Fatal(e.Start(viper.GetString("server.address"))) //nolint

}
