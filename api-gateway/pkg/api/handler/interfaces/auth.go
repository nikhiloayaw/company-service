package interfaces

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
}
