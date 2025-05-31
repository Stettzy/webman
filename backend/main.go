package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"webman/pkg/database"
	"webman/pkg/models"
	"webman/pkg/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ProxyRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

type ProxyResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func main() {
	// Initialize database
	dbPath := "./data/webman.db"
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}
	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:9091", "http://127.0.0.1:9091", "http://webman.stettzy.com", "https://webman.stettzy.com"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Content-Length",
		"Accept",
		"Accept-Encoding",
		"Authorization",
		"Access-Control-Request-Headers",
		"Access-Control-Request-Method",
	}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Initialize services
	collectionService := services.NewCollectionService()
	headerService := services.NewHeaderService()

	// Collections endpoints
	r.GET("/collections", func(c *gin.Context) {
		collections, err := collectionService.ListCollections()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, collections)
	})

	r.POST("/collections", func(c *gin.Context) {
		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		collection, err := collectionService.CreateCollection(req.Name, req.Description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, collection)
	})

	r.GET("/collections/:id", func(c *gin.Context) {
		collection, err := collectionService.GetCollection(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, collection)
	})

	r.PUT("/collections/:id", func(c *gin.Context) {
		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		collection, err := collectionService.UpdateCollection(c.Param("id"), req.Name, req.Description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, collection)
	})

	r.DELETE("/collections/:id", func(c *gin.Context) {
		if err := collectionService.DeleteCollection(c.Param("id")); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})

	r.POST("/collections/:id/requests", func(c *gin.Context) {
		var request models.Request
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := collectionService.AddRequest(c.Param("id"), request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusCreated)
	})

	r.PUT("/collections/:id/requests/:requestId", func(c *gin.Context) {
		var request models.Request
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		request.ID = c.Param("requestId")
		err := collectionService.UpdateRequest(c.Param("id"), request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	})

	r.DELETE("/collections/:id/requests/:requestId", func(c *gin.Context) {
		err := collectionService.DeleteRequest(c.Param("id"), c.Param("requestId"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})

	// Default headers endpoint
	r.GET("/headers/default", func(c *gin.Context) {
		headers := headerService.GetDefaultHeaders()
		c.JSON(http.StatusOK, headers)
	})

	// Proxy endpoint
	r.POST("/", func(c *gin.Context) {
		var req ProxyRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "fail", "error": err.Error()})
			return
		}

		proxyReq, err := http.NewRequest(req.Method, req.URL, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "fail", "error": err.Error()})
			return
		}

		// Set headers
		for key, value := range req.Headers {
			proxyReq.Header.Set(key, value)
		}

		// Set body if present
		if len(req.Body) > 0 {
			proxyReq.Body = io.NopCloser(bytes.NewReader(req.Body))
			proxyReq.ContentLength = int64(len(req.Body))
		}

		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "fail", "error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "fail", "error": err.Error()})
			return
		}

		// Convert headers
		headers := make(map[string]string)
		for key, values := range resp.Header {
			if len(values) > 0 {
				headers[key] = values[0]
			}
		}

		// Try to parse as JSON first
		var jsonBody interface{}
		if err := json.Unmarshal(body, &jsonBody); err == nil {
			// If it's valid JSON, return it directly
			c.JSON(http.StatusOK, ProxyResponse{
				StatusCode: resp.StatusCode,
				Headers:    headers,
				Body:       string(body),
			})
			return
		}

		// If not JSON, encode as base64
		encodedBody := base64.StdEncoding.EncodeToString(body)
		c.JSON(http.StatusOK, ProxyResponse{
			StatusCode: resp.StatusCode,
			Headers:    headers,
			Body:       encodedBody,
		})
	})

	if err := r.Run(":9090"); err != nil {
		log.Fatal(err)
	}
}
