package dashboard

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	image2 "image"
	_ "image/png"
	"log"
	"net/http"
	"waniKani/api"
)

func DrawStudyDashboard(data *api.SubjectData) {

	imageUrl := data.Data.CharacterImages[4].URL

	var images []image2.Image

	resp, err := http.Get(imageUrl)

	if err != nil {
		log.Fatalf("failed to fetch image: %v", err)
	}

	fmt.Println(resp.Body)

	image2, _, err := image2.Decode(resp.Body)
	if err != nil {
		log.Fatalf("failed to decode fetched image: %v", err)
	}
	images = append(images, image2)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()
	defer resp.Body.Close()

	img := widgets.NewImage(nil)
	img.SetRect(0, 0, 50, 25)

	index := 0

	render := func() {
		img.Image = images[index]

		if !img.Monochrome {
			img.Title = fmt.Sprintf("Color %d/%d - %d", index+1, len(images), data.Id)
		} else if !img.MonochromeInvert {
			img.Title = fmt.Sprintf("Monochrome(%d) %d/%d", img.MonochromeThreshold, index+1, len(images))
		} else {
			img.Title = fmt.Sprintf("InverseMonochrome(%d) %d/%d", img.MonochromeThreshold, index+1, len(images))
		}

		ui.Render(img)
	}

	render()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Left>", "h":
			index = (index + len(images) - 1) % len(images)
		case "<Right>", "l":
			index = (index + 1) % len(images)
		case "<Up>", "k":
			img.MonochromeThreshold++
		case "<Down>", "j":
			img.MonochromeThreshold--
		case "<Enter>":
			img.Monochrome = !img.Monochrome
		case "<Tab>":
			img.MonochromeInvert = !img.MonochromeInvert
		}
		render()
	}


}
