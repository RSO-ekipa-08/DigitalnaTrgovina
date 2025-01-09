package router

import (
	"authentication/src/platform/auth"
	"authentication/src/web/app/api"
	"authentication/src/web/app/user"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func New(authClient *auth.AuthClient) *gin.Engine {
	router := gin.Default()

	router.Static("/public", "src/web/static")
	router.LoadHTMLGlob("src/web/template/*")

	// Public routes
	router.GET("/login", func(c *gin.Context) {
		authURL, err := authClient.Login(c, c.Query("redirect_url"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, authURL)
	})

	router.GET("/callback", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		resp, err := authClient.Verify(c, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.HTML(http.StatusOK, "callback.html", gin.H{
			"access_token": resp.AccessToken,
			"id_token":     resp.IdToken,
			"profile":      resp.Profile,
		})
	})

	router.GET("/logout", func(c *gin.Context) {
		logoutURL, err := authClient.Logout(c, "/")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, logoutURL)
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"AppServiceURL":     os.Getenv("APP_SERVICE_URL"),
			"ReviewsServiceURL": os.Getenv("REVIEWS_SERVICE_URL"),
		})
	})

	router.GET("/app/:id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "app.html", gin.H{
			"AppServiceURL":     os.Getenv("APP_SERVICE_URL"),
			"ReviewsServiceURL": os.Getenv("REVIEWS_SERVICE_URL"),
		})
	})

	router.GET("/user", user.Handler)

	// API proxy routes
	apiGroup := router.Group("/api")
	{
		// App service proxy
		appGroup := apiGroup.Group("/app")
		appGroup.Any("/*path", api.ProxyHandler("APP_SERVICE_URL"))

		// Reviews service proxy
		reviewsGroup := apiGroup.Group("/reviews")
		reviewsGroup.Any("/*path", api.ProxyHandler("REVIEWS_SERVICE_URL"))
	}

	return router
}
