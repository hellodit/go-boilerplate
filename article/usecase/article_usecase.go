package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go-boilerplate/domain"
	"strings"
	"time"
)

type articleUsecase struct {
	ArticleRepository domain.ArticleRepository
	ContextTimeout time.Duration
}

func (a articleUsecase) CreateArticle(ctx context.Context, article *domain.Article) error {
	slug := strings.ReplaceAll(article.Title," ", "-")
	article.ID = uuid.New()
	article.Slug = slug

	err := a.ArticleRepository.Create(ctx, article)
	if err != nil {
		logrus.Warningln(err)
		return err
	}

	return nil
}

func (a articleUsecase) UpdateArticle(ctx context.Context, id uuid.UUID, article *domain.Article) (res interface{}, err error) {
	slug := strings.ReplaceAll(article.Title," ", "-")
	article.Slug = slug

	article, err = a.ArticleRepository.Update(ctx, id, article)
	if err != nil {
		logrus.Warnln(err)
		return nil, err
	}

	return article, nil

}

func (a articleUsecase) DeleteArticle(ctx context.Context, id uuid.UUID) error {
	err := a.ArticleRepository.Delete(ctx, id)
	if err != nil {
		logrus.Warnln(err)
		return err
	}
	return nil
}

func (a articleUsecase) GetArticleBySlug(ctx context.Context, slug string) (res interface{}, err error) {
	art, err := a.ArticleRepository.FindBy(ctx, "slug", slug)

	if err != nil {
		logrus.Warnln(err)
		return nil, err
	}

	return art, nil

}

func NewArticleUsecase(repository domain.ArticleRepository, duration time.Duration) domain.ArticleUsecase{
	return &articleUsecase{
		ArticleRepository: repository,
		ContextTimeout:    duration,
	}
}