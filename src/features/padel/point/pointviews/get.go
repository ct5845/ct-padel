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

//go:embed get.html
var getHTML string
var getComponent = utils.NewComponent("get.html", getHTML)

func RenderGet(point *pointmodel.Point, game *gamemodel.Game, set *setmodel.Set, match *matchmodel.MatchWithPlayers, playsListHTML template.HTML) (template.HTML, error) {
	return getComponent.Render(map[string]any{
		"Point": point, 
		"Game": game, 
		"Set": set, 
		"Match": match, 
		"PlaysListHTML": playsListHTML,
	})
}