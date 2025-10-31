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

type SaveOneDeps struct {
	DB     database.DB
	Logger *logger.Logger
}

func SaveOne(c context.Context, deps SaveOneDeps, user user_domain.UserEntity) (string, error) {
	existingUser, findUserErr := user_repository.FindOneByEmail(
		c,
		deps.DB,
		user.Email,
	)
	if findUserErr != nil {
		deps.Logger.Error(fmt.Errorf("find one user failed: %w", findUserErr))
		return "", errors.New("find user failed")
	}
	if len(existingUser.Id) > 0 {
		return existingUser.Id, errors.New("user already exists")
	}
	userId, saveOneUserErr := user_repository.SaveOne(
		c,
		deps.DB,
		user,
	)
	if saveOneUserErr != nil {
		deps.Logger.Error(fmt.Errorf("save one user failed: %w", saveOneUserErr))
		return "", errors.New("save user failed")
	}
	return userId, nil
}
