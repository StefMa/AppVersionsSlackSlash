package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func Lookup(
	operatingSystem string,
	appIds []string,
) (string, error) {
	if operatingSystem != "android" && operatingSystem != "ios" {
		return "", errors.New("Unknown Operating System name (" + operatingSystem + ")!")
	}

	if len(appIds) == 0 {
		return "", errors.New("Please specify appIds for Operating System '" + operatingSystem + "'")
	}

	appIdsString := strings.Join(appIds, ",")
	url := fmt.Sprintf(
		"https://appversions.vercel.app/?%s=%s&format=json",
		operatingSystem,
		appIdsString,
	)
	log.Println("Lookup with url: " + url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var jsonResponse lookupJsonResponse
	if err := json.Unmarshal(bodyBytes, &jsonResponse); err != nil {
		return "", err
	}
	return createSlackResponse(operatingSystem, jsonResponse), nil
}

func createSlackResponse(
	operatingSystem string,
	lookupResponse lookupJsonResponse,
) string {
	jsonResponse := `{"blocks": [`
	appBlock := ``
	for _, app := range lookupResponse.AndroidApps {
		appBlock = createAppBlock(app, appBlock, operatingSystem)
	}
	for _, app := range lookupResponse.IosApps {
		appBlock = createAppBlock(app, appBlock, operatingSystem)
	}
	jsonResponse += appBlock
	jsonResponse += `]}`
	return jsonResponse
}

func createAppBlock(
	app app,
	appBlock string,
	operatingSystem string,
) string {
	text := fmt.Sprintf(
		"Name: *%s*\nVersion: *%s*\nRating: *%s*\n",
		app.Name, app.Version, app.Rating,
	)
	appBlock += `
	{
		"type": "section",
		"text": {
			"type": "mrkdwn",
			"text": "` + text + `"
		}
	},
	`
	text = fmt.Sprintf(
		"<%s|*[Store]*> <https://appversions.vercel.app?%s=%s|*[AppVersions]*>",
		app.URL,
		operatingSystem,
		app.ID,
	)
	appBlock += `
	{
		"type": "section",
		"text": {
			"type": "mrkdwn",
			"text": "` + text + `"
		}
	},
	`
	return appBlock
}

type lookupJsonResponse struct {
	AndroidApps []app `json:"android"`
	IosApps     []app `json:"ios"`
}

type app struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Rating  string `json:"rating"`
	URL     string `json:"url"`
}
