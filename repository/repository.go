package repository

import (
	"context"

	"github.com/raulcv/goapiws/models"
)

type Repository interface {
	AddUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	Close() error
}

var implementation Repository

func SetRepository(repository Repository) { //Dependency Injection
	implementation = repository
}

func AddUser(ctx context.Context, user *models.User) error {
	return implementation.AddUser(ctx, user)
}
func GetUserById(ctx context.Context, id string) (*models.User, error) {
	return implementation.GetUserById(ctx, id)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx, email)
}

func Close() error {
	return implementation.Close()
}
