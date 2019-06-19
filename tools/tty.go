package tools

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type WindowSize struct {
	Width  int
	Height int
}

func GetTTYSize() *WindowSize {

	ws := &WindowSize{}
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	str := strings.Split(string(out), " ")
	width, _ := strconv.Atoi(strings.Replace(str[1], "\n", "", 1))
	height, _ := strconv.Atoi(str[0])

	ws.Width = width
	ws.Height = height

	return ws
}
