package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go-boilerplate/domain"
	"go-boilerplate/helper"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type userUsecase struct {
	UserRepo       domain.UserRepository
	ContextTimeout time.Duration
}

func (u *userUsecase) Fetch(ctx context.Context, limit, offset int) (res interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.ContextTimeout)
	defer cancel()
	users, err := u.UserRepo.Fetch(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"users":  users,
		"limit":  limit,
		"offset": offset,
		"row":    len(users),
	}, nil
}

func (u *userUsecase) Register(ctx context.Context, usr *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	usr.ID = uuid.New()
	usr.Password = string(hashedPassword)

	err = u.UserRepo.CreateUser(ctx, usr)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) Login(ctx context.Context, credential *domain.Credential) (res interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.ContextTimeout)
	defer cancel()

	user, err := u.UserRepo.Attempt(ctx, credential)
	if err != nil {
		return nil, errors.New("Email atau kata sandi tidak sesuai.\n Silakan tulis email terdaftar atau kata sandi yang sesuai.")
	}

	token, exp, err := helper.GenerateJwt(ctx, user)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token_type":   "Bearer",
		"access_token": token,
		"expires_in":   exp,
		"profile":      user,
	}, nil

}

func NewUserUsecase(userRepo domain.UserRepository, duration time.Duration) domain.UserUseCase {
	return &userUsecase{
		UserRepo:       userRepo,
		ContextTimeout: duration,
	}
}
