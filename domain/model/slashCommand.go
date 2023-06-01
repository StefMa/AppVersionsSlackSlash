package model

import (
	"stefma.guru/appVersionsSlackSlash/database"
)

type SlashCommand struct {
	DB database.Database

	Instruction     string
	OperatingSystem string
	Args            []string
}
