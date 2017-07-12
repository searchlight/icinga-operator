package json_path

import (
	"fmt"

	"github.com/appscode/go/net/httpclient"
	"github.com/appscode/searchlight/test/plugin"
)

type GithubOrg struct {
	PublicRepos int `json:"public_repos"`
}

func getPublicRepoNumber(url, uri string) (int, error) {
	httpClient := httpclient.Default().WithBaseURL(url)

	var githubOrg GithubOrg
	_, err := httpClient.Call("GET", uri, nil, &githubOrg, false)
	if err != nil {
		return 0, err
	}

	return githubOrg.PublicRepos, nil
}

func GetTestData() ([]plugin.TestData, error) {
	url := "https://api.github.com"
	uri := "/orgs/appscode"

	repoNumber, err := getPublicRepoNumber(url, uri)
	if err != nil {
		return nil, err
	}

	testDataList := []plugin.TestData{
		{
			Data: map[string]interface{}{
				"URL":     url + uri,
				"Query":   ".",
				"Warning": fmt.Sprintf(`.public_repos!=%v`, repoNumber),
			},
			ExpectedIcingaState: 0,
		},
		{
			Data: map[string]interface{}{
				"URL":     url + uri,
				"Query":   ".",
				"Warning": fmt.Sprintf(`.public_repos==%v`, repoNumber),
			},
			ExpectedIcingaState: 1,
		},
		{
			Data: map[string]interface{}{
				"URL":      url + uri,
				"Query":    ".",
				"Warning":  fmt.Sprintf(`.public_repos==%v`, repoNumber-1),
				"Critical": fmt.Sprintf(`.public_repos==%v`, repoNumber),
			},
			ExpectedIcingaState: 2,
		},
		{
			Data: map[string]interface{}{
				"URL":     url + uri + "fake",
				"Query":   ".",
				"Warning": fmt.Sprintf(`.public_repos==%v`, repoNumber-1),
			},
			ExpectedIcingaState: 3,
		},
	}

	return testDataList, nil
}
