package matchviews

import (
	"ct-padel-s/src/features/padel/match/matchmodel"
	"ct-padel-s/src/shared/utils"
	_ "embed"
	"html/template"
)

//go:embed getall.html
var getAllHTML string
var getAllComponent = utils.NewComponent("getall.html", getAllHTML)

type getAllViewModel struct {
	Matches []matchmodel.MatchWithPlayers
}

func RenderGetAll(matches []matchmodel.MatchWithPlayers) (template.HTML, error) {
	viewModel := getAllViewModel{Matches: matches}

	return getAllComponent.Render(viewModel)
}
