package interfaces

import "github.com/gin-gonic/gin"

type Middleware interface {
	Authenticate(role string) gin.HandlerFunc
}
