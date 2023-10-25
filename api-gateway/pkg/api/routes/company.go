package routes

import (
	handlerInterfaces "api-gateway/pkg/api/handler/interfaces"

	"github.com/gin-gonic/gin"
)

func RegisterCompanyRoutes(company *gin.RouterGroup, companyHandler handlerInterfaces.CompanyHandler) {

	company.POST("", companyHandler.Create)
}
