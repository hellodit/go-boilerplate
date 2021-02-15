package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type (
	Article struct {
		tableName 	struct{} 	`pg:"articles"`
		ID        	uuid.UUID `pg:"id,pk,type:uuid" json:"id"`
		Title      	string   `pg:"title,type:varchar(255)" json:"name" form:"title"`
		Slug     	string  `pg:"slug,type:varchar(255)" json:"slug"`
		Description string 	`pg:"description" json:"description"`
		CreatedAt time.Time `pg:"created_at" json:"createdAt"`
		UpdatedAt time.Time `pg:"updated_at" json:"updatedAt"`

	}

	ArticleRepository interface {
		Create(ctx context.Context, ar *Article) error
		Delete(ctx context.Context, id uuid.UUID) error
		FindBy(ctx context.Context, key, value string) (ar *Article, err error)
		Update(ctx context.Context, id uuid.UUID, art *Article) (ar *Article, err error)
	}

	ArticleUsecase interface {
		CreateArticle(ctx context.Context, article *Article) error
		UpdateArticle(ctx context.Context, id uuid.UUID, article *Article) (res interface{}, err error)
		DeleteArticle(ctx context.Context, id uuid.UUID) error
		GetArticleBySlug(ctx context.Context, id string) (res interface{}, err error)
	}
)


