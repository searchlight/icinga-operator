package client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func (ic *IcingaApiRequest) Do() *IcingaApiResponse {
	if ic.Err != nil {
		return &IcingaApiResponse{
			Err: ic.Err,
		}
	}
	ic.req.Header.Set("Accept", "application/json")

	if ic.userName != "" && ic.password != "" {
		ic.req.SetBasicAuth(ic.userName, ic.password)
	}

	ic.resp, ic.Err = ic.client.Do(ic.req)
	if ic.Err != nil {
		return &IcingaApiResponse{
			Err: ic.Err,
		}
	}

	ic.Status = ic.resp.StatusCode
	ic.ResponseBody, ic.Err = ioutil.ReadAll(ic.resp.Body)
	if ic.Err != nil {
		return &IcingaApiResponse{
			Err: ic.Err,
		}
	}
	return &IcingaApiResponse{
		Status:       ic.Status,
		ResponseBody: ic.ResponseBody,
	}
}

func (r *IcingaApiResponse) Into(to interface{}) (int, error) {
	if r.Err != nil {
		return r.Status, r.Err
	}
	err := json.Unmarshal(r.ResponseBody, to)
	if err != nil {
		return r.Status, err
	}
	return r.Status, nil
}

func (c *IcingaClient) newIcingaRequest(path string) *IcingaApiRequest {
	mTLSConfig := &tls.Config{}

	if c.config.CaCert != nil {
		certs := x509.NewCertPool()
		certs.AppendCertsFromPEM(c.config.CaCert)
		mTLSConfig.RootCAs = certs
	} else {
		mTLSConfig.InsecureSkipVerify = true
	}

	tr := &http.Transport{
		TLSClientConfig: mTLSConfig,
	}
	client := &http.Client{Transport: tr}

	c.pathPrefix = c.pathPrefix + path
	return &IcingaApiRequest{
		uri:      c.config.Endpoint + c.pathPrefix,
		client:   client,
		userName: c.config.BasicAuth.Username,
		password: c.config.BasicAuth.Password,
	}
}

func (ic *IcingaApiRequest) newRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	if strings.HasSuffix(urlStr, "/") {
		urlStr = strings.TrimRight(urlStr, "/")
	}

	return http.NewRequest(method, urlStr, body)
}
