package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	_ "github.com/gizak/termui/v3/widgets"
	"os"
	"time"
	"waniKani/api"
	"waniKani/dashboard"
	"waniKani/tools"
)

func displayDashboard(args []string) {

	if len(args) == 0 {
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
		ticker := time.NewTicker(5 * time.Minute).C

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
	} else if args[0] == "help" {
		fmt.Println("\nA dashboard displaying some graphical metrics for the user's WaniKani account.")
		fmt.Println("\nThis dashboard automatically refreshes every 5 minutes while open.")
		fmt.Println("\nThe dashboard contains:")
		fmt.Println("\n    Overview\t\tAn overview panel displaying the Username, Level, and Subscription status")
		fmt.Println("    Upcoming Reviews\tA bar graph showing the number of reviews available over time")
		fmt.Println("    Time Per Level\tA bar graph showing the time it took to level up for each level")
		fmt.Println("    Current Percentage\tA percentage gauge displaying the user's current percentage for both Radicals and Kanji\n")
	}

}

func helpMessage() {

	fmt.Println("A WaniKani cli tool to display user info and metrics.\n\n")

	fmt.Println("Usage:\n")
	fmt.Println("        wanikani-cli <command> [arguments]")

	fmt.Println("\nAvailable commands are:\n\n")
	fmt.Println("        dashboard\tDisplays a metrics dashboard for the user's account")
	fmt.Println("        assignments\tReturns WaniKani assignment counts for the user")
	fmt.Println("        search\t\tSearches the WaniKani subject content for the provided term")
	fmt.Println("        set-config\tUpdates the supplied configuration setting with the supplied value\n")
	fmt.Println("Use \"wanikani-cli <command> <help>\" for details on any of the above commands.\n")

}

func main() {

	args := os.Args[1:]

	api.LoadConfig()

	if len(args) == 0 {
		fmt.Println("usage: wanikani-cli [ dashboard | assignments | search \"<term>\" | set-config ] --help")
		os.Exit(0)
	}

	switch args[0] {
	case "--help":
		helpMessage()
	case "dashboard":
		displayDashboard(args[1:])
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
