package status

import "github.com/gin-gonic/gin"

type Handlers interface {
	GetAPIStatus() gin.HandlerFunc
}
