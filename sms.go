package main

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/jh-bate/intertidal/backend/platform"
	"github.com/jh-bate/intertidal/backend/store"
	"github.com/jh-bate/intertidal/flood"
)

const (
	//Health
	TIME     = "T"
	ACTIVITY = "A"
	BG       = "G"
	CARB     = "C"
	BASAL    = "L"
	BOLUS    = "S"
	NOTE     = "N"

	//Calcs
	LOG_LOW = "#LG"

	MMOLL = "mmol/L"
)

type (
	SmsApi interface {
		Load(userId string) (msgs []*Sms, err error)
		Send(sms SmsContent) error
	}
	Client struct {
		Messages  []*Sms
		raw       []TextData
		processed []interface{}
		api       SmsApi
	}
	TextData struct {
		text, date, device string
	}
	SmsRecived struct {
		Body, Date, From string
	}
)

func NewClient() *Client {
	return &Client{}
}

func newTextData(text, date, device string) TextData {
	return TextData{
		text:   text,
		date:   date,
		device: device,
	}
}

func (c *Client) AttachApi(api SmsApi) *Client {
	c.api = api
	return c
}

func (c *Client) Load() *Client {

	log.Println("loading sms messages")

	c.Messages, _ = c.api.Load("123")

	for i := range c.Messages {
		msg := messages.Messages[i]
		c.raw = append(c.raw, newTextData(msg.Body, msg.DateSent, msg.From))
	}

	c.transform()
	return c
}

func (c *Client) transform() {
	log.Println("transform text from sms client")

	for i := range c.raw {

		smsTxt := strings.Split(c.raw[i].text, " ")

	outer:
		for en := range smsTxt {

			log.Println("text ", smsTxt[en])

			switch {
			case strings.Index(strings.ToUpper(smsTxt[en]), BG) != -1:
				bg := strings.Split(smsTxt[en], BG)
				c.processed = append(c.processed, flood.MakeBg(bg[1], c.raw[i].date, c.raw[i].device))
				break
			case strings.Index(strings.ToUpper(smsTxt[en]), CARB) != -1:
				carb := strings.Split(smsTxt[en], CARB)
				c.processed = append(c.processed, flood.MakeCarb(carb[1], c.raw[i].date, c.raw[i].device))
				break
			case strings.Index(strings.ToUpper(smsTxt[en]), BASAL) != -1:
				basal := strings.Split(smsTxt[en], BASAL)
				c.processed = append(c.processed, flood.MakeBasal(basal[1], c.raw[i].date, c.raw[i].device))
				break
			case strings.Index(strings.ToUpper(smsTxt[en]), BOLUS) != -1:
				bolus := strings.Split(smsTxt[en], BOLUS)
				c.processed = append(c.processed, flood.MakeBolus(bolus[1], c.raw[i].date, c.raw[i].device))
				break
			case strings.Index(strings.ToUpper(smsTxt[en]), LOG_LOW) != -1:
				//hard code 'LOW'
				c.processed = append(c.processed, flood.MakeBg("3.9", c.raw[i].date, c.raw[i].device))
				break
			case strings.Index(strings.ToUpper(smsTxt[en]), ACTIVITY) != -1:
				log.Println("Will be an activity ", c.raw[i])
				break
			case strings.Index(strings.ToUpper(smsTxt[en]), NOTE) != -1:
				c.processed = append(c.processed, flood.MakeNote(smsTxt[en], c.raw[i].date, c.raw[i].device))
				break
			default:
				c.processed = append(c.processed, flood.MakeNote(c.raw[i].text, c.raw[i].date, c.raw[i].device))
				break outer
			}
		}
	}
	return
}

func (c *Client) StashLocal(key string, local store.Client) *Client {

	if len(c.processed) > 0 {

		log.Printf("to stash: [%v]", c.processed)

		err := local.StoreUserData(key, c.processed)

		if err != nil {
			log.Println("Error statshing data ", err)
		}
		return c
	}
	log.Println("No data to stash")
	return c
}

func (c *Client) StorePlatform(platform platform.Client) *Client {

	if len(c.processed) > 0 {

		err := platform.LoadInto(c.processed)

		if err != nil {
			log.Println("Error sending to platform ", err)
		}
	}
	log.Println("No data to send to the platform")
	return c
}
