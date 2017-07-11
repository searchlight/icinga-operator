package icinga

import (
	"bytes"
	"errors"
	"net/http"
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
	config     Config
	pathPrefix string
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

func (c *Client) Objects() *Client {
	c.pathPrefix = "/objects"
	return c
}

func (c *Client) Hosts(hostName string) *APIRequest {
	return c.newRequest("/hosts/" + hostName)
}

func (c *Client) HostGroups(hostName string) *APIRequest {
	return c.newRequest("/hostgroups/" + hostName)
}

func (c *Client) Service(hostName string) *APIRequest {
	return c.newRequest("/services/" + hostName)
}

func (c *Client) Actions(action string) *APIRequest {
	c.pathPrefix = ""
	return c.newRequest("/actions/" + action)
}

func (c *Client) Notifications(hostName string) *APIRequest {
	return c.newRequest("/notifications/" + hostName)
}

func (c *Client) Check() *APIRequest {
	c.pathPrefix = ""
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
		ic.Err = errors.New("Invalid request")
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
