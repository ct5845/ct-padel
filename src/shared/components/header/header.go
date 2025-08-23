package header

import (
	_ "ct-go-web-starter/src/infrastructure/config"
	"ct-go-web-starter/src/shared/utils"
	_ "embed"
	"html/template"
	"log/slog"
)

//go:embed header.html
var headerHTML string
var component = utils.NewComponent("header.html", headerHTML)

type Data struct {
	Title string
}

func init() {
	slog.Debug("Header component initialized", "component", "header")
}

func Render(data Data) (template.HTML, error) {
	slog.Debug("Rendering header component", "component", "header")
	result, err := component.Render(data)
	if err != nil {
		slog.Error("Failed to render header component", "error", err)
		return "", err
	}
	slog.Debug("Header component rendered successfully", "component", "header")
	return result, nil
}
