package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
	"ws-middle-cache/internal/middleware"
	"ws-middle-cache/internal/services"
	"ws-middle-cache/pkg/utils"

	"github.com/gin-gonic/gin"
)

var (
	// Initialize in-memory cache
	inMemoryCache *services.Cache

	// Azure Table Storage client (replace with your actual credentials)
	azureTableClient, _ = services.NewAzureTableClient(
		os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"),       // Replace with your account name
		os.Getenv("AZURE_STORAGE_ACCOUNT_KEY"),        // Replace with your account key
		os.Getenv("AZURE_STORAGE_ACCOUNT_TABLE_NAME"), // Replace with your table name
	)
)

func SetCacheInstance(cache *services.Cache) {
	logger := middleware.NewLogger()
	logger.Info("Setting in-memory cache instance")
	inMemoryCache = cache
}

// CacheHandler handles requests to fetch or store data in the cache.
func CacheHandler(c *gin.Context) {
	logger := middleware.NewLogger()
	requestedPath := c.Param("any")
	queryParams := c.Request.URL.Query()

	env := &utils.Environment{}
	expiration := 10 * 60                                       // Default expiration time in seconds (10 minutes)
	envExpiration := env.Get("CACHE_EXPIRATION_SECONDS", "600") // Default to 600 seconds if not set
	if parsedExpiration, err := strconv.Atoi(envExpiration); err == nil {
		expiration = parsedExpiration
	}

	// Sanitize function to remove unauthorized characters
	sanitize := func(input string) string {
		re := regexp.MustCompile(`[<>:"/\\|?*=]`)
		return re.ReplaceAllString(input, "_")
	}

	sanitizedPath := sanitize(requestedPath)
	sanitizedQuery := sanitize(queryParams.Encode())
	tokenizedFileName := sanitizedPath + "_" + sanitizedQuery

	// Check in-memory cache
	logger.Info("Checking in-memory cache for key:", tokenizedFileName)
	if value, found := inMemoryCache.Get(tokenizedFileName); found {
		logger.Info("Cache hit (in-memory) for key:", tokenizedFileName)
		c.Header("X-Cache-Status", "Hit")
		utils.JSONResponse(c, http.StatusOK, gin.H{
			"message":           "Cache hit (in-memory)",
			"tokenizedFileName": tokenizedFileName,
			"value":             json.RawMessage(value.(string)),
		})
		return
	}

	// Check Azure Table Storage
	logger.Info("Checking Azure Table Storage for key:", tokenizedFileName)
	entity, err := azureTableClient.GetEntity(context.Background(), sanitizedPath, tokenizedFileName)
	if err == nil {
		logger.Info("Cache hit (Azure Table Storage) for key:", tokenizedFileName)

		// Cache the result in memory
		inMemoryCache.Set(tokenizedFileName, entity.Properties["Value"].(string), time.Duration(expiration)*time.Second)

		c.Header("X-Cache-Status", "Table Hit")
		utils.JSONResponse(c, http.StatusOK, gin.H{
			"message":           "Cache hit (Azure Table Storage)",
			"tokenizedFileName": tokenizedFileName,
			"value":             json.RawMessage(entity.Properties["Value"].(string)),
		})
		return
	}

	// If not found, generate a new value (example: timestamp)
	cleanedPath := requestedPath
	if len(cleanedPath) > 0 && cleanedPath[0] == '/' {
		cleanedPath = cleanedPath[1:]
	}
	url := fmt.Sprintf("%s/%s?%s", env.Get("BACKEND_ENDPOINT", ""), cleanedPath, queryParams.Encode())
	logger.Info("Calling web service:", url)
	newValue, err := callWebService(url)

	if err != nil {
		logger.Error("Failed to call web service:", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to call web service")
		return
	}

	// Save to Azure Table Storage
	properties := map[string]any{
		"Value": func() string {
			jsonValue, err := json.Marshal(newValue)
			if err != nil {
				panic(fmt.Sprintf("Failed to marshal newValue: %v", err))
			}
			return string(jsonValue)
		}(),
	}
	err = azureTableClient.InsertEntity(context.Background(), sanitizedPath, tokenizedFileName, properties)
	if err != nil {
		logger.Error("Failed to save to Azure Table Storage:", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save to Azure Table Storage")
		return
	}

	// Save to in-memory cache
	logger.Info("Saving to in-memory cache for key:", tokenizedFileName)
	jsonValue, err := json.Marshal(newValue)
	if err != nil {
		logger.Error("Failed to convert value to JSON:", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to convert value to JSON")
		return
	}
	inMemoryCache.Set(tokenizedFileName, string(jsonValue), time.Duration(expiration)*time.Second)

	c.Header("X-Cache-Status", "Miss")
	utils.JSONResponse(c, http.StatusOK, gin.H{
		"message":           "Cache miss (generated new value)",
		"tokenizedFileName": tokenizedFileName,
		"value":             json.RawMessage(string(jsonValue)),
	})
}

// callWebService makes an HTTP GET request to the given URL and parses the JSON response.
func callWebService(url string) (interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call web service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("web service returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response body as JSON: %w", err)
	}

	return result, nil
}
