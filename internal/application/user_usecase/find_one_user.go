package user_usecase

import (
	"context"
	"errors"
	"fmt"

	"gocourse/common/logger"

	"gocourse/common/database"
	"gocourse/internal/domain/user_domain"
	"gocourse/internal/infrastructure/repository/user_repository"
)

type FindOneUserDeps struct {
	DB     database.DB
	Logger *logger.Logger
}

var (
	ErrFindOneUser = errors.New("find user failed")
)

func FindOneUser(c context.Context, deps FindOneUserDeps, id string) (user_domain.UserEntity, error) {
	user, findOneUserErr := user_repository.FindOne(c, deps.DB, id)
	if findOneUserErr != nil {
		deps.Logger.Error(fmt.Errorf("find one user failed: %w", findOneUserErr))
		return user_domain.UserEntity{}, ErrFindOneUser
	}
	return user, nil
}
