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

type Idcf struct {
	ApiKey     string `json:"api_key"`
	DeletePath string `json:"delete_path"`
	Expired    int64  `json:"expired"`
	SecretKey  string
	Method     string
	Uri        string
}

func (i *Idcf) url() string {
	return fmt.Sprintf("https://cdn.idcfcloud.com%s", i.Uri)
}

func (i *Idcf) Call(expired time.Time) {
	i.Expired = expired.Unix()
	request_body, _ := json.Marshal(i)
	req, _ := http.NewRequest(i.Method, i.url(), bytes.NewBuffer(request_body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("expired", fmt.Sprintf("%d", i.Expired))
	req.Header.Set("signature", i.signature())

	client := new(http.Client)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
}

func (i *Idcf) str() string {
	request_body, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s\n%s\n%s\n%d\n%s\n%s", i.Method, i.ApiKey, i.SecretKey, i.Expired, i.Uri, string(request_body))
}

func (i *Idcf) signature() string {
	mac := hmac.New(sha256.New, []byte(i.SecretKey))
	mac.Write([]byte(i.str()))

	return base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(mac.Sum(nil))))
}
