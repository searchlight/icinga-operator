package framework

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/gomega"
)

const (
	HTTPServerPort = "7033"
)

type Message struct {
	To   []string `json:"to,omitempty"`
	Body string   `json:"body,omitempty"`
}

func GetServer() *http.Server {
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
	return &http.Server{
		Addr:         ":" + HTTPServerPort,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      http.HandlerFunc(handler),
	}
}

func (f *Framework) EventuallyHTTPServerResponse() GomegaAsyncAssertion {
	return Eventually(
		func() string {
			resp, err := http.Get("http://127.0.0.1:" + HTTPServerPort)
			if err != nil {
				return err.Error()
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err.Error()
			}

			return string(data)
		},
		time.Minute*5,
		time.Second*5,
	)
}
