package domain

import (
	"stefma.guru/appVersionsSlackSlash/database"
)

func Add(
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
