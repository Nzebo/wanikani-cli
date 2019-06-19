package dashboard

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/widgets"
	"waniKani/api"
)


func Progress() (*widgets.Gauge, *widgets.Gauge) {

	return LevelProgress(UpdateLevelProgress())
}

func LevelProgress(radicalPct, kanjiPct int) (*widgets.Gauge, *widgets.Gauge) {

	g0 := widgets.NewGauge()
	g0.Title = "Radicals"
	g0.Percent = radicalPct
	g0.BarColor = api.ColorLightBlue
	g0.TitleStyle.Fg = ui.ColorWhite
	g0.SetRect(0, 48, 75, 45)

	g1 := widgets.NewGauge()
	g1.Title = "Kanji"
	g1.Percent = kanjiPct
	g1.BarColor = api.ColorPink
	g1.TitleStyle.Fg = ui.ColorWhite
	g1.SetRect(0, 51, 75, 48)

	return g0, g1
}

func UpdateLevelProgress() (int, int) {
	radicalPct, kanjiPct := api.GetLevelPercentage()

	return radicalPct, kanjiPct
}