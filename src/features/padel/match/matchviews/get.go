package matchviews

import (
	"ct-padel-s/src/features/padel/match/matchmodel"
	"ct-padel-s/src/features/padel/set/setviews"
	"ct-padel-s/src/features/padel/set/setmodel"
	"ct-padel-s/src/shared/utils"
	_ "embed"
	"html/template"
)

//go:embed get.html
var getHTML string
var getComponent = utils.NewComponent("get.html", getHTML)

func RenderGet(match *matchmodel.MatchWithPlayers, sets []*setmodel.Set) (template.HTML, error) {
	setsList, err := setviews.RenderSetList(match.ID, sets)

	if err != nil {
		return "", err
	}

	return getComponent.Render(map[string]any{
		"Match":    match,
		"SetsList": setsList,
	})
}
