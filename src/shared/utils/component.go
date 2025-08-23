package utils

import (
	"bytes"
	_ "ct-go-web-starter/src/infrastructure/config"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
)

// component represents a reusable template component
type component struct {
	name     string
	template *template.Template
}

// NewComponent creates a new component with the given name and HTML template string
func NewComponent(name, htmlTemplate string) *component {
	slog.Debug("Creating new component", "name", name)
	tmpl, err := template.New(name).Parse(htmlTemplate)
	if err != nil {
		slog.Error("Failed to parse template", "name", name, "error", err)
		panic(err)
	}

	slog.Debug("Component created successfully", "name", name)
	return &component{
		name:     name,
		template: tmpl,
	}
}

// Render executes the component template with the provided data and returns the HTML
func (c *component) Render(data interface{}) (template.HTML, error) {
	slog.Debug("Executing component template", "name", c.name)
	var buf bytes.Buffer
	err := c.template.Execute(&buf, data)
	if err != nil {
		slog.Error("Template execution failed", "name", c.name, "error", err)
		return "", err
	}
	slog.Debug("Component template executed successfully", "name", c.name)
	return template.HTML(buf.String()), nil
}

// MustRender executes the component template and panics on error (useful for compile-time safety)
func (c *component) MustRender(data interface{}) template.HTML {
	html, err := c.Render(data)
	if err != nil {
		slog.Error("Component render failed, panicking", "name", c.name, "error", err)
	}
	return html
}

// LoadComponent loads an HTML file from the src directory
func LoadComponent(relativeFilePath string) (string, error) {
	fullPath := filepath.Join("src", relativeFilePath)
	slog.Debug("Loading component file", "path", fullPath)
	
	content, err := os.ReadFile(fullPath)
	if err != nil {
		slog.Error("Failed to read component file", "path", fullPath, "error", err)
		return "", err
	}
	
	slog.Debug("Component file loaded successfully", "path", fullPath)
	return string(content), nil
}
