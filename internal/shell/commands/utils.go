package commands

import (
	"log"
	"os"

	"github.com/i582/phpstats/internal/shell"
)

func handleOutputInJson(c *shell.Context) (bool, *os.File) {
	toJson := c.Flags.Contains("--json")
	var jsonFile *os.File

	if c.Flags.Contains("--json") {
		var err error
		jsonFile, err = c.ValidateFile("--json")
		if err != nil {
			log.Fatal(err)
		}
	}
	return toJson, jsonFile
}
