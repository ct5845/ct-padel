package home

import (
	"ct-go-web-starter/src/shared/components/footer"
	"ct-go-web-starter/src/shared/components/header"
	"ct-go-web-starter/src/shared/templates"
	"ct-go-web-starter/src/shared/utils"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Rendering index page", "pages", "index", "path", r.URL.Path)
	if r.URL.Path != "/" {
		slog.Warn("Path not found", "path", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	header, err := header.Render(header.Data{Title: "CT Go Web Starter"})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	footer, err := footer.Render(footer.Data{})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	content, err := utils.LoadComponent("features/home/home.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	page, err := templates.Render(templates.Data{
		Title:       "CT Go Web Starter",
		HeaderHTML:  header,
		ContentHTML: template.HTML(content),
		FooterHTML:  footer,
	})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Info("Index page rendered successfully", "pages", "index")
	io.WriteString(w, string(page))
}
