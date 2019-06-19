package api

import (
	"bufio"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

func getConfigPath() string {

	configPath := ""
	GOOS := os.Getenv("GOOS")

	homeDirectory, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	switch GOOS {
	case "darwin":
		configPath = path.Join(homeDirectory, ".wanikani-cli")
	case "linux":
		configPath = path.Join(homeDirectory, ".config", "wanikani-cli")
	case "windows":
		configPath = path.Join(homeDirectory, "Documents", "WaniKani-cli")
	}

	return configPath

}

func LoadConfig() {

	configPath := getConfigPath()

	viper.SetConfigName("wkutil_config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	viperErr := viper.ReadInConfig()

	if viperErr != nil {
		fmt.Println(viperErr)
		fmt.Println("No config file detected, running first time setup...")
	}

	if viper.Get("apitoken") == nil {
		fmt.Println("No API token set. Please enter your WaniKani API v2 token:")
		InitialConfig()
	}
}

func InitialConfig() {

	baseConfigPath := getConfigPath()

	configPath := path.Join(baseConfigPath, "wkutil_config.yaml")

	if _, err := os.Stat(baseConfigPath); os.IsNotExist(err) {
		os.Mkdir(path.Join(baseConfigPath), 0755)
	}

	file, _ := os.OpenFile(configPath, os.O_RDONLY|os.O_CREATE, 0755)

	file.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')

		text = strings.Replace(text, "\n", "", -1)
		text = strings.Replace(text, "\t", "", -1)

		tokenRex, err := regexp.Compile("([A-Za-z0-9-]+)")

		if err != nil {
			log.Fatal(err)
		}

		if tokenRex.MatchString(text) {
			SetConfig("apitoken", text)
			fmt.Printf("\nAPI token set!\n\n")
			break
		} else {
			log.Fatal("Invalid API token format. Please use a v2 API token from your WaniKani account management page.")
		}
	}

	UpdateResetCache()

	UpdateSubjectsCache()

	fmt.Println("\nSetup complete! Run with --help to view usage")
	os.Exit(0)

}

func SetConfig(key, value string) {

	viper.Set(key, value)

	err := viper.WriteConfig()

	if err != nil {
		log.Fatal(err)
	}
}

const (
	ColorPink      ui.Color = 13
	ColorLightBlue ui.Color = 12
	ColorWhite     ui.Color = 15
	ColorBlack     ui.Color = 0
)
