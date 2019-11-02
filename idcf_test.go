package idcf

import (
	"testing"
)

func TestSignature(t *testing.T) {
	idcf := Idcf{
		APIKey:     "apikey",
		Method:     "DELETE",
		DeletePath: "http://aaaaa/.*",
		SecretKey:  "secret",
		Expired:    12345678,
		URI:        "/api/v0/caches",
	}

	sign, err := idcf.signature()
	if err != nil {
		t.Error(err)
	}

	if sign != "NzE4NWFlMTMyMjMwMzVmMGFjM2ZhYjI2NmM2OWE2N2E1MTQwZWQ5MzA3NTY4Y2Q4ZGQ2OTU1MGI1NDk0ZTdjOQ==" {
		t.Errorf("signature:%s", sign)
	}
}
