package logout

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(logoutURL string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		redirectURL := logoutURL
		ctx.Redirect(http.StatusFound, redirectURL)
	}
}
