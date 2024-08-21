package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/k0msak007/go-fiber-ecommerce/module/middlewares/middlewaresHandlers"
	"github.com/k0msak007/go-fiber-ecommerce/module/middlewares/middlewaresRepositories"
	"github.com/k0msak007/go-fiber-ecommerce/module/middlewares/middlewaresUsecases"
	"github.com/k0msak007/go-fiber-ecommerce/module/monitor/monitorHandlers"
	"github.com/k0msak007/go-fiber-ecommerce/module/users/usersHandlers"
	"github.com/k0msak007/go-fiber-ecommerce/module/users/usersRepositories"
	"github.com/k0msak007/go-fiber-ecommerce/module/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
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

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)
	router.Post("/signout", handler.SignOut)
}
