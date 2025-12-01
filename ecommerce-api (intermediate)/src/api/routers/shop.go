package routers

import (
	handlers "github.com/alielmi98/go-ecommerce-api/api/handlers"
	"github.com/alielmi98/go-ecommerce-api/api/middlewares"
	"github.com/alielmi98/go-ecommerce-api/config"
	"github.com/gin-gonic/gin"
)

// Category routes - manage product categories
func Category(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewCategoryHandler(cfg)

	r.POST("/", middlewares.Authentication(cfg), middlewares.Authorization([]string{"admin"}), h.Create)
	r.PUT("/:id", middlewares.Authentication(cfg), middlewares.Authorization([]string{"admin"}), h.Update)
	r.DELETE("/:id", middlewares.Authentication(cfg), middlewares.Authorization([]string{"admin"}), h.Delete)
	r.GET("/:id", h.GetById)
	r.POST("/get-by-filter", h.GetByFilter)
}

// Product routes - manage products
func Product(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewProductHandler(cfg)

	r.POST("/", middlewares.Authentication(cfg), middlewares.Authorization([]string{"admin"}), h.Create)
	r.PUT("/:id", middlewares.Authentication(cfg), middlewares.Authorization([]string{"admin"}), h.Update)
	r.DELETE("/:id", middlewares.Authentication(cfg), middlewares.Authorization([]string{"admin"}), h.Delete)
	r.GET("/:id", h.GetById)
	r.POST("/get-by-filter", h.GetByFilter)
}
