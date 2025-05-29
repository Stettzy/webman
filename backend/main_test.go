package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func getCollections(c *gin.Context) {
	c.JSON(http.StatusOK, []map[string]interface{}{})
}

func createCollection(c *gin.Context) {
	var collection map[string]interface{}
	if err := c.BindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, collection)
}

func saveRequest(c *gin.Context) {
	c.Status(http.StatusCreated)
}

func proxyRequest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// Request represents an API request
type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body,omitempty"`
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Setup routes
	router.GET("/collections", getCollections)
	router.POST("/collections", createCollection)
	router.POST("/collections/:id/requests", saveRequest)
	router.POST("/", proxyRequest)

	return router
}

func TestCollectionsEndpoint(t *testing.T) {
	router := setupTestRouter()

	t.Run("GET /collections returns 200", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/collections", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})

	t.Run("POST /collections creates new collection", func(t *testing.T) {
		w := httptest.NewRecorder()
		collection := map[string]string{
			"name": "Test Collection",
		}
		body, _ := json.Marshal(collection)
		req, _ := http.NewRequest("POST", "/collections", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "id")
		assert.Equal(t, collection["name"], response["name"])
	})
}

func TestRequestEndpoint(t *testing.T) {
	router := setupTestRouter()

	t.Run("POST / proxies request", func(t *testing.T) {
		w := httptest.NewRecorder()
		request := Request{
			Method: "GET",
			URL:    "https://api.github.com",
			Headers: map[string]string{
				"Accept": "application/json",
			},
		}
		body, _ := json.Marshal(request)
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})
}

func TestSaveRequest(t *testing.T) {
	router := setupTestRouter()

	// First create a collection
	w := httptest.NewRecorder()
	collection := map[string]string{
		"name": "Test Collection",
	}
	body, _ := json.Marshal(collection)
	req, _ := http.NewRequest("POST", "/collections", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var collectionResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &collectionResponse)
	collectionID := collectionResponse["id"].(string)

	t.Run("POST /collections/:id/requests saves request", func(t *testing.T) {
		w := httptest.NewRecorder()
		request := map[string]interface{}{
			"name":   "Test Request",
			"method": "GET",
			"url":    "https://api.github.com",
			"headers": map[string]string{
				"Accept": "application/json",
			},
		}
		body, _ := json.Marshal(request)
		req, _ := http.NewRequest("POST", "/collections/"+collectionID+"/requests", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
	})
}
