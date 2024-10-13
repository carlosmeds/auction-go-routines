package user_entity

import (
	"auctions-go-routines/internal/internal_error"
	"context"
)

type User struct {
	Id   string
	Name string
}

type UserRepositoryInterface interface {
	FindUserById(
		ctx context.Context, userId string) (*User, *internal_error.InternalError)
}
