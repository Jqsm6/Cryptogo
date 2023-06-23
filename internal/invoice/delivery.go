package invoice

import "github.com/gin-gonic/gin"

type Handlers interface {
	Info() gin.HandlerFunc
	Create() gin.HandlerFunc
}
