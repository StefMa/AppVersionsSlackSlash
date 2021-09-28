package domain

import (
	"net/http"
	"stefma.guru/appVersionsSlackSlash/database"
)

func Add(
	w http.ResponseWriter,
	r *http.Request,
	db database.Database,
	arguments []string,
) {
	osAndAppIds := arguments
	appIds := osAndAppIds[1:]
	for _, appId := range appIds {
		err := db.Add(osAndAppIds[0], appId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.Write([]byte("Ok"))
	return
}
