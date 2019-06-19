package api

import "encoding/json"


type UserData struct {
	Data struct {
		Username string `json:"username"`
		Level    int    `json:"level"`
		MaxLevelGrantedBySubscription int	`json:"max_level_granted_by_subscription"`
		Subscribed					  bool	`json:"subscribed"`
		Subscription struct {
			Active					  bool	`json:"active"`
			Type 					  string  `json:"type"`
			MaxLevelGranted			  string  `json:"max_level_granted"`
			PeriodEndsAt			  string  `json:"period_ends_at"`
		} `json:"subscription"`

	} `json:"data"`
}

func GetUserLevel() int {
	results := GetWaniKaniData("user")

	var data UserData

	_ = json.Unmarshal(results, &data)

	return data.Data.Level
}

func CheckSubscription() bool {
	results := GetWaniKaniData("user")

	var data UserData

	_ = json.Unmarshal(results, &data)

	return data.Data.Subscribed
}
