package interfaces

import "github.com/gin-gonic/gin"

type CompanyHandler interface {
	Create(ctx *gin.Context)
}
