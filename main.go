package main

import (
	"flag"
	"log"

	"github.com/jh-bate/intertidal-sms/clients"
	"github.com/jh-bate/intertidal/backend/platform"
	"github.com/jh-bate/intertidal/backend/store"
)

func loadFromSms(key, projectId string, stash *store.BoltClient) {
	log.Println("load from sms")

	tr := clients.NewTelerivetClient(clients.TelerivetConfig{ApiKey: key, ProjectId: projectId})

	sms := NewClient().
		AttachApi(api)

	p := platform.NewClient(
		&platform.Config{
			Auth:   "https://staging-api.tidepool.io/auth",
			Upload: "https://staging-uploads.tidepool.io/data",
		},
		"jamie@tidepool.org",
		"blip4life",
	)

	p.StashUserLocal(stash)

	sms.Load().StorePlatform(p)
}

func main() {

	key := flag.String("k", "", "api key")
	projId := flag.String("p", "", "projectId")

	flag.Parse()

	stash := store.NewBoltClient()

	loadFromSms(*key, *projId, stash)
}
