package router

import (
	"time"

	"github.com/AndrivA89/orders/internal/transport/http/handlers"
	"github.com/AndrivA89/orders/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type Router struct {
	userHandler    *handlers.UserHandler
	productHandler *handlers.ProductHandler
	orderHandler   *handlers.OrderHandler
	logger         *logrus.Logger
}

func NewRouter(
	userHandler *handlers.UserHandler,
	productHandler *handlers.ProductHandler,
	orderHandler *handlers.OrderHandler,
	logger *logrus.Logger,
) *Router {
	return &Router{
		userHandler:    userHandler,
		productHandler: productHandler,
		orderHandler:   orderHandler,
		logger:         logger,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(middleware.ErrorHandler(r.logger))
	router.Use(middleware.Logger(r.logger))
	router.Use(middleware.RequestLogger(r.logger))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "orders",
		})
	})

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			// Rate limiting для регистрации: 5 попыток в минуту с burst = 2
			users.POST("",
				middleware.RateLimitMiddleware(rate.Every(time.Minute/5), 2),
				r.userHandler.CreateUser)
			users.GET("/:id", r.userHandler.GetUser)
			users.GET("/:user_id/orders", r.orderHandler.GetOrdersByUser)
		}

		products := v1.Group("/products")
		{
			products.POST("", r.productHandler.CreateProduct)
			products.GET("", r.productHandler.GetProducts)
			products.GET("/:id", r.productHandler.GetProduct)
		}

		orders := v1.Group("/orders")
		{
			// Rate limiting для создания заказов: 10 попыток в минуту с burst = 3
			orders.POST("",
				middleware.RateLimitMiddleware(rate.Every(time.Minute/10), 3),
				r.orderHandler.CreateOrder)
			orders.GET("/:id", r.orderHandler.GetOrder)
			orders.PATCH("/:id/confirm", r.orderHandler.ConfirmOrder)
			orders.PATCH("/:id/cancel", r.orderHandler.CancelOrder)
		}
	}

	return router
}
