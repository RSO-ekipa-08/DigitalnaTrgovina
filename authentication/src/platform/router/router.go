package router

import (
	"authentication/src/platform/authenticator"
	"authentication/src/web/app/callback"
	"authentication/src/web/app/login"
	"authentication/src/web/app/logout"
	"authentication/src/web/app/user"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func New(auth *authenticator.Authenticator) *gin.Engine {
	router := gin.Default()

	router.Static("/public", "src/web/static")
	router.LoadHTMLGlob("src/web/template/*")

	// Public routes
	router.GET("/login", login.Handler(auth))
	router.GET("/callback", callback.Handler(auth))
	router.GET("/user", user.Handler) // Move this out of API group
	router.GET("/logout", logout.Handler)

	router.Use(func(c *gin.Context) {
        c.Set("AppServiceURL", os.Getenv("APP_SERVICE_URL"))
        c.Set("ReviewsServiceURL", os.Getenv("REVIEWS_SERVICE_URL"))
        c.Next()
    })

    // Update handlers to pass service URLs
    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "home.html", gin.H{
            "AppServiceURL":    os.Getenv("APP_SERVICE_URL"),
            "ReviewsServiceURL": os.Getenv("REVIEWS_SERVICE_URL"),
        })
    })

    router.GET("/app/:id", func(c *gin.Context) {
        c.HTML(http.StatusOK, "app.html", gin.H{
            "AppServiceURL":    os.Getenv("APP_SERVICE_URL"),
            "ReviewsServiceURL": os.Getenv("REVIEWS_SERVICE_URL"),
        })
    })

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	return router
}
