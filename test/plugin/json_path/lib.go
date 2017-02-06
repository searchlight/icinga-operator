package json_path

import "github.com/appscode/go/net/httpclient"

type GithubOrg struct {
	PublicRepos int `json:"public_repos"`
}

func GetPublicRepoNumber() (int, error) {
	httpClient := httpclient.Default().WithBaseURL("https://api.github.com")

	var githubOrg GithubOrg
	_, err := httpClient.Call("GET", "/orgs/appscode", nil, &githubOrg, false)
	if err != nil {
		return 0, err
	}

	return githubOrg.PublicRepos, nil
}
