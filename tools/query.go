package tools

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/thedevsaddam/gojsonq"
	"gopkg.in/gookit/color.v1"
	"regexp"
	"strconv"
	"strings"
	"waniKani/api"

	"log"
	"os"
	"path"
)

func findFromJapanese(term string, jq *gojsonq.JSONQ, sData *[]api.SubjectData) {

	jq.WhereEqual("data.characters", term).Out(sData)

}

func isEnglish(term string) bool {

	isAlpha := regexp.MustCompile(`^[A-Za-z ]+$`).MatchString

	if isAlpha(term) {
		return true
	} else {
		return false
	}

}

func findFromEnglish(term string, jq *gojsonq.JSONQ, sData *[]api.SubjectData) {

	jq.WhereEqual("data.meanings.[0].meaning", term).Out(sData)

}

func getReviewStastics(subjects *[]api.SubjectData, reviewData *api.ReviewCollection) {

	subjectIds := make([]string, 0)

	for _, subject := range *subjects {
		subjectIds = append(subjectIds, strconv.Itoa(subject.Id))
	}

	requestUrl := "https://api.wanikani.com/v2/review_statistics?subject_ids=" + strings.Join(subjectIds, ",")


	results := api.GetWaniKaniData(requestUrl)

	json.Unmarshal(results, reviewData)

}

func ProcessSearch(term string) {

	userSubsctibed := api.CheckSubscription()
	hiddenCount := 0

	var s []api.SubjectData
	var r api.ReviewCollection

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	subjectsPath := path.Join(cwd, "_assets", "subjects.json")

	jq := gojsonq.New().File(subjectsPath)

	term = strings.Title(term)

	if isEnglish(term) {
		findFromEnglish(term, jq, &s)

	} else {
		findFromJapanese(term, jq, &s)
	}

	getReviewStastics(&s, &r)

	for _, subject := range s {

		if userSubsctibed == false {
			if subject.Data.Level >= 4 {
				hiddenCount++
			}

		}


		fmt.Println("\n\n----------------------------------------")

		var p = color.RGBColor{}

		switch subject.Object {
		case "radical":
			p = color.RGB(0, 170, 255)
		case "kanji":
			p = color.RGB(255, 0, 170)
		case "vocabulary":
			p = color.RGB(170 ,0, 255)
		}

		p.Printf("\n%v\n\n", strings.Title(subject.Object))
		color.Printf("%v\n\n", subject.Data.Characters)
		color.Printf("<op=underscore;>Level</>: %v\n\n", subject.Data.Level)
		color.Printf("<op=underscore;>Meaning</>: %v\n\n", subject.Data.Meanings[0].Meaning)
		color.Printf("<op=underscore;>Meaning Mnemonic</>:\n\n%v\n\n", parseMnemonic(subject.Data.MeaningMnemonic))

		parseReadings(subject.Data.Readings)
		parsePartsOfSpeech(subject.Data.PartsOfSpeech)
		parseContextSentences(subject.Data.ContextSentences)
		parseReviewStatistics(subject, &r)


	}

	if len(s) == 0 {
		fmt.Printf("\nNo results for '%v'\n", term)
	}

	if userSubsctibed == false {
		if hiddenCount > 0 {
			fmt.Printf("%v results were hidden because you do not have an active WaniKani subscription.\n\n", hiddenCount)
		}
		fmt.Printf("Sign up for a WaniKani subscription to support this awesome service and gain access to all of these resources!\n")
		fmt.Println("\n\nSubscribe here: https://www.wanikani.com/account/subscription")
	}

	fmt.Println()

}

func parseReviewStatistics(subject api.SubjectData, r *api.ReviewCollection) {
	var reviewInfo api.ReviewData

	for _, review := range r.Reviews {
		if review.Data.SubjectId == subject.Id {
			reviewInfo = review
		}
	}

	if reviewInfo.Id != 0 {
		color.Printf("<op=underscore;>User Statistics</>\n\n")
		fmt.Printf("  * Overall Correct Percentage: %v\n\n", reviewInfo.Data.PercentageCorrect)
		fmt.Printf("  * Meaning Current Streak: %v\n", reviewInfo.Data.MeaningCurrentStreak)
		fmt.Printf("  * Meaning Correct: %v\n", reviewInfo.Data.MeaningCorrect)
		fmt.Printf("  * Meaning Incorrect: %v\n\n", reviewInfo.Data.MeaningIncorrect)
		fmt.Printf("  * Reading Current Streak: %v\n", reviewInfo.Data.ReadingCurrentStreak)
		fmt.Printf("  * Reading Correct: %v\n", reviewInfo.Data.ReadingCorrect)
		fmt.Printf("  * Reading Incorrect %v\n", reviewInfo.Data.ReadingIncorrect)
	}

}

func parseMnemonic(mnemonic string) string {

	radicalColor := color.RGB(0, 170, 255)
	kanjiColor := color.RGB(255, 0, 170)
	vocabularyColor := color.RGB(170, 0, 255)

	radicalRe := regexp.MustCompile(`<radical>([^<]+)</radical>`)
	kanjiRe := regexp.MustCompile(`<kanji>([^<]+)</kanji>`)
	vocabularyRe := regexp.MustCompile(`<vocabulary>([^<]+)</vocabulary>`)
	jaRe := regexp.MustCompile(`<ja>([^<]+)</ja>`)

	mnemonic = radicalRe.ReplaceAllString(mnemonic, radicalColor.Sprintf("$1"))
	mnemonic = kanjiRe.ReplaceAllString(mnemonic, kanjiColor.Sprintf("$1"))
	mnemonic = vocabularyRe.ReplaceAllString(mnemonic, vocabularyColor.Sprintf("$1"))
	mnemonic = jaRe.ReplaceAllString(mnemonic, fmt.Sprintf("$1"))

	return mnemonic

}

func parsePartsOfSpeech(parts []string) {

	if len(parts) == 0 {
		return
	}

	color.Printf("<op=underscore;>Parts of Speech:</>\n\n")

	fmt.Printf("%v\n\n", strings.Join(parts, ", "))
}

func parseContextSentences(sentences []api.ContextSentences) {

	if len(sentences) == 0 {
		return
	}

	color.Printf("<op=underscore;>Context Sentences:</>\n\n")

	for _, sentence := range sentences {
		fmt.Println(sentence.Japanese)
		fmt.Println(sentence.English)
		fmt.Println()
	}

}

func parseReadings(readings []api.Readings) {

	if len(readings) == 0 {
		return
	}

	color.Println("<op=underscore;>Readings</>\n")

	readingHeader := []string{"Type", "Reading", "Primary"}
	readingData := make([][]string, 0)

	for _, reading := range readings {

		data := []string{ reading.Type, reading.Reading, strconv.FormatBool(reading.Primary)}

		readingData = append(readingData, data)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(readingHeader)
	table.SetBorder(false)
	table.AppendBulk(readingData)
	table.Render()
	fmt.Println()

}
