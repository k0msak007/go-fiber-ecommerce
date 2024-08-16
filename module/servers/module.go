package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/k0msak007/go-fiber-ecommerce/module/middlewares/middlewaresHandlers"
	"github.com/k0msak007/go-fiber-ecommerce/module/middlewares/middlewaresRepositories"
	"github.com/k0msak007/go-fiber-ecommerce/module/middlewares/middlewaresUsecases"
	"github.com/k0msak007/go-fiber-ecommerce/module/monitor/monitorHandlers"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewareUsecase(repository)

	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}
