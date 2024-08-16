package middlewaresUsecases

import "github.com/k0msak007/go-fiber-ecommerce/module/middlewares/middlewaresRepositories"

type IMiddlewareUsecase interface {
}

type middlewareUsecase struct {
	middlewaresRepository middlewaresRepositories.IMiddlewareRepository
}

func MiddlewareUsecase(middlewaresRepository middlewaresRepositories.IMiddlewareRepository) IMiddlewareUsecase {
	return &middlewareUsecase{
		middlewaresRepository: middlewaresRepository,
	}
}
