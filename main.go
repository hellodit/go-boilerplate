package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-boilerplate/db/postgresql"
	MiddlewareCustom "go-boilerplate/middleware"
	"net/http"
	"os"
	"os/signal"
	"time"

	_articleHttpDelivery "go-boilerplate/article/delivery/http"
	_articlePostgreRepository "go-boilerplate/article/repository/postgresql"
	_articleUsecase "go-boilerplate/article/usecase"
	_userHttDelivery "go-boilerplate/user/delivery/http"
	_userPostgreRepository "go-boilerplate/user/repository/postgresql"
	_userUsecase "go-boilerplate/user/usecase"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		viper.AutomaticEnv()
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func main() {
	server := &http.Server{
		Addr:         ":" + viper.GetString("app_port"),
		ReadTimeout:  time.Duration(viper.GetInt("READ_TIMEOUT")) * time.Second,
		WriteTimeout: time.Duration(viper.GetInt("WRITE_TIMEOUT")) * time.Second,
	}

	postgreSQL := postgresql.Connect()

	timeoutCtx := time.Duration(viper.GetInt("CTX_TIMEOUT")) * time.Second

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	CustomMiddleware := MiddlewareCustom.Init()
	e.HTTPErrorHandler = CustomMiddleware.ErrorHandler
	MiddlewareCustom.Logger = logrus.New()
	e.Logger = MiddlewareCustom.GetEchoLogger()
	e.Use(MiddlewareCustom.Hook())

	go func() {
		if err := e.StartServer(server); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server up!")
	})

	userRepo := _userPostgreRepository.NewPsqlUserRepository(postgreSQL)
	userUsecase := _userUsecase.NewUserUsecase(userRepo, timeoutCtx)
	_userHttDelivery.NewUserHandler(e, userUsecase)

	articleRepo := _articlePostgreRepository.NewPsqlArticleRepository(postgreSQL)
	articleUsecase := _articleUsecase.NewArticleUsecase(articleRepo, timeoutCtx)
	_articleHttpDelivery.NewArticleHandler(e, articleUsecase)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
