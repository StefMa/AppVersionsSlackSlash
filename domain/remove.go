package domain

import (
	"stefma.guru/appVersionsSlackSlash/database"
)

func remove(
	db database.Database,
	operatingSystem string,
	appIds []string,
) error {
	for _, appId := range appIds {
		err := db.Remove(operatingSystem, appId)
		if err != nil {
			return err
		}
	}
	return nil
}
