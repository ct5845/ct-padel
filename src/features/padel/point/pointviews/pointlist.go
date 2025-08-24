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

//go:embed pointlist.html
var pointlistHTML string
var pointlistComponent = utils.NewComponent("pointlist.html", pointlistHTML)

func RenderPointList(points []*pointmodel.Point, game *gamemodel.Game, set *setmodel.Set, match *matchmodel.MatchWithPlayers) (template.HTML, error) {
	return pointlistComponent.Render(map[string]any{"Points": points, "Game": game, "Set": set, "Match": match})
}
