package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Idcf struct {
	ApiKey     string `json:"api_key"`
	DeletePath string `json:"delete_path"`
	Expired    int64  `json:"expired"`
	secretKey  string
	method     string
	uri        string
}

func (i *Idcf) url() string {
	return fmt.Sprintf("https://cdn.idcfcloud.com%s", i.uri)
}

func (i *Idcf) call(url string, expired time.Time) {
	i.Expired = expired.Unix()
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.Encode(&i)
	req, _ := http.NewRequest(i.method, url, &buf)
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
	return fmt.Sprintf("%s\n%s\n%s\n%d\n%s\n%s", i.method, i.ApiKey, i.secretKey, i.Expired, i.uri, string(request_body))
}

func (i *Idcf) signature() string {
	mac := hmac.New(sha256.New, []byte(i.secretKey))
	mac.Write([]byte(i.str()))

	return base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(mac.Sum(nil))))
}

func main() {
	fmt.Println("vim-go")

	apiKey := flag.String("api-key", "", "api key")
	secretKey := flag.String("secret-key", "", "secret key")
	deletePath := flag.String("delete-path", "", "delete path")
	day := flag.Int("day", 8, "extend day")
	flag.Parse()

	idcf := Idcf{
		ApiKey:     *apiKey,
		method:     "DELETE",
		DeletePath: *deletePath,
		secretKey:  *secretKey,
		uri:        "/api/v0/caches",
	}

	idcf.call(idcf.url(), time.Now().AddDate(0, 0, *day))
}
