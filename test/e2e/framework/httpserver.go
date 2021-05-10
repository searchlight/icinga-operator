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

package framework

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/gomega"
)

type Message struct {
	To   []string `json:"to,omitempty"`
	Body string   `json:"body,omitempty"`
}

func StartServer() *httptest.Server {
	var msg Message
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Write([]byte(msg.Body))
		case "POST":
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := json.Unmarshal(data, &msg); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}
	return httptest.NewServer(http.HandlerFunc(handler))
}

func (f *Framework) EventuallyHTTPServerResponse(serverURL string) GomegaAsyncAssertion {
	return Eventually(
		func() string {
			resp, err := http.Get(serverURL)
			Expect(err).NotTo(HaveOccurred())
			data, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())

			return string(data)
		},
		time.Minute*5,
		time.Second*5,
	)
}
