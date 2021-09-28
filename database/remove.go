package database

import (
	"errors"
)

func (db Database) Remove(operatingSystem string, appId string) error {
	if operatingSystem != "android" && operatingSystem != "ios" {
		return errors.New("Unknown Operating System name (" + operatingSystem + ")!")
	}

	_, err := db.client.
		Collection(operatingSystem).
		Doc(appId).
		Delete(db.context)
	return err
}
