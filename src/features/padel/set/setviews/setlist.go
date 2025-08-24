package setviews

import (
	"ct-padel-s/src/features/padel/set/setmodel"
	"ct-padel-s/src/shared/utils"
	_ "embed"
	"html/template"
)

//go:embed setlist.html
var setlistHTML string
var setlistComponent = utils.NewComponent("setlist.html", setlistHTML)

func RenderSetList(matchID int, sets []*setmodel.Set) (template.HTML, error) {
	return setlistComponent.Render(
		map[string]any{
			"MatchID": matchID,
			"Sets":    sets,
		})
}
