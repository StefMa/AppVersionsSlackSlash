package database

import (
	"errors"
)

func (db Database) Add(operatingSystem string, appId string) error {
	if operatingSystem != "android" && operatingSystem != "ios" {
		return errors.New("Unknown Operating System name (" + operatingSystem + ")!")
	}

	data := map[string]interface{}{
		"dummyData": true,
	}
	_, err := db.client.
		Collection(operatingSystem).
		Doc(appId).
		Set(db.context, data)
	return err
}
