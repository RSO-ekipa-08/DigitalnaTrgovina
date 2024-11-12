package callback

import (
	"encoding/json"
	"html/template"
	"net/http"

	"01-Login/platform/authenticator"

	"github.com/gin-gonic/gin"
)

func Handler(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := ctx.Query("code")
		if code == "" {
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		// Exchange code for token
		token, err := auth.Exchange(ctx.Request.Context(), code)
		if err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		// Verify token
		idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		// Get user profile
		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		// Convert profile to JSON string
		profileJSON, err := json.Marshal(profile)
		if err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		// Pass the data to template using template.JS for safe JavaScript execution
		ctx.HTML(http.StatusOK, "callback.html", gin.H{
			"access_token": template.JS(template.JSEscapeString(token.AccessToken)),
			"id_token":     template.JS(template.JSEscapeString(token.Extra("id_token").(string))),
			"profile":      string(profileJSON),
		})
	}
}
