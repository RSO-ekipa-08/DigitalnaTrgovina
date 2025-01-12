package api

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker/v2"
)

// ProxyConfig vsebuje konfiguracijo za proxy handler
type ProxyConfig struct {
	CircuitBreaker *gobreaker.CircuitBreaker[*http.Response]
	Timeout        time.Duration
}

// DefaultProxyConfig vrne privzeto konfiguracijo za proxy
func DefaultProxyConfig(serviceName string) ProxyConfig {
	return ProxyConfig{
		CircuitBreaker: NewCircuitBreaker[*http.Response](serviceName, DefaultConfig()),
		Timeout:        5 * time.Second,
	}
}

// ProxyHandler posreduje zahteve na ustrezno mikrostoritev
func ProxyHandler(targetService string, config ProxyConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := config.CircuitBreaker.Execute(func() (*http.Response, error) {
			// Preberi telo zahteve
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				return nil, err
			}

			// Uporabi url.Parse za pravilno rokovanje z URL-ji
			baseURL, err := url.Parse(os.Getenv(targetService))
			if err != nil {
				return nil, err
			}

			relPath, err := url.Parse(c.Param("path"))
			if err != nil {
				return nil, err
			}

			targetURL := baseURL.ResolveReference(relPath).String()

			log.Printf("Making '%s' request to '%s'", c.Request.Method, targetURL)

			req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(body))
			if err != nil {
				return nil, err
			}

			// Kopiraj vse headerje
			for name, values := range c.Request.Header {
				for _, value := range values {
					req.Header.Add(name, value)
				}
			}

			// Izvedi zahtevo z timeout-om
			client := &http.Client{
				Timeout: config.Timeout,
			}
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}

			return resp, nil
		})

		if err != nil {
			if err == gobreaker.ErrOpenState {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Storitev trenutno ni na voljo"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Napaka pri izvajanju zahteve"})
			}
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
