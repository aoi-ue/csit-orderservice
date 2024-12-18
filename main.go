package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
        OrderServiceHostOrIpAddress string `json:"orderServiceHostOrIpAddress"`
        SecretInput                 string `json:"secretInput"`
}

type ToyProductionKeyRequest struct {
        ToyProductionKey string `json:"toyProductionKey"`
		Secret string `json:"secret"`
}

type ToyProductionKeyResponse struct {
	Key string `json:"key"`
}


func handleGatekeeperRequest(c *gin.Context) {
        var reqBody RequestBody
        if err := c.BindJSON(&reqBody); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        fmt.Println("Received request:")
        fmt.Printf("OrderServiceHostOrIpAddress: %s\n", reqBody.OrderServiceHostOrIpAddress)
        fmt.Printf("SecretInput: %s\n", reqBody.SecretInput)

        // Send request to CSIT Mini Challenge API
        client := &http.Client{}
        targetURL := "https://dec-2024-mini-challenge.csit-events.sg/api/gatekeeper/access"
        jsonData, err := json.Marshal(reqBody)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshalling request body"})
                return
        }

        req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(jsonData))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
                return
        }

        req.Header.Set("Content-Type", "application/json")

        // Send request and handle response (implement logic to handle success/failure)
        resp, err := client.Do(req)
        if err != nil {
                // Handle error sending request
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending request"})
                return
        }
        defer resp.Body.Close()

        // Process response from CSIT Mini Challenge API (based on challenge requirements)
        // Respond back to the original request (optional)
        c.JSON(http.StatusOK, gin.H{"message": "Request processed successfully"}) // Adjust based on response
}

func handleToyProductionKeyRequest(c *gin.Context) {
	var request ToyProductionKeyRequest
	if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
	}

	// Validate the secret 
	if request.Secret != "Plush123!" { 
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid secret"})
			return
	}

	// Assuming the key is received in the request body from Postman
	var keyResponse ToyProductionKeyResponse
	if err := c.BindJSON(&keyResponse); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key response"})
			return
	}

	c.JSON(http.StatusOK, keyResponse)
}
func main() {
        r := gin.Default()

        r.POST("/api/gatekeeper/access", handleGatekeeperRequest)
        r.POST("/api/toyProductionKey", handleToyProductionKeyRequest)

        port := os.Getenv("PORT")
        if port == "" {
                port = "8080"
        }

        r.Run(":" + port)
}