package logout

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(logoutURL string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		returnURL := ctx.Query("url")

		if returnURL == "" {
			returnURL = logoutURL
		}

		// redirectURL := logoutURL
		ctx.Redirect(http.StatusFound, returnURL)
	}
}
