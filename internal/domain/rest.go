package domain

import "github.com/gin-gonic/gin"

type RestHandler interface {
	Register(router *gin.Engine)
}
