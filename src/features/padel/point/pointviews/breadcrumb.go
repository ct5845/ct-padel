package pointviews

import (
	"ct-padel-s/src/features/padel/game/gamemodel"
	"ct-padel-s/src/features/padel/match/matchmodel"
	"ct-padel-s/src/features/padel/point/pointmodel"
	"ct-padel-s/src/features/padel/set/setmodel"
	"ct-padel-s/src/shared/utils"
	_ "embed"
	"html/template"
)

//go:embed breadcrumb.html
var breadcrumbHTML string
var breadcrumbComponent = utils.NewComponent("breadcrumb.html", breadcrumbHTML)

func RenderBreadcrumb(match *matchmodel.MatchWithPlayers, set *setmodel.Set, game *gamemodel.Game, point *pointmodel.Point) (template.HTML, error) {
	return breadcrumbComponent.Render(map[string]any{"Match": match, "Set": set, "Game": game, "Point": point})
}