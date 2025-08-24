package footer

import (
	"ct-padel-s/src/shared/utils"
	_ "embed"
	"html/template"
	"log/slog"
)

//go:embed footer.html
var footerHTML string
var component = utils.NewComponent("footer.html", footerHTML)

type Data struct{}

func init() {
	slog.Debug("Footer component initialized", "component", "footer")
}

func Render(data Data) (template.HTML, error) {
	slog.Debug("Rendering footer component", "component", "footer")
	result, err := component.Render(data)
	if err != nil {
		slog.Error("Failed to render footer component", "error", err)
		return "", err
	}
	slog.Debug("Footer component rendered successfully", "component", "footer")
	return result, nil
}
