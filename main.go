package main

import (
	"flag"
	"log"

	"github.com/jh-bate/intertidal/backend/platform"
	"github.com/jh-bate/intertidal/backend/store"
)

func loadFromSms(token string, stash *store.BoltClient) {
	log.Println("load from trackthis")

	tt := trackthis.NewClient()
	p := platform.NewClient(
		&platform.Config{
			Auth:   "https://staging-api.tidepool.io/auth",
			Upload: "https://staging-uploads.tidepool.io/data",
		},
		"jamie@tidepool.org",
		"blip4life",
	)

	p.StashUserLocal(stash)

	tt.Init(trackthis.Config{AuthToken: token}).
		Load().
		StorePlatform(p).
		StashLocal(p.User.Token, stash)
}

func main() {

	authPtr := flag.String("t", "", "auth token for source")
	//destPtr := flag.String("d", "stash", "where the data will be put")

	flag.Parse()

	stash := store.NewBoltClient()

	loadFromSms(*authPtr, stash)
}
