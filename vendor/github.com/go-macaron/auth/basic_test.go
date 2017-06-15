package auth

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/macaron.v1"
)

func Test_BasicAuth(t *testing.T) {
	recorder := httptest.NewRecorder()

	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte("foo:bar"))

	m := macaron.New()
	m.Use(Basic("foo", "bar"))
	m.Use(func(res http.ResponseWriter, req *http.Request, u User) {
		res.Write([]byte("hello " + u))
	})

	r, _ := http.NewRequest("GET", "foo", nil)

	m.ServeHTTP(recorder, r)

	if recorder.Code != 401 {
		t.Error("Response not 401")
	}

	if recorder.Body.String() == "hello foo" {
		t.Error("Auth block failed")
	}

	recorder = httptest.NewRecorder()
	r.Header.Set("Authorization", auth)
	m.ServeHTTP(recorder, r)

	if recorder.Code == 401 {
		t.Error("Response is 401")
	}

	if recorder.Body.String() != "hello foo" {
		t.Error("Auth failed, got: ", recorder.Body.String())
	}
}

func Test_BasicFuncAuth(t *testing.T) {
	for auth, valid := range map[string]bool{
		"foo:spam":       true,
		"bar:spam":       true,
		"foo:eggs":       false,
		"bar:eggs":       false,
		"baz:spam":       false,
		"foo:spam:extra": false,
		"dummy:":         false,
		"dummy":          false,
		"":               false,
	} {
		recorder := httptest.NewRecorder()
		encoded := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

		m := macaron.New()
		m.Use(BasicFunc(func(username, password string) bool {
			return (username == "foo" || username == "bar") && password == "spam"
		}))
		m.Use(func(res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("hello"))
		})

		r, _ := http.NewRequest("GET", "foo", nil)

		m.ServeHTTP(recorder, r)

		if recorder.Code != 401 {
			t.Error("Response not 401, params:", auth)
		}

		if recorder.Body.String() == "hello" {
			t.Error("Auth block failed, params:", auth)
		}

		recorder = httptest.NewRecorder()
		r.Header.Set("Authorization", encoded)
		m.ServeHTTP(recorder, r)

		if valid && recorder.Code == 401 {
			t.Error("Response is 401, params:", auth)
		}
		if !valid && recorder.Code != 401 {
			t.Error("Response not 401, params:", auth)
		}

		if valid && recorder.Body.String() != "hello" {
			t.Error("Auth failed, got: ", recorder.Body.String(), "params:", auth)
		}
		if !valid && recorder.Body.String() == "hello" {
			t.Error("Auth block failed, params:", auth)
		}
	}
}
