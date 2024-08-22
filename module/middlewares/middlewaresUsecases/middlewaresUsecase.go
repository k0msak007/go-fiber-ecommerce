package middlewaresUsecases

import (
	"github.com/k0msak007/go-fiber-ecommerce/module/middlewares"
	"github.com/k0msak007/go-fiber-ecommerce/module/middlewares/middlewaresRepositories"
)

type IMiddlewareUsecase interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewareUsecase struct {
	middlewaresRepository middlewaresRepositories.IMiddlewareRepository
}

func MiddlewareUsecase(middlewaresRepository middlewaresRepositories.IMiddlewareRepository) IMiddlewareUsecase {
	return &middlewareUsecase{
		middlewaresRepository: middlewaresRepository,
	}
}

func (u *middlewareUsecase) FindAccessToken(userId, accessToken string) bool {
	return u.middlewaresRepository.FindAccessToken(userId, accessToken)
}

func (u *middlewareUsecase) FindRole() ([]*middlewares.Role, error) {
	roles, err := u.middlewaresRepository.FindRole()
	if err != nil {
		return nil, err
	}

	return roles, nil
}
