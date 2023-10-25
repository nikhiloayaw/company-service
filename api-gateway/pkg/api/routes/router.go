package routes

import (
	_ "api-gateway/cmd/api/docs"
	handlerInterfaces "api-gateway/pkg/api/handler/interfaces"
	middlewareInterfaces "api-gateway/pkg/api/middleware/interfaces"
	"net/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func NewGinRouter(
	middleware middlewareInterfaces.Middleware,
	authHandler handlerInterfaces.AuthHandler,
	companyHandler handlerInterfaces.CompanyHandler,

) http.Handler {

	router := gin.New()

	router.Use(gin.Logger())

	// swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// group with api version
	api := router.Group("/api/v1")

	auth := api.Group("auth")
	{ // for all routes under auth

		RegisterAuthRoutes(auth, authHandler)
	}

	// from here onward all the api should be under the authenticate middleware
	api.Use(middleware.Authenticate("user"))

	// register company routes
	company := api.Group("/company")
	{
		RegisterCompanyRoutes(company, companyHandler)
	}

	return router
}
