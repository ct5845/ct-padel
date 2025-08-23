package pages

import (
	"ct-go-web-starter/src/organisms/footer"
	"ct-go-web-starter/src/organisms/header"
	"ct-go-web-starter/src/templates/page"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Rendering index page", "pages", "index", "path", r.URL.Path)
	if r.URL.Path != "/" {
		slog.Warn("Path not found", "path", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	header, err := header.Render(header.Data{Title: "CT Padel"})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	footer, err := footer.Render(footer.Data{})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	content := template.HTML("content")

	page, err := page.Render(page.Data{
		Title:       "CT Padel",
		HeaderHTML:  header,
		ContentHTML: content,
		FooterHTML:  footer,
	})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Info("Index page rendered successfully", "pages", "index")
	io.WriteString(w, string(page))
}
