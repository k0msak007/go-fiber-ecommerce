package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/k0msak007/go-fiber-ecommerce/module/appinfo/appinfoHandlers"
	"github.com/k0msak007/go-fiber-ecommerce/module/appinfo/appinfoRepositories"
	"github.com/k0msak007/go-fiber-ecommerce/module/appinfo/appinfoUsecases"
	"github.com/k0msak007/go-fiber-ecommerce/module/files/filesHandlers"
	"github.com/k0msak007/go-fiber-ecommerce/module/files/filesUsecases"
	"github.com/k0msak007/go-fiber-ecommerce/module/middlewares/middlewaresHandlers"
	"github.com/k0msak007/go-fiber-ecommerce/module/middlewares/middlewaresRepositories"
	"github.com/k0msak007/go-fiber-ecommerce/module/middlewares/middlewaresUsecases"
	"github.com/k0msak007/go-fiber-ecommerce/module/monitor/monitorHandlers"
	"github.com/k0msak007/go-fiber-ecommerce/module/orders/ordersHandlers"
	"github.com/k0msak007/go-fiber-ecommerce/module/orders/ordersRepositories"
	"github.com/k0msak007/go-fiber-ecommerce/module/orders/ordersUsecases"
	"github.com/k0msak007/go-fiber-ecommerce/module/products/productsHandlers"
	"github.com/k0msak007/go-fiber-ecommerce/module/products/productsRepositories"
	"github.com/k0msak007/go-fiber-ecommerce/module/products/productsUsecases"
	"github.com/k0msak007/go-fiber-ecommerce/module/users/usersHandlers"
	"github.com/k0msak007/go-fiber-ecommerce/module/users/usersRepositories"
	"github.com/k0msak007/go-fiber-ecommerce/module/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule()
	ProductsModule()
	OrdersModule()
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

	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.mid.Authorize(2), handler.SignUpAdmin)

	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.s.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.s.cfg, usecase)

	router := m.r.Group("/appinfo")

	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddCategory)

	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)
	router.Get("/apiKey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)

	router.Delete("/:category_id/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.RemoveCategory)
}

func (m *moduleFactory) FilesModule() {
	usecase := filesUsecases.FilesUsecase(m.s.cfg)
	handler := filesHandlers.FilesHandler(m.s.cfg, usecase)

	router := m.r.Group("/files")

	router.Post("/upload", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UploadFile)
	router.Patch("/delete", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteFile)
}

func (m *moduleFactory) ProductsModule() {
	filesUsecase := filesUsecases.FilesUsecase(m.s.cfg)

	repository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, filesUsecase)
	usecase := productsUsecases.ProductsUsecase(repository)
	handler := productsHandlers.ProductsHandler(m.s.cfg, usecase, filesUsecase)

	router := m.r.Group("/products")

	router.Post("/", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddProduct)

	router.Patch("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UpdateProduct)

	router.Get("/", m.mid.ApiKeyAuth(), handler.FindProduct)
	router.Get("/:product_id", m.mid.ApiKeyAuth(), handler.FindOneProduct)

	router.Delete("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteProduct)
}

func (m *moduleFactory) OrdersModule() {
	filesUsecase := filesUsecases.FilesUsecase(m.s.cfg)
	productsRepository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, filesUsecase)

	ordersRepository := ordersRepositories.OrdersRepository(m.s.db)
	ordersUsecase := ordersUsecases.OrdersUsecase(ordersRepository, productsRepository)
	ordersHandler := ordersHandlers.OrdersHandler(m.s.cfg, ordersUsecase)

	router := m.r.Group("/orders")

	router.Post("/", m.mid.JwtAuth(), ordersHandler.InsertOrder)

	router.Get("/", m.mid.JwtAuth(), m.mid.Authorize(2), ordersHandler.FindOrder)
	router.Get("/:order_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), ordersHandler.FindOneOrder)
	router.Patch("/:user_id/:order_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), ordersHandler.UpdateOrder)
}
