package http

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
	"go-boilerplate/domain"
	"go-boilerplate/helper"
	"go-boilerplate/middleware"
	"net/http"
	"time"
)

type userHandler struct {
	userUsecase domain.UserUseCase
}

func NewUserHandler(e *echo.Echo, UserUsecase domain.UserUseCase){
	handler :=  &userHandler{
		userUsecase: UserUsecase,
	}
	user := e.Group("/user")
	customMiddleware := middleware.Init()

	user.POST("/register", handler.RegisterHandler)
	user.POST("/login", handler.LoginHandler)
	user.GET("/profile", handler.ProfileHandler, customMiddleware.Auth)
}

func (u userHandler) RegisterHandler(e echo.Context) error {
	rules := govalidator.MapData{
		"name":     []string{"required"},
		"password": []string{"required"},
		"email":    []string{"required"},
	}

	validate := govalidator.Options{
		Request: e.Request(),
		Rules:   rules,
	}

	if err := govalidator.New(validate).Validate(); len(err) > 0 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err).SetInternal(errors.New("invalid parameter"))
	}

	ctx := e.Request().Context()

	var usr domain.User

	usr.Name = e.FormValue("name")
	usr.Email = e.FormValue("email")
	usr.CreatedAt = time.Now()
	usr.UpdatedAt = time.Now()
	usr.Password = e.FormValue("password")

	if ctx == nil {
		ctx = context.Background()
	}

	err := u.userUsecase.Register(ctx, &usr)

	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   usr,
	})
}

func (u userHandler) LoginHandler(e echo.Context) error {
	rules := govalidator.MapData{
		"password": []string{"required"},
		"email":    []string{"required"},
	}

	validate := govalidator.Options{
		Request: e.Request(),
		Rules:   rules,
	}

	if err := govalidator.New(validate).Validate(); len(err) > 0 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err).SetInternal(errors.New("invalid parameter"))
	}

	var credential domain.Credential

	if err := e.Bind(&credential); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(errors.New("invalid parameter"))
	}

	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := u.userUsecase.Login(ctx, &credential)

	if err != nil {
		return echo.NewHTTPError(http.StatusFailedDependency, err.Error()).SetInternal(err)
	}

	return e.JSON(http.StatusOK, res)
}

func (u userHandler) ProfileHandler(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	claims, err := helper.ParseToken(e)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err).SetInternal(err)
	}
	profile := claims["user"]
	return e.JSON(http.StatusOK, profile)
}