package database

import "fmt"

func (db Database) Get(operatingSystem string) (map[string][]string, error) {
	switch operatingSystem {
	case "android":
		androidAppIds, err := db.get("android")
		if err != nil {
			return nil, err
		}
		data := make(map[string][]string)
		data["android"] = androidAppIds
		return data, nil
	case "ios":
		iosAppIds, err := db.get("ios")
		if err != nil {
			return nil, err
		}
		data := make(map[string][]string)
		data["ios"] = iosAppIds
		return data, nil
	case "all":
		androidAppIds, err := db.get("android")
		if err != nil {
			return nil, err
		}
		iosAppIds, err := db.get("ios")
		if err != nil {
			return nil, err
		}
		data := make(map[string][]string)
		data["android"] = androidAppIds
		data["ios"] = iosAppIds
		return data, nil
	default:
		return nil, fmt.Errorf("Not supported operating system")
	}
}

func (db Database) get(operatingSystem string) ([]string, error) {
	docSnapsAndroid, err := db.client.
		Collection(operatingSystem).
		Documents(db.context).
		GetAll()
	if err != nil {
		return nil, err
	}

	var androidAppIds []string
	for _, docSnap := range docSnapsAndroid {
		androidAppIds = append(androidAppIds, docSnap.Ref.ID)
	}
	return androidAppIds, nil
}
