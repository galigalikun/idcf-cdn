package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSignature(t *testing.T) {
	idcf := Idcf{
		ApiKey:     "apikey",
		method:     "DELETE",
		DeletePath: "http://aaaaa/.*",
		secretKey:  "secret",
		Expired:    12345678,
		uri:        "/api/v0/caches",
	}

	str := `DELETE
apikey
secret
12345678
/api/v0/caches
{"api_key":"apikey","delete_path":"http://aaaaa/.*","expired":12345678}`

	if idcf.str() != str {
		t.Errorf("str:%s", idcf.str())
	}

	if idcf.signature() != "ODIwZDc0YTBhN2VmMmMyYzNkMGIwNGZlMzAwZGM1YmJjYzVkYmI5NzU0ZWE5NmZjZjkxNDg2NmE4NGFmYmIyMg==" {
		t.Errorf("signature:%s", idcf.signature())
	}
}

func TestCall(t *testing.T) {
	bodyTest := `{"api_key":"apikey","delete_path":"http://aaaaa/.*","expired":1371646123}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("signature") != "ZDZmNDMyNWJjYzkyMmRiZmUwODk1ZGE2ZDFmZWI1NTI4ZGUwMmZlOGY4OWM1NWI5ZmRjNjhiN2QxYzFiNTUzZQ==" {
			t.Errorf("signature:%s", r.Header.Get("signature"))
		}
		if r.Header.Get("expired") != "1371646123" {
			t.Errorf("expired:%s", r.Header.Get("expired"))
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		if string(body) != bodyTest {
			t.Errorf("body:%s", string(body))
		}
	}))

	defer ts.Close()

	idcf := Idcf{
		ApiKey:     "apikey",
		method:     "DELETE",
		DeletePath: "http://aaaaa/.*",
		secretKey:  "secret",
		uri:        "/api/v0/caches",
	}

	idcf.call(ts.URL, time.Unix(1371646123, 0))
}
