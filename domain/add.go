package domain

import (
	"stefma.guru/appVersionsSlackSlash/database"
)

func add(
	db database.Database,
	operatingSystem string,
	appIds []string,
) error {
	for _, appId := range appIds {
		err := db.Add(operatingSystem, appId)
		if err != nil {
			return err
		}
	}
	return nil
}
