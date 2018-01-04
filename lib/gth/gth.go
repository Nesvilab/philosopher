package gth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/met"
)

// Release information from GitHub
type Release struct {
	URL             string `json:"url"`
	AssetsURL       string `json:"assets_url"`
	UploadURL       string `json:"upload_url"`
	HTMLURL         string `json:"html_url"`
	ID              int    `json:"id"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Draft           bool   `json:"draft"`
}

// Releases is a list of Release
type Releases []Release

// UpdateChecker reads GitHub API and reports if there is a new version available
func UpdateChecker() {

	// GET request
	res, err := http.Get("https://api.github.com/repos/prvst/philosopher/releases")
	if err != nil {
		logrus.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Fatal(err)
	}

	var rel Releases
	err = json.Unmarshal(body, &rel)
	if err != nil {
		logrus.Fatal(err)
	}

	if rel[0].TagName != met.GetVersion() {
		logrus.Warning("There is a new version of Philosopher available for download: https://github.com/prvst/philosopher/releases")
	}

}
