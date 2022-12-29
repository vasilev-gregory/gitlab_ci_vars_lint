package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type varDoc struct {
	gitlabVersion string
	runner        string
	descr         string
}

var (
	reVarFromDoc = regexp.MustCompile("(?m)^\\|\\s*`([^`]*)`\\s*\\|([^|]*)\\|([^|]*)\\|([^|]*)\\|")
	gitlabVars   = make(map[string]varDoc)
)

func parseGilabPredefinedVars() {
	url := "https://gitlab.com/gitlab-org/gitlab/-/raw/master/doc/ci/variables/predefined_variables.md"
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)

	for _, match := range reVarFromDoc.FindAllStringSubmatch(string(body), -1) {
		var doc varDoc
		doc.gitlabVersion = strings.TrimSpace(match[2])
		doc.runner = strings.TrimSpace(match[3])
		doc.descr = strings.TrimSpace(match[4])
		gitlabVars[strings.TrimSpace(match[1])] = doc
		//fmt.Printf("gitlabVars[%s] = %#v\n", strings.TrimSpace(match[1]), doc)
	}
}
