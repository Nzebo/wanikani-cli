package main

import (
	"fmt"
	"os"
	"time"
	"waniKani/api"
	"waniKani/dashboard"
	"waniKani/tools"
	ui "github.com/gizak/termui/v3"
	_ "github.com/gizak/termui/v3/widgets"
)

func displayDashboard() {
	if err := ui.Init(); err != nil {
		fmt.Printf("Failed to initialize termui: %v\n", err)
	}

	defer ui.Close()

	userText := dashboard.User()
	summaryBarChart := dashboard.Summary()
	levelsBarChart := dashboard.Level()
	radicalProgress, kanjiProgress := dashboard.Progress()

	draw := func() {
		userText.Title, userText.Text = dashboard.UpdateUserData()
		summaryBarChart.Labels, summaryBarChart.Data = dashboard.UpdateSummaryData()
		levelsBarChart.Labels, levelsBarChart.Data = dashboard.UpdateLevelData()
		radicalProgress.Percent, kanjiProgress.Percent = dashboard.UpdateLevelProgress()

		ui.Render(userText, summaryBarChart, levelsBarChart, radicalProgress, kanjiProgress)

	}

	uiEvents := ui.PollEvents()
	draw()
	ticker := time.NewTicker(1 * time.Minute).C

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			draw()
		}
	}
}

func main() {

	args := os.Args[1:]

	api.LoadConfig()

	if len(args) == 0 {
		fmt.Println("usage: waniKani [ dashboard | assignments | search \"<term>\" | set-config | help ]")
		os.Exit(0)
	}

	switch args[0] {
	case "dashboard":
		displayDashboard()
	case "assignments":
		api.CheckAssignments(args[1:])
	case "--update":
		tools.UpdateCache()
	case "search":
		tools.ProcessSearch(args[1])
	case "set-config":
		api.SetConfig(args[1], args[2])
	}

}
