package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"

	"github.com/i582/phpstats/internal/cli"
)

type CommandsMap struct {
	Commands map[string]*Command
}

func NewCommandsMap() *CommandsMap {
	return &CommandsMap{
		Commands: map[string]*Command{},
	}
}

func (c *CommandsMap) AddCommand(command *Command) {
	c.Commands[command.Name] = command
}

func (c *CommandsMap) GetCommand(name string) (*Command, bool) {
	command, ok := c.Commands[name]
	return command, ok
}

func (c *CommandsMap) GetSuggests(commands string) []prompt.Suggest {
	parts := strings.Fields(commands)
	if len(parts) == 0 {
		parts = []string{commands}
	}
	if strings.HasSuffix(commands, " ") {
		parts = append(parts, "")
	}

	return c.getSuggests(parts)
}

func (c *CommandsMap) getSuggests(commands []string) []prompt.Suggest {
	mainCommand := commands[0]

	if len(commands) == 1 {
		var suggests []prompt.Suggest

		commands := make([]*Command, 0, len(c.Commands))
		for _, command := range c.Commands {
			commands = append(commands, command)
		}
		sort.Slice(commands, func(i, j int) bool {
			return commands[i].Name < commands[j].Name
		})

		for _, command := range commands {
			if strings.HasPrefix(command.Name, mainCommand) {
				suggests = append(suggests, prompt.Suggest{
					Text:        command.Name,
					Description: command.Description,
				})
			}
		}
		return suggests
	}

	command, found := c.GetCommand(mainCommand)
	if found {
		return command.SubCommands.getSuggests(commands[1:])
	}

	return nil
}

type Command struct {
	Name        string
	Description string
	SubCommands *CommandsMap
}

func NewCommand(name string, description string) *Command {
	return &Command{
		Name:        name,
		Description: description,
		SubCommands: NewCommandsMap(),
	}
}

func (c *Command) AddSubCommand(command *Command) {
	c.SubCommands.AddCommand(command)
}

func completer(d prompt.Document) []prompt.Suggest {
	return commandsMap.GetSuggests(d.CurrentLine())
	//
	// firstLevelCommands := []prompt.Suggest{
	// 	{Text: "info", Description: "info about"},
	// 	{Text: "graph", Description: "graph for"},
	// 	{Text: "list", Description: "list of"},
	// }
	//
	// secondLevelCommands := []prompt.Suggest{
	// 	{Text: "class", Description: "list of"},
	// 	{Text: "func", Description: "list of"},
	// 	{Text: "namespace", Description: "list of"},
	// }
	//
	// return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

func executor(in string) {
	fmt.Println("Your input: " + in)
	if in == "" {
		LivePrefixState.IsEnable = false
		LivePrefixState.LivePrefix = in
		return
	}
	// LivePrefixState.LivePrefix = ">>> "
	// LivePrefixState.IsEnable = true
}

func changeLivePrefix() (string, bool) {
	return LivePrefixState.LivePrefix, LivePrefixState.IsEnable
}

var commandsMap = NewCommandsMap()

func main() {
	// infoCommand := NewCommand("info", "info about")
	// graphCommand := NewCommand("graph", "info about")
	// listCommand := NewCommand("list", "info about")
	//
	// classCommand := NewCommand("class", "info about")
	// funcCommand := NewCommand("func", "info about")
	// namespaceCommand := NewCommand("namespace", "info about")
	//
	// infoCommand.AddSubCommand(classCommand)
	// infoCommand.AddSubCommand(funcCommand)
	// infoCommand.AddSubCommand(namespaceCommand)
	//
	// graphCommand.AddSubCommand(classCommand)
	// graphCommand.AddSubCommand(funcCommand)
	// graphCommand.AddSubCommand(namespaceCommand)
	//
	// listCommand.AddSubCommand(classCommand)
	// listCommand.AddSubCommand(funcCommand)
	// listCommand.AddSubCommand(namespaceCommand)
	//
	// commandsMap.AddCommand(infoCommand)
	// commandsMap.AddCommand(graphCommand)
	// commandsMap.AddCommand(listCommand)
	//
	// p := prompt.New(
	// 	executor,
	// 	completer,
	// 	prompt.OptionPrefix(">>> "),
	// 	prompt.OptionLivePrefix(changeLivePrefix),
	// 	prompt.OptionSuggestionBGColor(prompt.Yellow),
	// 	prompt.OptionSuggestionTextColor(prompt.DarkGray),
	// 	prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
	// 	prompt.OptionSelectedSuggestionTextColor(prompt.White),
	//
	// 	prompt.OptionDescriptionBGColor(prompt.DarkGray),
	// 	prompt.OptionDescriptionTextColor(prompt.White),
	// 	prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
	// 	prompt.OptionSelectedDescriptionTextColor(prompt.White),
	// )
	// p.Run()
	cli.Run()
}
