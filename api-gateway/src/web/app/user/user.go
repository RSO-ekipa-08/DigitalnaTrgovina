package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(ctx *gin.Context) {
	// The JWT middleware will have already verified the token
	ctx.HTML(http.StatusOK, "user.html", nil)
}
