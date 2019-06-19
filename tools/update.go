package tools

import (
	"fmt"
	"os"
	"waniKani/api"
)

func UpdateCache() {

	api.UpdateResetCache()

	api.UpdateSubjectsCache()

	fmt.Println("\nCache update complete!")
	os.Exit(0)

}
