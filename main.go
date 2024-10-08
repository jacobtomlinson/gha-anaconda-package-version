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
	"strconv"
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

func writeOutput(key, value string) error {
	output_filename := os.Getenv("GITHUB_OUTPUT")
	f, err := os.OpenFile(output_filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf(`%s=%s`, key, value)); err != nil {
		return err
	}
	return nil
}

func main() {

	var semtags []*semver.Version
	var caltags []string

	orgName := os.Getenv("INPUT_ORG")
	pkgName := os.Getenv("INPUT_PACKAGE")
	verSys := getEnvDefault("INPUT_VERSION_SYSTEM", "SemVer")
	retries, err := strconv.Atoi(getEnvDefault("INPUT_RETRIES", "3"))

	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf(`https://api.anaconda.org/package/%s/%s/files`, orgName, pkgName)

	dhClient := http.Client{
		Timeout: time.Second * 4, // Maximum of 4 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "gha-anaconda-package-version")
	var res *http.Response
	getErr := error(nil)
	for i := 1; i <= retries; i++ {
		res, getErr = dhClient.Do(req)
		if getErr == nil {
			break
		}
	}
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	pkg := make([]PkgFile, 0)
	unmarshalErr := json.Unmarshal(body, &pkg)
	//fmt.Println(fmt.Sprintf(`body:: %s`, body))
	if unmarshalErr != nil {
		log.Fatal(unmarshalErr)
	}

	if verSys == "SemVer" {
		for _, tag := range pkg {
			matched, _ := regexp.MatchString(`.*\..*\..*`, tag.Version)
			if matched {
				version, semverErr := semver.NewVersion(tag.Version)
				if semverErr == nil {
					semtags = append(semtags, version)
				} else {
					fmt.Println("incompatible semver found:", tag.Version, semverErr)
				}
			}
		}
		semver.Sort(semtags)

		if len(semtags) == 0 {
			log.Fatal(fmt.Sprintf(`Unable to find files for %s/%s`, orgName, pkgName))
		}
		writeErr := writeOutput("version", fmt.Sprintf(`%s`, semtags[len(semtags)-1]))
		if writeErr != nil {
			log.Fatal(readErr)
		}

	} else { // CalVer
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
		writeErr := writeOutput("version", fmt.Sprintf(`%s`, caltags[len(caltags)-1]))
		if writeErr != nil {
			log.Fatal(readErr)
		}
	}
}
