package gth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"philosopher/lib/msg"

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
			msg.Custom(e, "fatal")
		}

		var rel Releases
		e = json.Unmarshal(body, &rel)
		if e != nil {

			logrus.Warning("GitHub unreachable for the moment, can't check for versions right now.")

		} else {

			local := strings.Split(v, ".")
			local[0] = strings.Replace(local[0], "v", "", 1)

			remote := strings.Split(rel[0].TagName, ".")
			remote[0] = strings.Replace(remote[0], "v", "", 1)

			outdatedMajorVersion := remote[0] > local[0]
			outdatedMinorVersion := (remote[0] == local[0]) && (remote[1] > local[1])
			outdated := outdatedMajorVersion || outdatedMinorVersion

			if outdated {
				logrus.Warning("There is a new version of Philosopher available for download: https://github.com/Nesvilab/philosopher/releases")
			}

		}
	}

}
