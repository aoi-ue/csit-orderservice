package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ToyProductionKeyRequest struct {
	ToyProductionKey string `json:"toyProductionKey"`
}

type GatekeeperAccessRequest struct {
	OrderServiceHostOrIpAddress string `json:"orderServiceHostOrIpAddress"`
	SecretInput                 string `json:"secretInput"`
}

func main() {
	r := gin.Default()

	r.POST("/api/toyProductionKey", handleToyProductionKey)
	r.POST("/api/gatekeeper/access", handleGatekeeperAccess)

	r.Run(":8080")
}

func handleToyProductionKey(c *gin.Context) {
	var req ToyProductionKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate the toy production key format
	if !validateToyProductionKey(req.ToyProductionKey) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid toy production key format"})
		return
	}

	// Process the toy production key
	// ...

	c.Status(http.StatusOK)
}

func handleGatekeeperAccess(c *gin.Context) {
	var req GatekeeperAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate the secret input format
	if !validateSecretInput(req.SecretInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid secret input format"})
		return
	}

	// Process the order service host or IP address and secret input
	// ...

	c.Status(http.StatusOK)
}

func validateToyProductionKey(key string) bool {
	// Validate the toy production key format
	// The correct format is: <Most popular toy name (capitalize the first letter)>123! (no spaces)
	parts := strings.Split(key, "123!")
	if len(parts) != 2 || parts[1] != "" {
		return false
	}

	// Capitalize the first letter of the toy name
	toyName := strings.ToUpper(string(parts[0][0])) + parts[0][1:]
	return toyName == "Plush"
}

func validateSecretInput(input string) bool {
	// Validate the secret input format
	// The correct format is: <Most popular toy name (capitalize the first letter)>123! (no spaces)
	return input == "Plush123!"
}