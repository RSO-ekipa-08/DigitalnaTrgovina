package api

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

// ProxyHandler posreduje zahteve na ustrezno mikrostoritev
func ProxyHandler(targetService string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Preberi telo zahteve
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Napaka pri branju zahteve"})
			return
		}

		// Uporabi url.Parse za pravilno rokovanje z URL-ji
		baseURL, err := url.Parse(os.Getenv(targetService))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Napaka pri razčlenjevanju URL-ja"})
			return
		}

		relPath, err := url.Parse(c.Param("path"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Napaka pri razčlenjevanju poti"})
			return
		}

		targetURL := baseURL.ResolveReference(relPath).String()

		req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(body))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Napaka pri ustvarjanju zahteve"})
			return
		}

		// Kopiraj vse headerje
		for name, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}

		// Izvedi zahtevo
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Napaka pri izvajanju zahteve"})
			return
		}
		defer resp.Body.Close()

		// Kopiraj odgovor nazaj klientu
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Napaka pri branju odgovora"})
			return
		}

		// Kopiraj headerje odgovora
		for name, values := range resp.Header {
			for _, value := range values {
				c.Header(name, value)
			}
		}

		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), responseBody)
	}
}
