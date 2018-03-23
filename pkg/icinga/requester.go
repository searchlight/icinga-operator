package icinga

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
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

	if c.config.CACert != nil {
		certs := x509.NewCertPool()
		certs.AppendCertsFromPEM(c.config.CACert)
		mTLSConfig.RootCAs = certs
	} else {
		mTLSConfig.InsecureSkipVerify = true
	}

	// ref: https://github.com/golang/go/blob/release-branch.go1.9/src/net/http/transport.go#L35
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       mTLSConfig,
	}
	client := &http.Client{Transport: tr}

	return &APIRequest{
		uri:      c.config.Endpoint + path,
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
