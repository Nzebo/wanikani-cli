package dashboard

import (
	"encoding/json"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"math"
	"strconv"
	"time"
	"waniKani/api"
	"waniKani/tools"
)

type LevelInfo []float64
type LevelLabels []string

type LevelData struct {
	Levels []LevelObj `json:"data"`
}

type LevelObj struct {
	Id      string `json:"id"`
	Details struct {
		Level    int    `json:"level"`
		Unlocked string `json:"unlocked_at"`
		Started  string `json:"started_at"`
		Passed   string `json:"passed_at"`
	} `json:"data"`
}

func UpdateLevelData() (LevelLabels, LevelInfo) {
	results := api.GetWaniKaniData("level_progressions")

	var resetTimestamp = api.GetLatestReset()

	return levelParse(results, resetTimestamp)
}

func Level() *widgets.BarChart {

	return levelBarChart(UpdateLevelData())

}

func levelParse(results []byte, resetTimestamp time.Time) (LevelLabels, LevelInfo) {

	var data LevelData

	var ldata = make(LevelInfo, 0)
	var llabel = make(LevelLabels, 0)

	_ = json.Unmarshal(results, &data)

	for i := 0; i < len(data.Levels); i++ {
		level := data.Levels[i].Details.Level

		if data.Levels[i].Details.Started == "" {
			continue
		}

		levelStartTime, err := time.Parse(time.RFC3339, data.Levels[i].Details.Started)

		if err != nil {
			log.Fatal(err)
		}

		if levelStartTime.Unix() < resetTimestamp.Unix() {

			continue
		} else {

			if data.Levels[i].Details.Passed != "" {
				levelPassedTime, err := time.Parse(time.RFC3339, data.Levels[i].Details.Passed)

				if err != nil {
					log.Fatal(err)
				}

				levelDuration := math.Round((float64(levelPassedTime.Unix()-levelStartTime.Unix())/60/60/24)*100) / 100

				ldata = append(ldata, levelDuration)
				llabel = append(llabel, "L"+strconv.Itoa(level))

			} else {

				levelDuration := math.Round((float64(time.Now().Unix()-levelStartTime.Unix())/60/60/24)*100) / 100

				ldata = append(ldata, levelDuration)
				llabel = append(llabel, "L"+strconv.Itoa(level))

			}

		}

	}

	return llabel, ldata

}

func levelBarChart(labels LevelLabels, data LevelInfo) *widgets.BarChart {

	ws := tools.GetTTYSize()

	y1 := ws.Height / 4
	x2 := ws.Width
	y2 := ws.Height / 2

	bc := widgets.NewBarChart()
	bc.Data = data
	bc.Labels = labels
	bc.Title = "Time Per Level (in Days)"
	bc.SetRect(0, y1, x2, y2)
	bc.BarWidth = 8
	bc.BarColors = []ui.Color{api.ColorPink, api.ColorLightBlue}
	bc.LabelStyles = []ui.Style{ui.NewStyle(api.ColorWhite)}
	bc.NumStyles = []ui.Style{ui.NewStyle(0)}

	return bc

}
