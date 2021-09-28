package database

func (db Database) Get() (map[string][]string, error) {
	docSnapsAndroid, err := db.client.
		Collection("android").
		Documents(db.context).
		GetAll()
	if err != nil {
		return nil, err
	}

	var androidAppIds []string
	for _, docSnap := range docSnapsAndroid {
		androidAppIds = append(androidAppIds, docSnap.Ref.ID)
	}

	docSnapsIos, err := db.client.
		Collection("ios").
		Documents(db.context).
		GetAll()
	if err != nil {
		return nil, err
	}

	var iosAppIds []string
	for _, docSnap := range docSnapsIos {
		iosAppIds = append(iosAppIds, docSnap.Ref.ID)
	}

	data := make(map[string][]string)
	data["ios"] = iosAppIds
	data["android"] = androidAppIds
	return data, nil
}
