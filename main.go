package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/coreos/go-semver/semver"
)

type PkgFile struct {
	Version string `json:"version"`
}

func getEnvDefault(key, default_value string) string {
    if val, found := os.LookupEnv(key); found {
        return val
    }
    return default_value
}

func main() {

	var semtags []*semver.Version
	var caltags []string

	orgName := os.Getenv("INPUT_ORG")
	pkgName := os.Getenv("INPUT_PACKAGE")
	verSys := getEnvDefault("INPUT_VERSION_SYSTEM", "SemVer")

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
	fmt.Println(fmt.Sprintf(`body:: %s`, body))
	if unmarshalErr != nil {
		log.Fatal(unmarshalErr)
	}

	if (verSys == "SemVer") {
		for _, tag := range pkg {
			matched, _ := regexp.MatchString(`.*\..*\..*`, tag.Version)
			if matched {
				semtags = append(semtags, semver.New(tag.Version))
			}
		}
		semver.Sort(semtags)

		if len(semtags) == 0 {
			log.Fatal(fmt.Sprintf(`Unable to find files for %s/%s`, orgName, pkgName))
		}

		fmt.Println(fmt.Sprintf(`::set-output name=version::%s`, semtags[len(semtags)-1]))

	} else {  // CalVer
		for _, tag := range pkg {
			matched, _ := regexp.MatchString(`.*\..*\..*`, tag.Version)
			if matched {
				caltags = append(caltags, tag.Version)
			}
		}
		sort.Strings(caltags)

		if len(caltags) == 0 {
			log.Fatal(fmt.Sprintf(`Unable to find files for %s/%s`, orgName, pkgName))
		}

		fmt.Println(fmt.Sprintf(`::set-output name=version::%s`, caltags[len(caltags)-1]))

	}
}
