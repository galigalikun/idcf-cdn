package idcf

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Idcf struct
type Idcf struct {
	APIKey     string `json:"api_key"`
	DeletePath string `json:"delete_path"`
	Expired    int64  `json:"expired"`
	SecretKey  string
	Method     string
	URI        string
}

func (i *Idcf) url() string {
	return fmt.Sprintf("https://cdn.idcfcloud.com%s", i.URI)
}

// Call func Call
func (i *Idcf) Call(expired time.Time) error {
	i.Expired = expired.Unix()
	body, err := json.Marshal(i)
	if err != nil {
		return err
	}
	sign, err := i.signature()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(i.Method, i.url(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("expired", fmt.Sprintf("%d", i.Expired))
	req.Header.Set("signature", sign)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	byteArray, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(string(byteArray))
	}

	return nil
}

func (i *Idcf) signature() (string, error) {
	mac := hmac.New(sha256.New, []byte(i.SecretKey))

	body, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	mac.Write([]byte(fmt.Sprintf("%s\n%s\n%s\n%d\n%s\n%s", i.Method, i.APIKey, i.SecretKey, i.Expired, i.URI, string(body))))

	return base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(mac.Sum(nil)))), nil
}
