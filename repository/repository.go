package repository

import (
	"context"

	"github.com/raulcv/goapiws/models"
)

type Repository interface {
	AddUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	AddPost(ctx context.Context, user *models.Post) error
	// GetPostByUserId(ctx context.Context, id string, userId string) (int64, error)
	GetPostById(ctx context.Context, id string) (*models.Post, error)
	UpdatePost(ctx context.Context, post *models.Post) error
	DeletePost(ctx context.Context, id string, userId string) (int64, error)
	ActivatePost(ctx context.Context, id string) error
	ListPost(ctx context.Context, page uint64) ([]*models.Post, error)
	ListPostTwo(ctx context.Context, page uint64) ([]*models.Post, int64, error)

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

func AddPost(ctx context.Context, post *models.Post) error {
	return implementation.AddPost(ctx, post)
}
func GetPostById(ctx context.Context, id string) (*models.Post, error) {
	return implementation.GetPostById(ctx, id)
}
func UpdatePost(ctx context.Context, post *models.Post) error {
	return implementation.UpdatePost(ctx, post)
}

//	func GetPostByUserId(ctx context.Context, id string, userId string) (int64, error) {
//		return implementation.GetPostByUserId(ctx, id, userId)
//	}
func DeletePost(ctx context.Context, id string, userId string) (int64, error) {
	return implementation.DeletePost(ctx, id, userId)
}
func ActivatePost(ctx context.Context, id string) error {
	return implementation.ActivatePost(ctx, id)
}

func ListPost(ctx context.Context, page uint64) ([]*models.Post, error) {
	return implementation.ListPost(ctx, page)
}

func ListPostTwo(ctx context.Context, page uint64) ([]*models.Post, int64, error) {
	return implementation.ListPostTwo(ctx, page)
}
func Close() error {
	return implementation.Close()
}
