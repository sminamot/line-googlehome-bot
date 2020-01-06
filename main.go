package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ikasamah/homecast"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					fmt.Println(event.Source.UserID)
					profile, err := bot.GetProfile(event.Source.UserID).Do()
					var userName string
					if err == nil {
						userName = profile.DisplayName
					}
					if err := speak(message.Text, userName); err != nil {
						log.Print(err)
					}
					/*
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
							log.Print(err)
						}
					*/
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func speak(text, user string) error {
	ctx := context.Background()
	devices := homecast.LookupAndConnect(ctx)

	if user != "" {
		text = fmt.Sprintf("%sからのメッセージが届きました。。。%s", user, text)
	}

	for _, device := range devices {
		return device.Speak(ctx, text, "ja")
	}

	return errors.New("not found devices")
}
