package playerui

import (
	"ct-padel-s/src/features/padel/player/playermodel"
	"ct-padel-s/src/shared/utils"
	_ "embed"
	"html/template"
)

//go:embed playerchip.html
var playerchipHTML string
var PlayerchipComponent = utils.NewComponent("playerchip.html", playerchipHTML)

func RenderPlayerChip(player playermodel.Player) (template.HTML, error) {
	return PlayerchipComponent.Render(player)
}

func RenderPlayerWithTeamChip(player playermodel.Player, team int) (template.HTML, error) {
	return PlayerchipComponent.Render(player)
}
