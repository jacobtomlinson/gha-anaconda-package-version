package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/coreos/go-semver/semver"
)

type PkgFile struct {
	Version string `json:"version"`
}

func main() {

	orgName := os.Getenv("INPUT_ORG")
	pkgName := os.Getenv("INPUT_PACKAGE")

	url := fmt.Sprintf(`https://api.anaconda.org/package/%s/%s/files`, orgName, pkgName)

	dhClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "gha-anaconda-package-version")

	res, getErr := dhClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	pkg := make([]PkgFile, 0)
	unmarshalErr := json.Unmarshal(body, &pkg)
	if unmarshalErr != nil {
		log.Fatal(unmarshalErr)
	}

	var tags []semver.Version
	for _, tag := range pkg {
		matched, _ := regexp.MatchString(`.*\..*\..*`, tag.Version)
		if matched {
			tags = append(tags, *semver.New(tag.Version))
		}
	}

	if len(tags) == 0 {
		log.Fatal(fmt.Sprintf(`Unable to find files for %s/%s`, orgName, pkgName))
	}

	semver.Sort(tags)
	fmt.Println(fmt.Sprintf(`::set-output name=version::%s`, tags[len(tags)-1]))
}
