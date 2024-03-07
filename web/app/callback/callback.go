package callback

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"01-Login/platform/authenticator"
)

func filterMap(inputMap map[string]interface{}, fields []string) map[string]interface{} {
	// Create a new map to store the filtered fields
	filteredMap := make(map[string]interface{})

	// Iterate over the list of fields to keep
	for _, field := range fields {
		if value, exists := inputMap[field]; exists {
			// If the field exists in the input map, add it to the filtered map
			filteredMap[field] = value
		}
	}

	return filteredMap
}

// Handler for our callback.
func Handler(auth *authenticator.Authenticator, loginURL string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if ctx.Query("state") != session.Get("state") {
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		token, err := auth.Exchange(ctx.Request.Context(), ctx.Query("code"))
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Failed to convert an authorization code into a token.")
			return
		}

		idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Println(profile)

		fieldsToKeep := []string{"name", "nickname", "picture", "sub"}

		profile = filterMap(profile, fieldsToKeep)

		profileJSON, err := json.Marshal(profile)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		encodedProfile := url.QueryEscape(string(profileJSON))
		queryString := "data=" + encodedProfile

		returnURL := session.Get("url").(string)
		if returnURL == "" {
			returnURL = loginURL
		}

		redirectURL := returnURL + "?" + queryString
		ctx.Redirect(http.StatusFound, redirectURL)

	}
}
