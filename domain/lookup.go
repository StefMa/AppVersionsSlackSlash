package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"stefma.guru/appVersionsSlackSlash/domain/model"
	"strings"
)

func lookup(
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
	var jsonResponse model.AppVersionsJsonResponse
	if err := json.Unmarshal(bodyBytes, &jsonResponse); err != nil {
		return "", err
	}
	return createslackResponse(operatingSystem, jsonResponse)
}

func createslackResponse(
	operatingSystem string,
	appVersionsJsonResponse model.AppVersionsJsonResponse,
) (string, error) {
	jsonRequest := model.SlackJsonRequest{
		Blocks: []model.SlackBlock{},
	}
	for _, app := range appVersionsJsonResponse.AndroidApps {
		jsonRequest.Blocks = append(
			jsonRequest.Blocks,
			createAppBlock(app, operatingSystem)...,
		)
	}
	for _, app := range appVersionsJsonResponse.IosApps {
		jsonRequest.Blocks = append(
			jsonRequest.Blocks,
			createAppBlock(app, operatingSystem)...,
		)
	}

	jsonBytes, err := json.Marshal(&jsonRequest)
	if err != nil {
		return "", err
	}
	log.Println("Sending json string: " + string(jsonBytes))
	return string(jsonBytes), nil
}

func createAppBlock(
	app model.AppVersionsApp,
	operatingSystem string,
) []model.SlackBlock {
	text := fmt.Sprintf(
		"Name: *%s*\nVersion: *%s*\nRating: *%s*\n",
		app.Name, app.Version, app.Rating,
	)
	informationSectionBlock := model.SlackBlock{
		Type: "section",
		Text: model.SlackText{
			Type: "mrkdwn",
			Text: text,
		},
	}

	text = fmt.Sprintf(
		"<%s|*[Store]*> <https://appversions.vercel.app?%s=%s|*[AppVersions]*>",
		app.URL,
		operatingSystem,
		app.ID,
	)
	linksSectionBlock := model.SlackBlock{
		Type: "section",
		Text: model.SlackText{
			Type: "mrkdwn",
			Text: text,
		},
	}
	return []model.SlackBlock{
		informationSectionBlock,
		linksSectionBlock,
	}
}
