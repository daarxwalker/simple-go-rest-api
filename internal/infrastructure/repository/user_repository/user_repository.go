package user_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"gocourse/common/database"
	"gocourse/internal/domain/user_domain"
)

func FindOne(
	c context.Context, db database.DB, id string,
) (user_domain.UserEntity, error) {
	var user user_domain.UserEntity
	sql, args, createSqlErr := squirrel.Select().
		Columns(
			"id", "name", "email",
		).
		From(user_domain.Table).
		Where("id = ?", id).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if createSqlErr != nil {
		return user, fmt.Errorf("failed to create find one user sql: %w", createSqlErr)
	}
	if scanErr := pgxscan.Get(c, db, &user, sql, args...); scanErr != nil && !errors.Is(scanErr, pgx.ErrNoRows) {
		return user, fmt.Errorf("failed to scan user: %w", scanErr)
	}
	return user, nil
}

func FindOneByEmail(
	c context.Context, db database.DB, email string,
) (user_domain.UserEntity, error) {
	var user user_domain.UserEntity
	sql, args, createSqlErr := squirrel.Select().
		Columns(
			"id", "name", "email2",
		).
		From(user_domain.Table).
		Where("email = ?", email).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if createSqlErr != nil {
		return user, fmt.Errorf("failed to create find one user by email sql: %w", createSqlErr)
	}
	if scanErr := pgxscan.Get(c, db, &user, sql, args...); scanErr != nil && !errors.Is(scanErr, pgx.ErrNoRows) {
		return user, fmt.Errorf("failed to scan user: %w", scanErr)
	}
	return user, nil
}

func SaveOne(c context.Context, db database.DB, user user_domain.UserEntity) (string, error) {
	if len(user.Id) == 0 {
		return insertOne(c, db, user)
	}
	return updateOne(c, db, user)
}

func insertOne(c context.Context, db database.DB, user user_domain.UserEntity) (string, error) {
	var id string
	sql, args, createSqlErr := squirrel.Insert(user_domain.Table).
		Columns("name", "email").
		Values(
			user.Name, user.Email,
		).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if createSqlErr != nil {
		return "", fmt.Errorf("failed to create insert user sql: %w", createSqlErr)
	}
	if scanErr := db.QueryRow(c, sql, args...).Scan(&id); scanErr != nil {
		return "", fmt.Errorf("failed to insert user: %w", scanErr)
	}
	return id, nil
}

func updateOne(c context.Context, db database.DB, user user_domain.UserEntity) (string, error) {
	var id string
	sql, args, createSqlErr := squirrel.Update(user_domain.Table).
		Set("name", user.Name).
		Set("email", user.Email).
		Where("id = ?", user.Id).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if createSqlErr != nil {
		return "", fmt.Errorf("failed to create update user sql: %w", createSqlErr)
	}
	if scanErr := db.QueryRow(c, sql, args...).Scan(&id); scanErr != nil {
		return "", fmt.Errorf("failed to update user: %w", scanErr)
	}
	return id, nil
}
