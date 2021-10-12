package domain

import (
	"errors"
	"log"
	"stefma.guru/appVersionsSlackSlash/database"
	"stefma.guru/appVersionsSlackSlash/domain/model"
	"strings"
)

func Execute(command model.SlashCommand) (string, error) {
	switch command.(type) {
	case model.GetSlashCommand:
		return ExecuteGetSlashCommand(command.(model.GetSlashCommand))
	case model.AppVersionsSlashCommand:
		return ExecuteAppVersionSlashCommand(command.(model.AppVersionsSlashCommand))
	default:
		panic("Unknown SlashCommand type!")
	}
}

func ExecuteGetSlashCommand(command model.GetSlashCommand) (string, error) {
	return get(command.DB)
}

func ExecuteAppVersionSlashCommand(command model.AppVersionsSlashCommand) (string, error) {
	switch command.Instruction {
	case "add":
		err := add(command.DB, command.OperatingSystem, command.AppIds)
		return "Ok", err
	case "remove":
		err := remove(command.DB, command.OperatingSystem, command.AppIds)
		return "Ok", err
	case "lookup":
		return lookup(command.OperatingSystem, command.AppIds)
	default:
		panic("We construct an 'appVersionsSlashCommand' with an unknown instructions: " + command.Instruction)
	}
}

// BuildSlashCommand will construct an SlashCommand
func BuildSlashCommand(
	text string,
	db database.Database,
) (model.SlashCommand, error) {
	log.Println("Building slack command for: " + text)
	instructionAndArguments := strings.Split(text, " ")
	instruction := instructionAndArguments[0]
	switch instruction {
	case "get":
		log.Println("Build 'getSlashCommand'...")
		return model.GetSlashCommand{
			DB:          db,
			Instruction: text,
		}, nil
	case "add":
		fallthrough
	case "lookup":
		fallthrough
	case "remove":
		log.Printf("Build 'appVersionsSlashCommand' for '%s'\n", instruction)
		operatingSystem := instructionAndArguments[1]
		appIds := instructionAndArguments[2:]
		return model.AppVersionsSlashCommand{
			DB:              db,
			Instruction:     instruction,
			OperatingSystem: operatingSystem,
			AppIds:          appIds,
		}, nil
	default:
		return nil, errors.New("Unknown command. Only 'add', 'remove' & 'get' allowed")
	}
}
