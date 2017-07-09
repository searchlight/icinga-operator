package client

import (
	"bytes"
	"errors"
	"net/http"

	_ "github.com/appscode/searchlight/api/install"
	clientset "k8s.io/client-go/kubernetes"
)

type IcingaConfig struct {
	Endpoint  string
	BasicAuth struct {
		Username string
		Password string
	}
	CaCert []byte
}

type IcingaClient struct {
	config     *IcingaConfig
	pathPrefix string
}

type IcingaApiRequest struct {
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

type IcingaApiResponse struct {
	Err          error
	Status       int
	ResponseBody []byte
}

func newClient(icingaConfig *IcingaConfig) *IcingaClient {
	c := &IcingaClient{
		config: icingaConfig,
	}
	return c
}

func (c *IcingaClient) SetEndpoint(endpoint string) *IcingaClient {
	c.config.Endpoint = endpoint
	return c
}

func (c *IcingaClient) Objects() *IcingaClient {
	c.pathPrefix = "/objects"
	return c
}

func (c *IcingaClient) Hosts(hostName string) *IcingaApiRequest {
	return c.newIcingaRequest("/hosts/" + hostName)
}

func (c *IcingaClient) HostGroups(hostName string) *IcingaApiRequest {
	return c.newIcingaRequest("/hostgroups/" + hostName)
}

func (c *IcingaClient) Service(hostName string) *IcingaApiRequest {
	return c.newIcingaRequest("/services/" + hostName)
}

func (c *IcingaClient) Actions(action string) *IcingaApiRequest {
	c.pathPrefix = ""
	return c.newIcingaRequest("/actions/" + action)
}

func (c *IcingaClient) Notifications(hostName string) *IcingaApiRequest {
	return c.newIcingaRequest("/notifications/" + hostName)
}

func (c *IcingaClient) Check() *IcingaApiRequest {
	c.pathPrefix = ""
	return c.newIcingaRequest("")
}

func addUri(uri string, name []string) string {
	for _, v := range name {
		uri = uri + "!" + v
	}
	return uri
}

func (ic *IcingaApiRequest) Get(name []string, jsonBody ...string) *IcingaApiRequest {
	if len(jsonBody) == 0 {
		ic.req, ic.Err = ic.newRequest("GET", addUri(ic.uri, name), nil)
	} else if len(jsonBody) == 1 {
		ic.req, ic.Err = ic.newRequest("GET", addUri(ic.uri, name), bytes.NewBuffer([]byte(jsonBody[0])))
	} else {
		ic.Err = errors.New("Invalid request")
	}
	return ic
}

func (ic *IcingaApiRequest) Create(name []string, jsonBody string) *IcingaApiRequest {
	ic.req, ic.Err = ic.newRequest("PUT", addUri(ic.uri, name), bytes.NewBuffer([]byte(jsonBody)))
	return ic
}

func (ic *IcingaApiRequest) Update(name []string, jsonBody string) *IcingaApiRequest {
	ic.req, ic.Err = ic.newRequest("POST", addUri(ic.uri, name), bytes.NewBuffer([]byte(jsonBody)))
	return ic
}

func (ic *IcingaApiRequest) Delete(name []string, jsonBody string) *IcingaApiRequest {
	ic.req, ic.Err = ic.newRequest("DELETE", addUri(ic.uri, name), bytes.NewBuffer([]byte(jsonBody)))
	return ic
}

func (ic *IcingaApiRequest) Params(param map[string]string) *IcingaApiRequest {
	p := ic.req.URL.Query()
	for k, v := range param {
		p.Add(k, v)
	}
	ic.req.URL.RawQuery = p.Encode()
	return ic
}

func NewIcingaClient(kubeClient clientset.Interface, secretName, secretNamespace string) (*IcingaClient, error) {
	config, err := getIcingaConfig(kubeClient, secretName, secretNamespace)
	if err != nil {
		return nil, err
	}
	c := newClient(config)
	return c, nil
}
