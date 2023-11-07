package util

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetPage(t *testing.T) {

	// Create a Gin router instance for testing.
	router := gin.Default()

	// Define a route that uses the GetPage function.
	router.GET("/test", func(ctx *gin.Context) {
		offset, limit := GetPage(ctx)
		ctx.JSON(http.StatusOK, gin.H{"offset": offset, "limit": limit})
	})

	// Create a mock Gin context with the desired query parameters.
	req := httptest.NewRequest(http.MethodGet, "/test?page=2&pageSize=15", nil)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req

	// Serve the request using the router.
	router.ServeHTTP(rec, req)

	// Check the response status code.
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Parse the JSON response.
	responseData := struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	}{}

	// Unmarshal the response JSON.
	if err := json.Unmarshal(rec.Body.Bytes(), &responseData); err != nil {
		t.Errorf("Error parsing JSON response: %v", err)
	}

	// Check the offset and limit values in the response.
	expectedOffset := 15
	expectedLimit := 15

	if responseData.Offset != expectedOffset {
		t.Errorf("Expected offset %d, got %d", expectedOffset, responseData.Offset)
	}

	if responseData.Limit != expectedLimit {
		t.Errorf("Expected limit %d, got %d", expectedLimit, responseData.Limit)
	}

}
