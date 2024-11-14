package router

import (
	"authentication/src/platform/authenticator"
	"authentication/src/web/app/callback"
	"authentication/src/web/app/home"
	"authentication/src/web/app/login"
	"authentication/src/web/app/logout"
	"authentication/src/web/app/user"

	"github.com/gin-gonic/gin"
)

func New(auth *authenticator.Authenticator) *gin.Engine {
	router := gin.Default()

	router.Static("/public", "web/static")
	router.LoadHTMLGlob("web/template/*")

	// Public routes
	router.GET("/", home.Handler)
	router.GET("/login", login.Handler(auth))
	router.GET("/callback", callback.Handler(auth))
	router.GET("/user", user.Handler) // Move this out of API group
	router.GET("/logout", logout.Handler)

	return router
}
