package gth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Nesvilab/philosopher/lib/msg"

	"github.com/sirupsen/logrus"
)

// Release information from GitHub
type Release struct {
	URL             string `json:"url"`
	AssetsURL       string `json:"assets_url"`
	UploadURL       string `json:"upload_url"`
	HTMLURL         string `json:"html_url"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	ID              int    `json:"id"`
	Draft           bool   `json:"draft"`
}

// Releases is a list of Release
type Releases []Release

// UpdateChecker reads GitHub API and reports if there is a new version available
func UpdateChecker(v, b string) {

	// GET request
	res, e := http.Get("https://api.github.com/repos/prvst/philosopher/releases")
	if e != nil {
		msg.Custom(errors.New("can't check for updates, server unreachable"), "warning")
	} else {

		body, e := ioutil.ReadAll(res.Body)
		if e != nil {
			msg.Custom(e, "error")
		}

		var rel Releases
		e = json.Unmarshal(body, &rel)
		if e != nil {

			logrus.Warning("GitHub unreachable for the moment, can't check for versions right now.")

		} else {

			local := strings.Replace(v, ".", "", -1)
			local = strings.Replace(local, "v", "", 1)

			remote := strings.Replace(rel[0].TagName, ".", "", -1)
			remote = strings.Replace(remote, "v", "", 1)

			if remote > local {
				logrus.Warning("There is a new version of Philosopher available for download: https://github.com/Nesvilab/philosopher/releases")
			}

		}
	}
}
