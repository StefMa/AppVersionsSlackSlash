package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"stefma.guru/appVersionsSlackSlash/database"
	"stefma.guru/appVersionsSlackSlash/domain/model"
	"strings"
)

func HandleSlashCommand(
	writer http.ResponseWriter,
	request *http.Request,
) {
	rawBody, err := readBody(request)
	if err != nil {
		writer.Write([]byte("Error while reading HTTP body."))
		return
	}

	if valid := verifySignature(request.Header, rawBody); !valid {
		writer.Write([]byte("Signature doesn't match with calculated one."))
		return
	}

	text, err := getTextParam(rawBody)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	db, err := database.CreateDatabase()
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	defer db.Close()

	command, err := model.BuildSlashCommand(text, db)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	result, err := command.Execute()
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.Write([]byte(result))
}

func readBody(request *http.Request) ([]byte, error) {
	requestBody := request.Body
	rawBody, err := io.ReadAll(requestBody)
	if err != nil {
		return []byte{}, err
	}
	defer request.Body.Close()
	return rawBody, nil
}

// verifySignature will verify request comes from slack.
// See also https://api.slack.com/authentication/verifying-requests-from-slack
func verifySignature(headers http.Header, rawBody []byte) bool {
	versionNumber := "v0"

	timestamp := headers.Get("X-Slack-Request-Timestamp")
	log.Println("TimeStamp: " + timestamp)

	expectedSignature := headers.Get("X-Slack-Signature")
	log.Println("Expected signature: " + expectedSignature)

	textToEncrypt := fmt.Sprintf(
		"%s:%s:%s",
		versionNumber,
		timestamp,
		string(rawBody),
	)
	log.Println("Text to encrypt: " + textToEncrypt)

	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	hash := hmac.New(sha256.New, []byte(signingSecret))
	hash.Write([]byte(textToEncrypt))
	encryptedResult := string(hex.EncodeToString(hash.Sum(nil)))
	encryptedResultWithVersionNumber := "v0=" + encryptedResult
	log.Println("EncryptedResult: " + encryptedResultWithVersionNumber)

	return hmac.Equal(
		[]byte(encryptedResultWithVersionNumber),
		[]byte(expectedSignature),
	)
}

// getTextParam returns the "text" string
// from the given rawBody payload.
// See also https://api.slack.com/interactivity/slash-commands#app_command_handling
func getTextParam(rawBody []byte) (string, error) {
	keyValuesSlice := strings.Split(string(rawBody), "&")
	keyValues := make(map[string]string)
	for _, keyValue := range keyValuesSlice {
		kV := strings.Split(keyValue, "=")
		keyValues[kV[0]] = kV[1]
	}
	log.Printf("KeyValues: %v\n", keyValues)

	commandValue := keyValues["command"]
	commandText := strings.Join(strings.Split(keyValues["text"], "+"), " ")
	if commandValue == "" && commandText == "" {
		return "", errors.New("Wrong query params provided")
	}

	return commandText, nil
}
