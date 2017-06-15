package cli

import (
	"strings"

	api "github.com/appscode/api/health"
	_env "github.com/appscode/go/env"
)

type Auth struct {
	ApiServer string           `json:"apiserver,omitempty"`
	Env       _env.Environment `json:"env,omitempty"`
	TeamId    string           `json:"team_id,omitempty"`
	UserName  string           `json:"username,omitempty"`
	Email     string           `json:"email,omitempty"`
	Token     string           `json:"token,omitempty"`
	Phid      string           `json:"phid,omitempty"`
	Settings  struct {
		CollectAnalytics bool   `json:"collect_analytics,omitempty"`
		TimeZone         string `json:"time_zone,omitempty"`
		TimeFormat       string `json:"time_format,omitempty"`
		DateFormat       string `json:"date_format,omitempty"`
	} `json:"settings,omitempty"`
	Network struct {
		ClusterUrls api.URLBase `json:"cluster_urls,omitempty"`
		PublicUrls  api.URLBase `json:"public_urls,omitempty"`
		TeamUrls    api.URLBase `json:"team_urls,omitempty"`
	} `json:"network,omitempty"`
}

func (a *Auth) TeamAddr() string {
	if a.Env.IsHosted() {
		return a.TeamId + "." + a.Network.TeamUrls.BaseAddr
	} else {
		return a.Network.TeamUrls.BaseAddr
	}
}

func (s *Auth) TeamURL(trails ...string) string {
	return strings.TrimRight(s.TeamEndpoint()+"/"+strings.Join(trails, "/"), "/")
}

func (a *Auth) TeamEndpoint() string {
	return a.Network.TeamUrls.Scheme + "://" + a.TeamAddr()
}

func NewAnonAUth() *Auth {
	a := &Auth{
		ApiServer: _env.ProdApiServer,
		Env:       _env.Prod,
	}
	a.Settings.CollectAnalytics = false
	a.Network.TeamUrls = api.URLBase{
		Scheme:   "https",
		BaseAddr: "appscode.com",
	}
	return a
}
