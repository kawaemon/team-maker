package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type LineClient struct {
	ChannelAccessToken string
	OnMessage          func(string) string
}

const defaultPort = "8080"

func (c *LineClient) Start() error {
	port := os.Getenv("PORT")

	if port == "" {
		log.Printf("PORT env variable is not specified. Falling back to %s\n", defaultPort)
		port = defaultPort
	}

	if _, err := strconv.Atoi(port); err != nil {
		log.Printf("PORT env variable is not valid. Falling back to %s\n", defaultPort)
		port = defaultPort
	}

	http.HandleFunc("/", c.onRequest)

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Webhook listening at %s", addr)

	return http.ListenAndServe(addr, nil)
}

type lineEvent struct {
	Type       string      `json:"type"`
	ReplyToken string      `json:"replyToken"`
	Message    lineMessage `json:"message"`
}

type lineMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type lineRequest struct {
	Events []lineEvent `json:"events"`
}

type toLineRequest struct {
	Messages               []lineMessage `json:"messages"`
	ReplyToken             string        `json:"replyToken"`
	IsNotificationDisabled bool          `json:"notificationDisabled"`
}

const replyEndPoint = "https://api.line.me/v2/bot/message/reply"

func (c *LineClient) onRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	if r.ContentLength > 1000*1000*10 {
		log.Printf("Request Content Length is more than 10MB(%d). Ignoring.\n", r.ContentLength)
		return
	}

	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Failed to fetch request body: %s\n", err)
		return
	}

	deserData := lineRequest{}
	err = json.Unmarshal(data, &deserData)

	if err != nil {
		log.Printf("Failed to deserialize request: %s\n", err)
		return
	}

	events := deserData.Events

	if events == nil {
		return
	}

	for _, event := range events {
		if event.Type != "message" || event.Message.Type != "text" {
			continue
		}

		reply := strings.TrimSpace(c.OnMessage(event.Message.Text))

		if reply == "" {
			continue
		}

		err = c.replyToLine(event.ReplyToken, reply)
		if err != nil {
			log.Printf("Failed to reply to line: %s\n", err)
			continue
		}
	}
}

func (c *LineClient) replyToLine(replyToken string, msg string) error {
	reqStruct := toLineRequest{
		Messages: []lineMessage{
			{
				Type: "text",
				Text: msg,
			},
		},
		ReplyToken:             replyToken,
		IsNotificationDisabled: false,
	}

	reqJson, err := json.Marshal(&reqStruct)

	if err != nil {
		return fmt.Errorf("failed to serialize request to json: %w", err)
	}

	req, err := http.NewRequest("POST", replyEndPoint, bytes.NewBuffer(reqJson))

	if err != nil {
		return fmt.Errorf("failed to construct request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ChannelAccessToken))

	httpClient := http.Client{}
	response, err := httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("failed to send request to line reply api endpoint: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("LINE reply API responded with %s", response.Status)
	}

	return nil
}
