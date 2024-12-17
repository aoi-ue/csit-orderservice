package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

type GatekeeperResponse struct {
	Status  string      `json:"status"`  // e.g., "success" or "error"
	Message string      `json:"message"` // A message providing more details
	Data    interface{} `json:"data"`    // Any additional data returned by the service
}

func main() {
	r := gin.Default()

	// to receive 
	r.POST("/api/toyProductionKey", handleToyProductionKey)

	// to send 
	r.POST("/api/gatekeeper/access", handleGatekeeperAccess)
	
	// Use environment variable for port, default to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
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

	// Return the received toy production key in the response
	c.JSON(http.StatusOK, gin.H{"toyProductionKey": req.ToyProductionKey})
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

	// Send request to the Gatekeeper Service
	gatekeeperResponse, err := sendToGatekeeperService(req.OrderServiceHostOrIpAddress, req.SecretInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to access Gatekeeper Service"})
		return
	}

	// Process the response from the Gatekeeper Service
	c.JSON(http.StatusOK, gatekeeperResponse)
}

func sendToGatekeeperService(orderServiceHostOrIpAddress, secretInput string) (GatekeeperResponse, error) {
	gatekeeperURL := "https://dec-2024-mini-challenge.csit-events.sg/api/gatekeeper/access"

	// Create the request body
	requestBody, err := json.Marshal(GatekeeperAccessRequest{
		OrderServiceHostOrIpAddress: orderServiceHostOrIpAddress,
		SecretInput:                 secretInput,
	})
	if err != nil {
		return GatekeeperResponse{}, err
	}

	// Send the POST request to the Gatekeeper Service
	resp, err := http.Post(gatekeeperURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return GatekeeperResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GatekeeperResponse{}, fmt.Errorf("failed to access Gatekeeper Service, status code: %d", resp.StatusCode)
	}

	var gatekeeperResponse GatekeeperResponse
	if err := json.NewDecoder(resp.Body).Decode(&gatekeeperResponse); err != nil {
		return GatekeeperResponse{}, err
	}

	return gatekeeperResponse, nil
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