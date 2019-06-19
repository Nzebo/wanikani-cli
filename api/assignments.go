package api

import (
	"encoding/json"
	"fmt"
	"gopkg.in/gookit/color.v1"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type AssignmentCollection struct {
	Pages struct {
		PerPage int    `json:"per_page"`
		NextURL string `json:"next_url"`
	} `json:"pages"`

	TotalCount  int          `json:"total_count"`
	Assignments []Assignment `json:"data"`
}

type Assignment struct {
	Id   int    `json:"id"`
	URL  string `json:"url"`
	Data struct {
		SubjectID     int    `json:"subject_id"`
		SubjectType   string `json:"subject_type"`
		SRSStage      int    `json:"srs_stage"`
		SRSStageName  string `json:"srs_stage_name"`
		UnlockedAt    string `json:"unlocked_at"`
		StartedAt     string `json:"started_at"`
		PassedAt      string `json:"passed_at"`
		BurnedAt      string `json:"burned_at"`
		AvailableAt   string `json:"available_at"`
		ResurrectedAt string `json:"resurrected_at"`
		Passed        bool   `json:"passed"`
		Resurrected   bool   `json:"resurrected"`
		Hidden        bool   `json:"hidden"`
	} `json:"data"`
}

func CheckAssignments(args []string) {

	var assignmentResults []byte
	var arguments = []string{"immediately_available_for_review"}

	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	baseUrl := "https://api.wanikani.com/v2/assignments?"

	for _, arg := range args {
		switch arg {
		case "--today":
			year, month, day := time.Now().Date()
			date := time.Date(year, month, day, 23, 59, 59, 0, time.Local).UTC().Format("2006-01-02T15:04:05Z")
			arguments = append(arguments[1:], "available_before="+date+"&available_after="+now)
		case "--radicals": //todo fix the subject type filter
			arguments = append(arguments, "subject_types=radical")
		case "--kanji":
			arguments = append(arguments, "subject_types=kanji")
		default:
			fmt.Printf("Invalid flag '%v'. Valid flags are [--today, --lessons]\n", arg) //todo add list of flags
			os.Exit(1)
		}

	}

	requestUrl := baseUrl + strings.Join(arguments, "&")

	assignmentResults = GetWaniKaniData(requestUrl)

	var assignmentData AssignmentCollection

	_ = json.Unmarshal(assignmentResults, &assignmentData)

	radicalCount := 0
	kanjiCount := 0
	vocabularyCount := 0

	for _, item := range assignmentData.Assignments {

		switch item.Data.SubjectType {
		case "radical":
			radicalCount++
		case "kanji":
			kanjiCount++
		case "vocabulary":
			vocabularyCount++
		}
	}

	fmt.Printf("\nReviews available:\n\n")
	color.RGB(0, 170, 255).Printf("Radicals: %v\n", radicalCount)
	color.RGB(255, 0, 170).Printf("Kanji: %v\n", kanjiCount)
	color.RGB(170, 0, 255).Printf("Vocabulary: %v\n\n", vocabularyCount)

}

func GetLevelPercentage() (int, int) {

	var assignments AssignmentCollection

	reset := GetLatestReset()

	level := GetUserLevel()

	totalRadicalAssignments := 0.0
	totalKanjiAssignments := 0.0

	passedRadicalAssignments := 0.0
	passedKanjiAssignments := 0.0

	url := "https://api.wanikani.com/v2/assignments?levels=" + strconv.Itoa(level)

	currentProgress := GetWaniKaniData(url)

	_ = json.Unmarshal(currentProgress, &assignments)

	for _, review := range assignments.Assignments {

		if review.Data.StartedAt == "" {
			continue
		}

		startTime, err := time.Parse(time.RFC3339, review.Data.StartedAt)

		if err != nil {
			log.Fatal(err)
		}

		if startTime.Unix() < reset.Unix() {
			continue
		}

		switch review.Data.SubjectType {
		case "radical":
			totalRadicalAssignments++
			if review.Data.Passed == true {
				passedRadicalAssignments++
			}
		case "kanji":
			totalKanjiAssignments++
			if review.Data.Passed == true {
				passedKanjiAssignments++
			}
		}

	}

	return int(passedRadicalAssignments / totalRadicalAssignments * 100), int(passedKanjiAssignments / totalKanjiAssignments * 100)

}
