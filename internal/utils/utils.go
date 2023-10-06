package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/milanvthakor/task-manager-api/pkg/config"
)

// HandlerFuncWithApp is a type for handler functions that require the application object.
type HandlerFuncWithApp func(*gin.Context, *config.Application)

// InjectApp injects the application object into a handler function.
func InjectApp(app *config.Application, handler HandlerFuncWithApp) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler(ctx, app)
	}
}
