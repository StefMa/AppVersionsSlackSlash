package domain

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"stefma.guru/appVersionsSlackSlash/database"
	"strings"
)

func Get(
	w http.ResponseWriter,
	r *http.Request,
	db database.Database,
) {
	storedData, err := db.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	urlToAppVersions := generaterUrl(storedData["android"], storedData["ios"])
	shortUrl, err := generateShortUrl(urlToAppVersions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(shortUrl))
}

func generaterUrl(androidAppIds []string, iosAppIds []string) string {
	params := url.Values{}
	if len(androidAppIds) > 0 {
		params.Add("android", strings.Join(androidAppIds, ","))
	}
	if len(iosAppIds) > 0 {
		params.Add("ios", strings.Join(iosAppIds, ","))
	}
	appVersionsLink := "https://appversions.vercel.app?" + url.QueryEscape(params.Encode())
	return appVersionsLink
}

func generateShortUrl(longUrl string) (string, error) {
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
	if err != nil {
		return "", err
	}
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jsonResponse); err != nil {
		return "", err
	}
	return jsonResponse["shortLink"].(string), nil
}
