package domain

import (
	"net/http"
	"stefma.guru/appVersionsSlackSlash/database"
)

func Remove(
	w http.ResponseWriter,
	r *http.Request,
	db database.Database,
	arguments []string,
) {
	osAndAppIds := arguments
	appIds := osAndAppIds[1:]
	for _, appId := range appIds {
		err := db.Remove(osAndAppIds[0], appId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.Write([]byte("Ok"))
}
