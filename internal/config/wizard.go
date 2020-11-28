package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/i582/cfmt"
	"github.com/manifoldco/promptui"

	"github.com/i582/phpstats/internal/utils"
)

var configTemplate = `# Project name
project-name: "%s"

# Directories and files for analysis relative to the configuration files directory.
# By default, it is "."
include:
  - "%s"

# Directories and files excluded from analysis.
# Please note that for correct work, you need to add a slash at the end
# so that only the necessary folders are excluded, and not all that have the same value in the path.
#
# For example:
#    "src/utils" can exclude both the desired "src/utils" folder and the "src/utilsForMe" folder for example.
#
# By default, it is empty
# exclude:
#   - ""

# The port on which the server will be launched
# to interact with the analyzer from other programs.
# By default, it is 8080
port: 8080

# The path where the cache will be stored.
# Caching can significantly speed up data collection.
# By default, it is set to the value of the temporary folder + /phpstats.
cacheDir: ""

# Disables caching.
# By default, it is false
disableCache: false

# Path to the project relative to which all imports are allowed.
# By default, it is equal to the analyzed directory.
projectPath: ""

# File extensions to be included in the analysis.
# By default, it is php, inc, php5, phtml.
extensions:
  - "php"
  - "inc"
  - "php5"
  - "phtml"
`

func ConfigureConfig() {
	fmt.Println("Welcome to the phpstats configuration wizard.")
	fmt.Println("You will need to answer a few questions to complete.")

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Some unexpected error, try again")
	}

	projectName := projectNamePrompt(err)
	sourcePath := sourcePathPrompt(workingDir, err)

	cfmt.Println("{{v}}::green Thanks for answers. The config was {{successfully}}::green created.")

	file, err := os.OpenFile(filepath.Join(workingDir, "phpstats.yml"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0677)
	if err != nil {
		log.Fatalf("file not open %v", err)
	}
	defer file.Close()

	configData := fmt.Sprintf(configTemplate, projectName, sourcePath)

	fmt.Fprint(file, configData)
}

func sourcePathPrompt(workingDir string, err error) string {
	validatePathToSource := func(input string) error {
		if len(input) == 0 {
			return errors.New("Path to source must not be empty (for current dir enter ./)")
		}

		path := filepath.Join(workingDir, input)

		if exists, err := utils.Exists(path); !exists && err == nil {
			return fmt.Errorf("The specified path (%s) is invalid", path)
		}

		return nil
	}

	promptSourcePath := promptui.Prompt{
		Label:       "Path to source (the relative path to the folder with the code)",
		Validate:    validatePathToSource,
		HideEntered: true,
	}

	sourcePath, err := promptSourcePath.Run()
	if err != nil {
		log.Fatal(err)
	}

	return sourcePath
}

func projectNamePrompt(err error) string {
	validate := func(input string) error {
		if len(input) == 0 {
			return errors.New("Project name must not be empty")
		}
		return nil
	}

	promptProjectName := promptui.Prompt{
		Label:       "Project name",
		Validate:    validate,
		HideEntered: true,
	}

	projectName, err := promptProjectName.Run()
	if err != nil {
		log.Fatalf("Some unexpected error, try again: %v", err)
	}

	return projectName
}
