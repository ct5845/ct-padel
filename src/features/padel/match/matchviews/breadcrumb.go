package matchviews

import (
	"ct-padel-s/src/features/padel/match/matchmodel"
	"ct-padel-s/src/shared/utils"
	_ "embed"
	"html/template"
)

//go:embed getbreadcrumb.html
var getBreadcrumbHTML string
var getBreadcrumbComponent = utils.NewComponent("getbreadcrumb.html", getBreadcrumbHTML)

//go:embed getallbreadcrumb.html
var getAllBreadcrumbHTML string
var getAllBreadcrumbComponent = utils.NewComponent("getallbreadcrumb.html", getAllBreadcrumbHTML)

func RenderGetBreadcrumb(match *matchmodel.MatchWithPlayers) (template.HTML, error) {
	return getBreadcrumbComponent.Render(map[string]any{"Match": match})
}

func RenderGetAllBreadcrumb() (template.HTML, error) {
	return getAllBreadcrumbComponent.Render(map[string]any{})
}
