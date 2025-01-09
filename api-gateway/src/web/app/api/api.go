package api

import (
	"bytes"
	"io"
	"net/http"
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

		// Ustvari novo zahtevo na ciljno storitev
		targetURL := os.Getenv(targetService) + c.Param("path")
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
