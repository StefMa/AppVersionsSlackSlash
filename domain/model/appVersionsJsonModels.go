package model

type AppVersionsJsonResponse struct {
	AndroidApps []AppVersionsApp `json:"android"`
	IosApps     []AppVersionsApp `json:"ios"`
}

type AppVersionsApp struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Rating  string `json:"rating"`
	URL     string `json:"url"`
}
