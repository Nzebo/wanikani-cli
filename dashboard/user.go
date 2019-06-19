package dashboard

import (
	"encoding/json"
	"fmt"
	"github.com/gizak/termui/v3/widgets"
	"waniKani/api"
	"waniKani/tools"
)

func UpdateUserData() (string, string) {

	results := api.GetWaniKaniData("user")

	return userParse(results)
}

func User() *widgets.Paragraph {

	results := api.GetWaniKaniData("user")

	title, text := userParse(results)

	return userCreateElement(title, text)

}

func userParse(results []byte) (string, string) {

	var data api.UserData
	var title = "WaniKani Overview"
	var text string

	_ = json.Unmarshal(results, &data)

	username := data.Data.Username
	level := data.Data.Level
	subscription := ""

	switch data.Data.Subscribed {
	case true:
		subscription = "Active"
	case false:
		subscription = "Inactive"
	}

	api.GetLevelPercentage()

	text = fmt.Sprintf("Username: %s\nLevel: %v\nSubscription: %s", username, level, subscription)

	return title, text

}

func userCreateElement(title, text string) *widgets.Paragraph {

	ws := tools.GetTTYSize()

	x2 := ws.Width / 4

	t := widgets.NewParagraph()
	t.Title = title
	t.Text = text
	t.SetRect(0, 0, x2, 5)

	return t

}
