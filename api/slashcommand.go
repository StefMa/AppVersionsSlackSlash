package api

import (
	"errors"
	"net/http"
	"stefma.guru/appVersionsSlackSlash/database"
	"stefma.guru/appVersionsSlackSlash/domain"
	"strings"
)

func HandleSlashCommand(
	w http.ResponseWriter,
	r *http.Request,
) {

	text, err := getQueryParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	commandAndArguments := strings.Split(text, " ")

	db, err := database.CreateDatabase()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	switch strings.TrimSpace(commandAndArguments[0]) {
	case "get":
		domain.Get(w, r, db)
		break
	case "add":
		domain.Add(w, r, db, commandAndArguments[1:])
		break
	case "remove":
		domain.Remove(w, r, db, commandAndArguments[1:])
		break
	default:
		http.Error(w, "Unknown command. Only 'add', 'remove' & 'add' allowed", http.StatusBadRequest)
	}
}

func getQueryParams(r *http.Request) (string, error) {
	if r.FormValue("command") == "" && r.FormValue("text") == "" {
		return "", errors.New("Wrong query params provided")
	}
	return r.FormValue("text"), nil
}
