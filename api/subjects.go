package api

import (
	"encoding/json"
	"fmt"
	"github.com/kyokomi/emoji"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

type SubjectCollection struct {
	Url   string `json:"url"`
	Pages struct {
		NextURL string `json:"next_url"`
	} `json:"pages"`

	Subjects []SubjectData `json:"data"`
}

type SubjectData struct {
	Id     int    `json:"id"`
	Object string `json:"object"`
	Data   struct {
		AmalgamationSubjectIds []int               `json:"amalgamation_subject_ids"`
		AuxiliaryMeanings      []AuxiliaryMeanings `json:"auxiliary_meanings"`
		Characters             string              `json:"characters"`
		CharacterImages        []CharacterImage    `json:"character_images"`
		ContextSentences       []ContextSentences  `json:"context_sentences"`
		DocumentURL            string              `json:"document_url"`
		HiddenAt               string              `json:"hidden_at"`
		LessonPosition         int                 `json:"lesson_position"`
		Level                  int                 `json:"level"`
		Meanings               []Meanings          `json:"meanings"`
		MeaningMnemonic        string              `json:"meaning_mnemonic"`
		MeaningHint            string              `json:"meaning_hint"`
		Readings               []Readings          `json:"readings"`
		ReadingMnemonic        string              `json:"reading_mnemonic"`
		ReadingHint            string              `json:"reading_hint"`
		PartsOfSpeech          []string            `json:"parts_of_speech"`
		Slug                   string              `json:"slug"`
	} `json:"data"`
}

type AuxiliaryMeanings struct {
	Meaning string `json:"meaning"`
	Type    string `json:"type"`
}

type CharacterImage struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
	Metadata    struct {
		InlineStyles string `json:"inline_styles"`
	} `json:"metadata"`
}

type Meanings struct {
	Meaning        string `json:"meaning"`
	Primary        bool   `json:"primary"`
	AcceptedAnswer bool   `json:"accepted_answer"`
}

type Readings struct {
	Type           string `json:"type"`
	Primary        bool   `json:"primary"`
	AcceptedAnswer string `json:"accepted_answer"`
	Reading        string `json:"reading"`
}

type ContextSentences struct {
	English  string `json:"en"`
	Japanese string `json:"ja"`
}

func UpdateSubjectsCache() {

	var fullSubjectCache []SubjectData

	fmt.Println("Updating subjects cache, this may take a few minutes...")

	contentLimit := 0

	baseConfigPath := getConfigPath()

	configPath := path.Join(baseConfigPath, "_assets")

	// Check the user's subscription status to respect WaniKani's content restrictions
	isSubscribed := CheckSubscription()

	if isSubscribed {
		emoji.Println("\n* You have an active subscription! :tada: Caching content for every level.\n\n")
		contentLimit = 60
	} else {
		fmt.Printf("\n* You do not have an active subscription. Only content for levels 1-3 will be cached.\n\n")
		contentLimit = 3
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.Mkdir(configPath, 0755)
	}

	for i := 1; i <= contentLimit; i++ {
		var subjectCache []SubjectData

		getSubjects(i, &subjectCache, &fullSubjectCache)

		outputData, err := json.Marshal(subjectCache)

		if err != nil {
			log.Fatal(err)
		}

		outputFile := path.Join(configPath, "level_"+strconv.Itoa(i)+"_subjects.json")
		err = ioutil.WriteFile(outputFile, outputData, 0755)

		fmt.Printf("Updated level %v subject data\n", i)

		time.Sleep(3 * time.Second)

	}

	outputData, err := json.Marshal(fullSubjectCache)

	if err != nil {
		log.Fatal(err)
	}

	outputFile := path.Join(configPath, "subjects.json")
	err = ioutil.WriteFile(outputFile, outputData, 0755)

	fmt.Println("\nSuccessfully updated subjects cache!")

}

func getSubjects(level int, subjectCache *[]SubjectData, fullSubjectCache *[]SubjectData, url ...string) {

	var data SubjectCollection
	var subjectData []byte

	baseUrl := "https://api.wanikani.com/v2/subjects?"

	if len(url) == 0 {
		subjectData = GetWaniKaniData(baseUrl + "levels=" + strconv.Itoa(level))
	} else {
		subjectData = GetWaniKaniData(url[0])
	}

	_ = json.Unmarshal(subjectData, &data)

	nextURL := data.Pages.NextURL

	for i := 0; i < len(data.Subjects); i++ {
		*subjectCache = append(*subjectCache, data.Subjects[i])
		*fullSubjectCache = append(*fullSubjectCache, data.Subjects[i])
	}

	if nextURL != "" {
		getSubjects(level, subjectCache, fullSubjectCache, nextURL)
	}
}
