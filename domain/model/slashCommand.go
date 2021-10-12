package model

import (
	"stefma.guru/appVersionsSlackSlash/database"
)

type SlashCommand interface{}

// Use this for commands like
// `/appversions [add|remove] android appId appId1 ...`
type AppVersionsSlashCommand struct {
	DB database.Database

	Instruction     string
	OperatingSystem string
	AppIds          []string
}

// Use this for the "special" `get` instruction.
// `/appversions get`
type GetSlashCommand struct {
	DB database.Database

	Instruction string
}
