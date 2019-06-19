package dashboard

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"waniKani/api"
	"waniKani/tools"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type SummaryData struct {
	Data struct {
		Reviews []Review `json:"reviews"`
	} `json:"data"`
}

type Review struct {
	AvailableAt string    `json:"available_at"`
	SubjectIds  []Subject `json:"subject_ids"`
}

type SubjectTypes struct {
	radical    int
	kanji      int
	vocabulary int
}

type Subject int
type BarData []float64
type BarLabels []string

func UpdateSummaryData() (BarLabels, BarData) {

	results := api.GetWaniKaniData("summary")

	return summaryParse(results)
}

func Summary() *widgets.BarChart {

	results := api.GetWaniKaniData("summary")

	l, d := summaryParse(results)

	return createBarChart(l, d)

}

func summaryParse(results []byte) (BarLabels, BarData) {

	var data SummaryData
	var bdata = make(BarData, 0)
	var btime = make(BarLabels, 0)

	//Subjects := readSubjectCache()

	_ = json.Unmarshal(results, &data)

	for i := 0; i < len(data.Data.Reviews); i++ {

		if len(data.Data.Reviews[i].SubjectIds) > 0 {

			timestamp, err := time.Parse(time.RFC3339, data.Data.Reviews[i].AvailableAt)
			reviewCount := len(data.Data.Reviews[i].SubjectIds)

			if err != nil {
				fmt.Printf("Error parsing timestamp: %v\n", err)
			}

			timestampString := timestamp.Local().Format("03:04 PM")

			if timestamp.Unix() <= time.Now().Unix() {
				timestampString = "Now"
			}

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			} else {
				bdata = append(bdata, float64(reviewCount))
				btime = append(btime, timestampString)
			}
		}
	}

	return btime, bdata

}

func createBarChart(labels BarLabels, data BarData) *widgets.BarChart {

	ws := tools.GetTTYSize()

	x2 := ws.Width
	y2 := ws.Height / 4

	bc := widgets.NewBarChart()
	bc.Data = data
	bc.Labels = labels
	bc.Title = "Upcoming Reviews"
	bc.SetRect(0, 5, x2, y2)
	bc.BarWidth = 10
	bc.BarColors = []ui.Color{api.ColorPink, api.ColorLightBlue}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	return bc

}
