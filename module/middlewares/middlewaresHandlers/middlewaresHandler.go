package middlewaresHandlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/k0msak007/go-fiber-ecommerce/config"
	"github.com/k0msak007/go-fiber-ecommerce/module/entities"
	"github.com/k0msak007/go-fiber-ecommerce/module/middlewares/middlewaresUsecases"
	"github.com/k0msak007/go-fiber-ecommerce/pkg/kawaiiauth"
	"github.com/k0msak007/go-fiber-ecommerce/pkg/utils"
)

type middlewareHandlersErrCode string

const (
	routerCheckErr middlewareHandlersErrCode = "middleware-001"
	jwtAuthErr     middlewareHandlersErrCode = "middleware-002"
	paramsCheckErr middlewareHandlersErrCode = "middlware-003"
	authorizeErr   middlewareHandlersErrCode = "middlware-004"
	apiKeyErr      middlewareHandlersErrCode = "middlware-005"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorize(exprectRoleId ...int) fiber.Handler
	ApiKeyAuth() fiber.Handler
}

type middlewaresHandler struct {
	cfg                config.IConfig
	middlewaresUsecase middlewaresUsecases.IMiddlewareUsecase
}

func MiddlewaresHandler(cfg config.IConfig, middlewaresUsecase middlewaresUsecases.IMiddlewareUsecase) IMiddlewaresHandler {
	return &middlewaresHandler{
		cfg:                cfg,
		middlewaresUsecase: middlewaresUsecase,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// check router
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := kawaiiauth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		if !h.middlewaresUsecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"no permission to access",
			).Res()
		}

		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		return c.Next()
	}
}

func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Locals("userRoleId").(int) == 2 {
			return c.Next()
		}
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErr),
				"never gonna give you up",
			).Res()
		}
		return c.Next()
	}
}

func (h *middlewaresHandler) Authorize(exprectRoleId ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals("userRoleId").(int)
		if !ok {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizeErr),
				"user_id is not int type",
			).Res()
		}

		roles, err := h.middlewaresUsecase.FindRole()
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(authorizeErr),
				err.Error(),
			).Res()
		}

		sum := 0
		for _, v := range exprectRoleId {
			sum += v
		}

		expectedValueBinary := utils.BinaryConverter(sum, len(roles))
		userValueBinary := utils.BinaryConverter(userRoleId, len(roles))

		for i := range userValueBinary {
			if userValueBinary[i]&expectedValueBinary[i] == 1 {
				return c.Next()
			}
		}
		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizeErr),
			"no permission to access",
		).Res()
	}
}

func (h *middlewaresHandler) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-Api-Key")
		if _, err := kawaiiauth.ParseApiKey(h.cfg.Jwt(), key); err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(apiKeyErr),
				"apikey is invalid or required",
			).Res()
		}
		return c.Next()
	}
}
