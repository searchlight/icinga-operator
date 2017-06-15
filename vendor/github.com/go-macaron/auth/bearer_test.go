package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/macaron.v1"
)

func Test_BearerAuth(t *testing.T) {
	recorder := httptest.NewRecorder()

	auth := "Bearer foobar"

	m := macaron.New()
	m.Use(Bearer("foobar"))
	m.Use(func(res http.ResponseWriter, req *http.Request, u User) {
		res.Write([]byte("hello " + u))
	})

	r, _ := http.NewRequest("GET", "foo", nil)

	m.ServeHTTP(recorder, r)

	if recorder.Code != 401 {
		t.Error("Response not 401")
	}

	if recorder.Body.String() == "hello " {
		t.Error("Auth block failed")
	}

	recorder = httptest.NewRecorder()
	r.Header.Set("Authorization", auth)
	m.ServeHTTP(recorder, r)

	if recorder.Code == 401 {
		t.Error("Response is 401")
	}

	if recorder.Body.String() != "hello " {
		t.Error("Auth failed, got: ", recorder.Body.String())
	}
}

func Test_BearerFuncAuth(t *testing.T) {
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
		encoded := "Bearer " + auth

		m := macaron.New()
		m.Use(BearerFunc(func(token string) bool {
			return valid
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
