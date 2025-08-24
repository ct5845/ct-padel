package setviews

import (
	"ct-padel-s/src/features/padel/game/gamemodel"
	"ct-padel-s/src/features/padel/match/matchmodel"
	"ct-padel-s/src/features/padel/set/setmodel"
	"ct-padel-s/src/shared/utils"
	_ "embed"
	"html/template"
)

//go:embed get.html
var getHTML string
var getComponent = utils.NewComponent("get.html", getHTML)

func RenderGet(set *setmodel.Set, match *matchmodel.MatchWithPlayers, games []*gamemodel.Game, gamesListHTML template.HTML) (template.HTML, error) {
	return getComponent.Render(map[string]any{
		"Set": set, 
		"Match": match, 
		"Games": games, 
		"GamesListHTML": gamesListHTML,
	})
}
