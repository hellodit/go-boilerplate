package http

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
	"go-boilerplate/domain"
	"go-boilerplate/middleware"
	"net/http"
	"time"
)

type articleHandler struct {
	articleUsecase domain.ArticleUsecase
}

func NewArticleHandler(e *echo.Echo, usecase domain.ArticleUsecase)  {
	handler := &articleHandler{articleUsecase: usecase}
	article := e.Group("/article")
	customMiddleware := middleware.Init()

	article.GET("/:slug", handler.GetArticleHandler)
	article.DELETE("/destroy", handler.DestroyArticleHandler, customMiddleware.Auth)
	article.POST("/store", handler.StoreArticleHandler, customMiddleware.Auth)
	article.PUT("/update", handler.UpdateArticleHandler, customMiddleware.Auth)
}

func (a articleHandler) GetArticleHandler(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	res, err := a.articleUsecase.GetArticleBySlug(ctx,e.Param("slug"))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return e.JSON(http.StatusOK, res)
}

func (a articleHandler) StoreArticleHandler(e echo.Context) error {
	rules := govalidator.MapData{
		"title":     []string{"required"},
		"description": []string{"required"},
	}

	validate := govalidator.Options{
		Request: e.Request(),
		Rules:   rules,
	}

	if err := govalidator.New(validate).Validate(); len(err) > 0 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err).SetInternal(errors.New("invalid parameter"))
	}

	ctx := e.Request().Context()

	var article domain.Article

	article.Title = e.FormValue("title")
	article.Description = e.FormValue("description")
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	if ctx == nil{
		ctx = context.Background()
	}

	err := a.articleUsecase.CreateArticle(ctx, &article)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
	})
}

func (a articleHandler) DestroyArticleHandler(e echo.Context) error {

	rules := govalidator.MapData{
		"id":     []string{"required"},
	}

	validate := govalidator.Options{
		Request: e.Request(),
		Rules:   rules,
	}

	if err := govalidator.New(validate).Validate(); len(err) > 0 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err).SetInternal(errors.New("invalid parameter"))
	}

	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	articleId, err := uuid.Parse(e.FormValue("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	err = a.articleUsecase.DeleteArticle(ctx, articleId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}
	return e.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
	})

}

func (a articleHandler) UpdateArticleHandler(e echo.Context) error {
	rules := govalidator.MapData{
		"title":     []string{"required"},
		"description": []string{"required"},
	}

	validate := govalidator.Options{
		Request: e.Request(),
		Rules:   rules,
	}

	if err := govalidator.New(validate).Validate(); len(err) > 0 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err).SetInternal(errors.New("invalid parameter"))
	}

	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	articleId, err := uuid.Parse(e.FormValue("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	var article domain.Article
	article.Title = e.FormValue("title")
	article.Description = e.FormValue("description")
	article.UpdatedAt = time.Now()


	res, err := a.articleUsecase.UpdateArticle(ctx, articleId, &article)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"data" 	: res,
		"status": "success",
	})
}
