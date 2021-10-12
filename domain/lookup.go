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
	var jsonResponse lookupJsonResponse
	if err := json.Unmarshal(bodyBytes, &jsonResponse); err != nil {
		return "", err
	}
	return createslackResponse(operatingSystem, jsonResponse)
}

func createslackResponse(
	operatingSystem string,
	lookupResponse lookupJsonResponse,
) (string, error) {
	jsonRequest := slackJsonRequest{
		Blocks: []slackBlock{},
	}
	for _, app := range lookupResponse.AndroidApps {
		jsonRequest.Blocks = append(
			jsonRequest.Blocks,
			createAppBlock(app, operatingSystem)...,
		)
	}
	for _, app := range lookupResponse.IosApps {
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
	app app,
	operatingSystem string,
) []slackBlock {
	text := fmt.Sprintf(
		"Name: *%s*\nVersion: *%s*\nRating: *%s*\n",
		app.Name, app.Version, app.Rating,
	)
	informationSectionBlock := slackBlock{
		Type: "section",
		Text: slackText{
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
	linksSectionBlock := slackBlock{
		Type: "section",
		Text: slackText{
			Type: "mrkdwn",
			Text: text,
		},
	}
	return []slackBlock{
		informationSectionBlock,
		linksSectionBlock,
	}
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

type slackJsonRequest struct {
	Blocks []slackBlock `json:"blocks"`
}

type slackBlock struct {
	Type string    `json:"type"`
	Text slackText `json:"text"`
}

type slackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
