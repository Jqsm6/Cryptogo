package invoice

import "github.com/gin-gonic/gin"

type Handlers interface {
	Create() gin.HandlerFunc
	Info() gin.HandlerFunc
}
