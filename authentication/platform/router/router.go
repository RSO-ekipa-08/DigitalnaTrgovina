package router

import (
	"01-Login/platform/authenticator"
	"01-Login/web/app/callback"
	"01-Login/web/app/home"
	"01-Login/web/app/login"
	"01-Login/web/app/logout"
	"01-Login/web/app/user"

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
