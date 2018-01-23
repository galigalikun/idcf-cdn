package main

import (
	"flag"
	"log"
	"time"

	idcf "github.com/galigalikun/idcf-cdn"
)

func main() {

	var (
		apiKey     string
		secretKey  string
		deletePath string
		day        int
	)

	flag.StringVar(&apiKey, "api-key", "", "api key")
	flag.StringVar(&secretKey, "secret-key", "", "secret key")
	flag.StringVar(&deletePath, "delete-path", "", "delete path")
	flag.IntVar(&day, "day", 8, "extend day")
	flag.Parse()

	idcf := idcf.Idcf{
		ApiKey:     apiKey,
		Method:     "DELETE",
		DeletePath: deletePath,
		SecretKey:  secretKey,
		Uri:        "/api/v0/caches",
	}

	err := idcf.Call(time.Now().AddDate(0, 0, day))
	if err != nil {
		log.Fatal(err)
	}
}
