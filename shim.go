package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

type IncomingEvent struct {
	RecordType    string
	Type          string
	TypeCode      int
	Name          string
	Tag           string
	MessageStream string
	Description   string
	Email         string
	From          string
	BouncedAt     string
}

type SlackAlert struct {
	Text string `json:"text"` // Slack expects a lowercase key name
}

var SlackAddress string

func main() {
	flag.StringVar(&SlackAddress, "url", "",
		"The Slack API endpoint to post alerts to.")

	flag.Parse()

	if SlackAddress == "" {
		log.Panic("No address provided to post alerts to.  Unable to continue.")
	}

	fmt.Println("Event handler running!")
	http.HandleFunc("/event", handleEvent)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusBadRequest)
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusBadRequest)
	}
	var content IncomingEvent
	err := json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		http.Error(w, "Error parsing request", http.StatusBadRequest)
	}
	if content.Type == "SpamNotification" {
		go sendEvent(content)
	}
	// everything worked fine with the call, 200 is automatically returned.
}

func sendEvent(event IncomingEvent) {
	eventBody, err := json.Marshal(map[string]string{
		"text": fmt.Sprintf("A spam notification has been detected! ```\rEmail: %s\rFrom: %s\rDescription: %s\r```",
			event.Email,
			event.From,
			event.Description),
	})
	if err != nil {
		fmt.Printf("Error marshalling response: {%s}\r", err.Error())
		return
	}

	resp, err := http.Post(SlackAddress, "application/json", bytes.NewBuffer(eventBody))
	if err != nil {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error posting to slack: %s\rStatus: %d\rBody: %s\r", err.Error(), resp.StatusCode, body)
	}
}
