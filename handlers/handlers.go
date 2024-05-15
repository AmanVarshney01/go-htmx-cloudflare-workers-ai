package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type AIResponse struct {
	Result struct {
		Response string `json:"response"`
	} `json:"result"`
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
}

func HandleIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func HandlePrompt(c echo.Context) error {
	prompt := c.FormValue("prompt")
	accountId := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/@cf/meta/llama-2-7b-chat-int8", accountId)
	requestData := map[string]string{
		"prompt": prompt,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		log.Println("Error marshaling request data:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending HTTP request:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	defer resp.Body.Close()

	var response AIResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("Error decoding response:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.Render(http.StatusOK, "response", response.Result.Response)
}
