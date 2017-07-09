package icinga

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func (ic *APIRequest) Do() *APIResponse {
	if ic.Err != nil {
		return &APIResponse{
			Err: ic.Err,
		}
	}
	ic.req.Header.Set("Accept", "application/json")

	if ic.userName != "" && ic.password != "" {
		ic.req.SetBasicAuth(ic.userName, ic.password)
	}

	ic.resp, ic.Err = ic.client.Do(ic.req)
	if ic.Err != nil {
		return &APIResponse{
			Err: ic.Err,
		}
	}

	ic.Status = ic.resp.StatusCode
	ic.ResponseBody, ic.Err = ioutil.ReadAll(ic.resp.Body)
	if ic.Err != nil {
		return &APIResponse{
			Err: ic.Err,
		}
	}
	return &APIResponse{
		Status:       ic.Status,
		ResponseBody: ic.ResponseBody,
	}
}

func (r *APIResponse) Into(to interface{}) (int, error) {
	if r.Err != nil {
		return r.Status, r.Err
	}
	err := json.Unmarshal(r.ResponseBody, to)
	if err != nil {
		return r.Status, err
	}
	return r.Status, nil
}

func (c *Client) newRequest(path string) *APIRequest {
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
	return &APIRequest{
		uri:      c.config.Endpoint + c.pathPrefix,
		client:   client,
		userName: c.config.BasicAuth.Username,
		password: c.config.BasicAuth.Password,
	}
}

func (ic *APIRequest) newRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	if strings.HasSuffix(urlStr, "/") {
		urlStr = strings.TrimRight(urlStr, "/")
	}

	return http.NewRequest(method, urlStr, body)
}
