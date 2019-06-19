package api

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

func LoadLevelSubject(level int) *[]SubjectData {

	var Subjects *[]SubjectData

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	subjectJsonFile := path.Join(cwd, "_assets", "level_" + strconv.Itoa(level) + "_subjects.json")


	subjectFile, err := os.Open(subjectJsonFile)

	if err != nil {
		log.Fatal(err)
	}

	defer subjectFile.Close()

	subjectData, _ := ioutil.ReadAll(subjectFile)

	json.Unmarshal(subjectData, &Subjects)

	return Subjects

}

func ThrowError(message string, err error) {
	fmt.Printf("%v: %v\n", message, err)
	os.Exit(1)
}


// first argument should be 'resource'. Second arg should be 'optional_arg'
// have to use a variadic argument to support the optional second argument
func GetWaniKaniData(resource string) []byte {

	url := ""

	if strings.Contains(resource, "api.wanikani.com") {
		url = resource
	} else {
		url = fmt.Sprintf("https://api.wanikani.com/v2/%v", resource)
	}

	client := &http.Client{}
	var responseBody []byte

	req, err := http.NewRequest("GET", url, nil)

	headerValue := fmt.Sprintf("Bearer %s", viper.Get("ApiToken"))

	if err != nil {
		ThrowError("Error creating new request", err)
	}

	req.Header.Set("Authorization", headerValue)
	resp, err := client.Do(req)


	if err != nil {
		ThrowError("Error making initial request", err)
	}

	if resp.StatusCode != 401 {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			ThrowError("Error reading response body", err)
		}

		responseBody = body
	}

	return responseBody

}