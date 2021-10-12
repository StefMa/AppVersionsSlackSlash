package model

import (
	"errors"
	"log"
	"stefma.guru/appVersionsSlackSlash/database"
	"stefma.guru/appVersionsSlackSlash/domain"
	"strings"
)

type SlashCommand interface {
	Execute() (string, error)
}

// Use this for commands like
// `/appversions [add|remove] android appId appId1 ...`
type appVersionsSlashCommand struct {
	db database.Database

	instruction     string
	operatingSystem string
	appIds          []string
}

// Use this for the "special" `get` instruction.
// `/appversions get`
type getSlashCommand struct {
	db database.Database

	instruction string
}

// BuildSlashCommand will construct an SlashCommand
func BuildSlashCommand(
	text string,
	db database.Database,
) (SlashCommand, error) {
	log.Println("Building slack command for: " + text)
	instructionAndArguments := strings.Split(text, " ")
	instruction := instructionAndArguments[0]
	switch instruction {
	case "get":
		log.Println("Build 'getSlashCommand'...")
		return getSlashCommand{
			db:          db,
			instruction: text,
		}, nil
	case "add":
		fallthrough
	case "lookup":
		fallthrough
	case "remove":
		log.Printf("Build 'appVersionsSlashCommand' for '%s'\n", instruction)
		operatingSystem := instructionAndArguments[1]
		appIds := instructionAndArguments[2:]
		return appVersionsSlashCommand{
			db:              db,
			instruction:     instruction,
			operatingSystem: operatingSystem,
			appIds:          appIds,
		}, nil
	default:
		return nil, errors.New("Unknown command. Only 'add', 'remove' & 'get' allowed")
	}
}

func (command getSlashCommand) Execute() (string, error) {
	return domain.Get(command.db)
}

func (command appVersionsSlashCommand) Execute() (string, error) {
	switch command.instruction {
	case "add":
		err := domain.Add(command.db, command.operatingSystem, command.appIds)
		return "Ok", err
	case "remove":
		err := domain.Remove(command.db, command.operatingSystem, command.appIds)
		return "Ok", err
	case "lookup":
		return domain.Lookup(command.operatingSystem, command.appIds)
	default:
		panic("We construct an 'appVersionsSlashCommand' with an unknown instructions: " + command.instruction)
	}
}
