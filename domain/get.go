package domain

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"stefma.guru/appVersionsSlackSlash/database"
)

func get(
	db database.Database,
	operatingSystem string,
) (string, error) {
	switch operatingSystem {
	case "android":
		break
	case "ios":
		break
	default:
		if operatingSystem != "" {
			log.Printf("Not supported operatingSystem '%s'. Fallback to 'all'", operatingSystem)
		}
		operatingSystem = "all"
	}
	storedData, err := db.Get(operatingSystem)
	if err != nil {
		return "", err
	}
	urlToAppVersions := generaterUrl(storedData["android"], storedData["ios"])
	shortUrl, err := generateShortUrl(urlToAppVersions)
	if err != nil {
		return "", err
	}
	return shortUrl, nil
}

func generaterUrl(androidAppIds []string, iosAppIds []string) string {
	params := url.Values{}
	if len(androidAppIds) > 0 {
		params.Add("android", strings.Join(androidAppIds, ","))
	}
	if len(iosAppIds) > 0 {
		params.Add("ios", strings.Join(iosAppIds, ","))
	}
	appVersionsLink := AppVersioBaseUrl + "/lookup?" + url.QueryEscape(params.Encode())
	return appVersionsLink
}

func generateShortUrl(longUrl string) (string, error) {
	log.Println("LongUrl: " + longUrl)
	apiKey := os.Getenv("FIREBASE_DYNAMIC_LINKS_API_KEY")
	url := "https://firebasedynamiclinks.googleapis.com/v1/shortLinks?key=" + apiKey
	shortUrlDomain := os.Getenv("FIREBASE_DYNAMIC_LINKS_DOMAIN")
	var jsonStr = []byte(`
    {
      "longDynamicLink": "` + shortUrlDomain + `?link=` + longUrl + `",
      "suffix": {
        "option": "SHORT"
      }
    }
  `)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	log.Println("Body: " + string(bodyBytes))
	if err != nil {
		return "", err
	}
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jsonResponse); err != nil {
		return "", err
	}
	return jsonResponse["shortLink"].(string), nil
}
