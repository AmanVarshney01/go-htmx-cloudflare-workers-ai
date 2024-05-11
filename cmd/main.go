package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func main() {
	e := echo.New()
	e.Static("/assets", "assets")
	e.Use(middleware.Logger())
	e.Renderer = newTemplates()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	// e.GET("/github", func(c echo.Context) error {
	// 	resp, err := http.Get("https://api.github.com/repos/amanvarshney01/gla-app-reimagined")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer resp.Body.Close()

	// 	var repo struct {
	// 		Name string `json:"name"`
	// 	}

	// 	err = json.NewDecoder(resp.Body).Decode(&repo)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return c.Render(http.StatusOK, "github", repo.Name)
	// })

	e.GET("/ai", func(c echo.Context) error {
		accountId := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
		apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")

		url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/@cf/meta/llama-2-7b-chat-int8", accountId)
		requestData := map[string]string{
			"prompt": "Cons of java",
		}
		requestBody, err := json.Marshal(requestData)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+apiToken)
		req.Header.Set("Content-Type", "application/json")

		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		type AIResponse struct {
			Result struct {
				Response string `json:"response"`
			} `json:"result"`
			Success  bool     `json:"success"`
			Errors   []string `json:"errors"`
			Messages []string `json:"messages"`
		}

		var response AIResponse

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "ai", response.Result.Response)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
