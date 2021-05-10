/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package icinga

import (
	"bytes"
	"net/http"

	"github.com/pkg/errors"
)

type Config struct {
	Endpoint  string
	BasicAuth struct {
		Username string
		Password string
	}
	CACert []byte
}

type Client struct {
	config Config
}

type APIRequest struct {
	client *http.Client

	uri      string
	suffix   string
	params   map[string]string
	userName string
	password string
	verb     string

	Err  error
	req  *http.Request
	resp *http.Response

	Status       int
	ResponseBody []byte
}

type APIResponse struct {
	Err          error
	Status       int
	ResponseBody []byte
}

func NewClient(cfg Config) *Client {
	return &Client{config: cfg}
}

func (c *Client) SetEndpoint(endpoint string) *Client {
	c.config.Endpoint = endpoint
	return c
}

func (c *Client) Hosts(hostName string) *APIRequest {
	return c.newRequest("/objects/hosts/" + hostName)
}

func (c *Client) HostGroups(hostName string) *APIRequest {
	return c.newRequest("/objects/hostgroups/" + hostName)
}

func (c *Client) Service(hostName string) *APIRequest {
	return c.newRequest("/objects/services/" + hostName)
}

func (c *Client) Notifications(hostName string) *APIRequest {
	return c.newRequest("/objects/notifications/" + hostName)
}

func (c *Client) Actions(action string) *APIRequest {
	return c.newRequest("/actions/" + action)
}

func (c *Client) Check() *APIRequest {
	return c.newRequest("")
}

func addUri(uri string, name []string) string {
	for _, v := range name {
		uri = uri + "!" + v
	}
	return uri
}

func (ic *APIRequest) Get(name []string, jsonBody ...string) *APIRequest {
	if len(jsonBody) == 0 {
		ic.req, ic.Err = ic.newRequest("GET", addUri(ic.uri, name), nil)
	} else if len(jsonBody) == 1 {
		ic.req, ic.Err = ic.newRequest("GET", addUri(ic.uri, name), bytes.NewBuffer([]byte(jsonBody[0])))
	} else {
		ic.Err = errors.New("invalid request")
	}
	return ic
}

func (ic *APIRequest) Create(name []string, jsonBody string) *APIRequest {
	ic.req, ic.Err = ic.newRequest("PUT", addUri(ic.uri, name), bytes.NewBuffer([]byte(jsonBody)))
	return ic
}

func (ic *APIRequest) Update(name []string, jsonBody string) *APIRequest {
	ic.req, ic.Err = ic.newRequest("POST", addUri(ic.uri, name), bytes.NewBuffer([]byte(jsonBody)))
	return ic
}

func (ic *APIRequest) Delete(name []string, jsonBody string) *APIRequest {
	ic.req, ic.Err = ic.newRequest("DELETE", addUri(ic.uri, name), bytes.NewBuffer([]byte(jsonBody)))
	return ic
}

func (ic *APIRequest) Params(param map[string]string) *APIRequest {
	p := ic.req.URL.Query()
	for k, v := range param {
		p.Add(k, v)
	}
	ic.req.URL.RawQuery = p.Encode()
	return ic
}
