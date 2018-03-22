// +build ignore
package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	buffer "github.com/aybabtme/buffer-go"
)

func main() {
	cfg := &buffer.ClientConfig{}
	flag.StringVar(&cfg.ClientID, "ClientID", os.Getenv("BUFFER_CLIENTID"), "")
	flag.StringVar(&cfg.ClientSecret, "ClientSecret", os.Getenv("BUFFER_CLIENTSECRET"), "")
	flag.StringVar(&cfg.RedirectURL, "RedirectURL", os.Getenv("BUFFER_REDIRECTURL"), "")
	flag.StringVar(&cfg.AccessToken, "AccessToken", os.Getenv("BUFFER_ACCESSTOKEN"), "")
	flag.Parse()

	client, err := buffer.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profiles, err := client.Profiles().List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(profiles); err != nil {
		log.Fatal(err)
	}
}
