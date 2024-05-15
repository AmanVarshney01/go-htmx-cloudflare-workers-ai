package main

import (
	"html/template"
	"io"
	"log"

	"github.com/amanvarshney01/go-htmx-cloudflare-workers-ai/handlers"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	e := echo.New()
	e.Static("/assets", "assets")
	e.Use(middleware.Logger())
	e.Renderer = newTemplates()

	e.GET("/", handlers.HandleIndex)

	e.POST("/prompt", handlers.HandlePrompt)

	e.Logger.Fatal(e.Start(":1323"))
}
