package main

import (
	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := app.oauth2Server.NewResponse()
		defer resp.Close()

		if ir := app.oauth2Server.HandleInfoRequest(resp, c.Request); ir != nil {
			c.Set("token", ir.AccessData)
		}

		if resp.IsError {
			app.responseHandler.Unauthorised(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
