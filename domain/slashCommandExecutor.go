package domain

import (
	"errors"
	"log"
	"strings"

	"stefma.guru/appVersionsSlackSlash/database"
	"stefma.guru/appVersionsSlackSlash/domain/model"
)

// BuildSlashCommand will construct an SlashCommand
func BuildSlashCommand(
	text string,
	db database.Database,
) (*model.SlashCommand, error) {
	log.Println("Building slack command for: " + text)

	instructionAndArguments := strings.Split(text, " ")
	instruction := instructionAndArguments[0]

	switch instruction {
	case "get":
		fallthrough
	case "add":
		fallthrough
	case "lookup":
		fallthrough
	case "remove":
		log.Printf("Build 'appVersionsSlashCommand' for '%s'\n", instruction)
		var operatingSystem string
		var args []string
		if len(instructionAndArguments) != 1 {
			operatingSystem = instructionAndArguments[1]
			args = instructionAndArguments[2:]
		}
		return &model.SlashCommand{
			DB:              db,
			Instruction:     instruction,
			OperatingSystem: operatingSystem,
			Args:            args,
		}, nil
	default:
		return nil, errors.New("Unknown command. Only 'add', 'remove' & 'get' allowed")
	}
}

func Execute(command *model.SlashCommand) (string, error) {
	switch command.Instruction {
	case "get":
		return get(command.DB, command.OperatingSystem)
	case "add":
		err := add(command.DB, command.OperatingSystem, command.Args)
		return "Ok", err
	case "remove":
		err := remove(command.DB, command.OperatingSystem, command.Args)
		return "Ok", err
	case "lookup":
		return lookup(command.OperatingSystem, command.Args)
	default:
		panic("We construct an 'appVersionsSlashCommand' with an unknown instructions: " + command.Instruction)
	}
}
