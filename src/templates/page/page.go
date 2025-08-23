package page

import (
	_ "ct-go-web-starter/src/config"
	"ct-go-web-starter/src/utils"
	_ "embed"
	"html/template"
	"log/slog"
)

//go:embed page.html
var pageHTML string
var component = utils.NewComponent("page.html", pageHTML)

type Data struct {
	Title       string
	HeaderHTML  template.HTML
	ContentHTML template.HTML
	FooterHTML  template.HTML
}

func init() {
	slog.Debug("Page template initialized", "component", "page")
}

func Render(data Data) (template.HTML, error) {
	slog.Debug("Rendering page template", "component", "page", "title", data.Title)
	result, err := component.Render(data)
	if err != nil {
		slog.Error("Failed to render page template", "error", err, "title", data.Title)
		return "", err
	}
	slog.Debug("Page template rendered successfully", "component", "page", "title", data.Title)
	return result, nil
}