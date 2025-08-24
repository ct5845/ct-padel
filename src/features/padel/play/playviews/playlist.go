package playviews

import (
	"ct-padel-s/src/features/padel/game/gamemodel"
	"ct-padel-s/src/features/padel/match/matchmodel"
	"ct-padel-s/src/features/padel/play/playmodel"
	"ct-padel-s/src/features/padel/point/pointmodel"
	"ct-padel-s/src/features/padel/set/setmodel"
	"ct-padel-s/src/shared/utils"
	_ "embed"
	"html/template"
)

//go:embed playlist.html
var playlistHTML string
var playlistComponent = utils.NewComponent("playlist.html", playlistHTML)

func RenderPlayList(plays []*playmodel.Play, point *pointmodel.Point, game *gamemodel.Game, set *setmodel.Set, match *matchmodel.MatchWithPlayers) (template.HTML, error) {
	return playlistComponent.Render(map[string]any{"Plays": plays, "Point": point, "Game": game, "Set": set, "Match": match})
}