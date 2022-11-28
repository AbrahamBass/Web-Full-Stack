package repository

import (
	"context"

	"githuh.com/go/rest-crud/models"
)

type Repository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	InsertPost(ctx context.Context, post *models.Post) error
	GetPostById(ctx context.Context, id string) (*models.Post, error)
	UpdatePost(ctx context.Context, post *models.Post) error
	DeletePost(ctx context.Context, id string, userId string) error
	ListPost(ctx context.Context, page uint64) ([]*models.Post, error)
	Close() error
}

var implementations Repository

func SetRepository(repository Repository) {
	implementations = repository
}

func InsertUser(ctx context.Context, user *models.User) error {
	return implementations.InsertUser(ctx, user)
}

func GetUserById(ctx context.Context, id string) (*models.User, error) {
	return implementations.GetUserById(ctx, id)
}

func Close() error {
	return implementations.Close()
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementations.GetUserByEmail(ctx, email)
}

func InsertPost(ctx context.Context, post *models.Post) error {
	return implementations.InsertPost(ctx, post)
}

func GetPostById(ctx context.Context, id string) (*models.Post, error) {
	return implementations.GetPostById(ctx, id)
}

func UpdatePost(ctx context.Context, post *models.Post) error {
	return implementations.UpdatePost(ctx, post)
}

func DeletePost(ctx context.Context, id string, userId string) error {
	return implementations.DeletePost(ctx, id, userId)
}

func ListPost(ctx context.Context, page uint64) ([]*models.Post, error) {
	return implementations.ListPost(ctx, page)
}
