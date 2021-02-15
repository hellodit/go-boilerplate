package postgresql

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go-boilerplate/domain"
	"golang.org/x/crypto/bcrypt"
)

type psqlUserRepository struct {
	DB *pg.DB
}


func (u *psqlUserRepository) CreateUser(ctx context.Context, usr *domain.User) error {
	_, err := u.DB.Model(usr).Insert()
	if err != nil {
		logrus.Warnln(err)
		return err
	}
	return nil
}

func (u *psqlUserRepository) Attempt(ctx context.Context, credential *domain.Credential) (user *domain.User, err error) {
	user = new(domain.User)
	err = u.DB.Model(user).Where("email = ?", credential.Email).Select()
	if err != nil {
		logrus.Warnln(err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credential.Password))
	if err != nil {
		logrus.Warnln(err)
		return nil, err
	}

	return user, nil
}

// Update Query for Update  user
func (u *psqlUserRepository) Update(ctx context.Context, usr *domain.User) error {
	_, err := u.DB.Model(usr).
		WherePK().
		UpdateNotZero()

	if err != nil {
		logrus.Warnln(err)
		return err
	}
	return nil
}

func (u *psqlUserRepository) Find(ctx context.Context, id uuid.UUID) (user *domain.User, err error) {
	user = new(domain.User)
	err = u.DB.Model(user).Where("id = ? ", id).First()
	if err != nil {
		logrus.Warnln(err)
		return nil, err
	}

	return user, nil
}

func (u *psqlUserRepository) FindBy(ctx context.Context, key, value string) (user *domain.User, err error) {
	user = new(domain.User)
	if err := u.DB.Model(user).Where(key+"=?", value).First(); err != nil {
		return nil, err
		logrus.Warnln(err)
	}
	return user, nil
}

func NewPsqlUserRepository(db *pg.DB) domain.UserRepository {
	return &psqlUserRepository{DB: db}
}
