package gth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
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
func UpdateChecker(v, b string) {

	// GET request
	res, err := http.Get("https://api.github.com/repos/prvst/philosopher/releases")
	if err != nil {
		logrus.Warning("Can't check for updates, server unreachable: ", err)
	} else {

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logrus.Fatal(err)
		}

		var rel Releases
		err = json.Unmarshal(body, &rel)
		if err != nil {

			logrus.Warning("GitHub unreachable for the moment, can't check for versions right now.")

		} else {

			local := strings.Split(v, ".")
			local[0] = strings.Replace(local[0], "v", "", 1)

			remote := strings.Split(rel[0].TagName, ".")
			remote[0] = strings.Replace(remote[0], "v", "", 1)

			if remote[0] > local[0] {
				logrus.Warning("There is a new version of Philosopher available for download: https://github.com/prvst/philosopher/releases")
			}

			if (remote[0] == local[0]) && (remote[1] > local[1]) {
				logrus.Warning("There is a new version of Philosopher available for download: https://github.com/prvst/philosopher/releases")
			}

		}
	}

}
