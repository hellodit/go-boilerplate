package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type (
	//Credential user
	Credential struct {
		Email    string `json:"email" form:"email" validate:"required,email"`
		Password string `json:"password" form:"password" validate:"required"`
	}

	//User struct
	User struct {
		tableName struct{}  `pg:"users"`
		ID        uuid.UUID `pg:"id,pk,type:uuid" json:"id"`
		Name      string    `pg:"name,type:varchar(255)" json:"name" form:"name" validate:"required"`
		Email     string    `pg:"email,type:varchar(255)" json:"email" form:"email" validate:"required,email"`
		Password  string    `pg:"password,type:varchar(255)" json:"-" form:"password" validate:"required"`
		CreatedAt time.Time `pg:"created_at" json:"createdAt"`
		UpdatedAt time.Time `pg:"updated_at" json:"updatedAt"`
	}
)

//UserRepository interface
type UserRepository interface {
	CreateUser(ctx context.Context, usr *User) error
	Attempt(ctx context.Context, credential *Credential) (user *User, err error)
	Update(ctx context.Context, usr *User) error
	Find(ctx context.Context, id uuid.UUID) (user *User, err error)
	FindBy(ctx context.Context, key, value string) (user *User, err error)
}

//UserUseCase interface
type UserUseCase interface {
	Register(ctx context.Context, usr *User) error
	Login(ctx context.Context, credential *Credential) (res interface{}, err error)
}
