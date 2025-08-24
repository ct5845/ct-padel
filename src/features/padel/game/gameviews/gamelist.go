package gameviews

import (
	"ct-padel-s/src/features/padel/game/gamemodel"
	"ct-padel-s/src/features/padel/match/matchmodel"
	"ct-padel-s/src/features/padel/set/setmodel"
	"ct-padel-s/src/shared/utils"
	_ "embed"
	"html/template"
)

//go:embed gamelist.html
var gamelistHTML string
var gamelistComponent = utils.NewComponent("gamelist.html", gamelistHTML)

func RenderGameList(games []*gamemodel.Game, set *setmodel.Set, match *matchmodel.MatchWithPlayers) (template.HTML, error) {
	return gamelistComponent.Render(map[string]any{"Games": games, "Set": set, "Match": match})
}