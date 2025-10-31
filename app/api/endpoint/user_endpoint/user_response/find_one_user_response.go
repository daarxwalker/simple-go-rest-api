package user_response

import (
	"gocourse/internal/domain/user_domain"
)

type FindOne struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func Map(user user_domain.UserEntity) FindOne {
	return FindOne{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
}
