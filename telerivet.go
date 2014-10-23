package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type (
	TelerivetClient struct {
		config     TelerivetConfig
		httpClient *http.Client
	}
	TelerivetConfig struct {
		ApiKey    string `json:"apiKey"`
		ProjectId string `json:"projectId"`
	}
	SmsContent struct {
		Text string `json:"content"`
		To   string `json:"to_number"`
	}
	Sms struct {
		Id   string `json:"id"`
		Text string `json:"content"`
		From string `json:"from_number"`
		To   string `json:"to_number"`
		Date int    `json:"time_sent"`
	}
)

const (
	base_url = "https://api.telerivet.com/v1/projects/"
)

var (
	loadMsgsUrl = base_url + "%s/contacts/%s/messages?direction=incoming&message_type=sms"
	sendMsgUrl  = base_url + "%s/messages/send"
)

func NewClient(config interface{}) *TelerivetClient {
	return &TelerivetClient{httpClient: &http.Client{}, config: config.(TelerivetConfig)}
}

func (c *TelerivetClient) Load(userId string) (msgs []*Sms, err error) {

	log.Println("loading from telerivet")

	url := fmt.Sprintf(loadMsgsUrl, c.config.ProjectId, userId)

	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(c.config.ApiKey)
	if resp, err := tc.httpClient.Do(req); err != nil {
		return nil, err
	} else {
		if resp.StatusCode == http.StatusOK {
			defer req.Body.Close()
			if err := json.NewDecoder(req.Body).Decode(&msgs); err != nil {
				log.Printf("Error after trying to load messages: %v", err)
				return err
			}
			return msgs, nil
		}
		return nil, errors.New("Issue loading messages: " + string(resp.StatusCode))
	}
}

func (c *TelerivetClient) Send(sms SmsContent) error {

	log.Println("sending via telerivet")

	url := fmt.Sprintf(sendMsgUrl, c.config.ProjectId)

	jsonSms, _ := json.Marshal(sms)
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(string(jsonSms)))
	req.SetBasicAuth(c.config.ApiKey)
	if resp, err := tc.httpClient.Do(req); err != nil {
		return err
	} else {
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		return errors.New("Issue loading messages: " + string(resp.StatusCode))
	}

	/*
			curl -s -u YOUR_API_KEY: \
		 "https://api.telerivet.com/v1/projects/PROJECT_ID/messages/send" \
		 -H "Content-Type: application/json" \
		 -d '{
		    "content": "hello world",
		    "to_number": "+16505550123"
		  }'
	*/

	return
}
