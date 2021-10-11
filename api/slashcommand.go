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
	"stefma.guru/appVersionsSlackSlash/domain"
	"strings"
)

func HandleSlashCommand(
	writer http.ResponseWriter,
	request *http.Request,
) {
	rawBody, err := readBody(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if valid := verifySignature(request.Header, rawBody); !valid {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	text, err := getTextParam(rawBody)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.CreateDatabase()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	commandAndArguments := strings.Split(text, " ")
	switch strings.TrimSpace(commandAndArguments[0]) {
	case "get":
		domain.Get(writer, request, db)
		break
	case "add":
		domain.Add(writer, request, db, commandAndArguments[1:])
		break
	case "remove":
		domain.Remove(writer, request, db, commandAndArguments[1:])
		break
	default:
		http.Error(writer, "Unknown command. Only 'add', 'remove' & 'add' allowed", http.StatusBadRequest)
	}
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
	slashCommandText := keyValues["text"]
	if commandValue == "" && slashCommandText == "" {
		return "", errors.New("Wrong query params provided")
	}
	return slashCommandText, nil
}
