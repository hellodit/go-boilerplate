package postgresql

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go-boilerplate/domain"
)

type psqlArticleRepository struct {
	DB *pg.DB
}

func NewPsqlArticleRepository(db *pg.DB) domain.ArticleRepository{
	return psqlArticleRepository{DB: db}
}

func (p psqlArticleRepository) Create(ctx context.Context, ar *domain.Article) error {
	_, err := p.DB.Model(ar).Insert()
	if err != nil {
		logrus.Warnln(err)
		return err
	}
	return nil
}

func (p psqlArticleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	article := new(domain.Article)
	_, err := p.DB.Model(article).Where("id=?", id).Delete()
	if err != nil {
		logrus.Warnln(err)
		return err
	}
	return nil
}

func (p psqlArticleRepository) FindBy(ctx context.Context, key, value string) (ar *domain.Article, err error) {
	ar = new(domain.Article)
	if err := p.DB.Model(ar).Where(key+"=?", value).First(); err != nil {
		return nil, err
		logrus.Warnln(err)
	}
	return ar, nil
}

func (p psqlArticleRepository) Update(ctx context.Context, id uuid.UUID, art *domain.Article) (ar *domain.Article, err error) {
	_, err = p.DB.Model(art).Where("id = ?", id).UpdateNotZero()
	if err != nil {
		logrus.Warnln(err)
		return nil, err
	}
	return art, nil
}

