package api

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"
)

type ResetData struct {
	Resets 				[]ResetObj `json:"data"`
}

type ResetObj struct {
	ResetDetails struct {
		Created 	  	string 	   `json:"created_at"`
		OriginalLevel 	int 	   `json:"original_level"`
		TargetLevel	  	int 	   `json:"target_level"`
		Confirmed	  	string 	   `json:"confirmed_at"`
	} `json:"data"`
}

func UpdateResetCache() {

	results := parseResetTime(GetWaniKaniData("resets"))

	if results != "" {

		fmt.Println("\nDetected an account reset, caching this timestamp to provide accurate statistics.\n")

		viper.Set("latest_reset", results)
		err := viper.WriteConfig()

		if err != nil {
			log.Fatal(err)
		}

	}

}

func parseResetTime(results []byte) string {

	var data ResetData

	_ = json.Unmarshal(results, &data)

	if len(data.Resets) > 0 {
		resetConfirmTime := data.Resets[len(data.Resets)-1].ResetDetails.Confirmed

		return resetConfirmTime
	}

	return ""
}

func GetLatestReset() time.Time {
	resetString := viper.GetString("latest_reset")

	if resetString != "" {
		resetTime, err := time.Parse(time.RFC3339, resetString)

		if err != nil {
			log.Fatal(err)
		}

		return resetTime

	} else {
		t, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00")

		return t
	}

}