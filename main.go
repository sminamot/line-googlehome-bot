package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/evalphobia/google-home-client-go/googlehome"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/yosssi/go-voicetext"
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
					profile, err := bot.GetProfile(event.Source.UserID).Do()
					var userName string
					if err == nil {
						userName = profile.DisplayName
					}
					if err := saveWav(message.Text, userName); err != nil {
						log.Print(err)
						continue
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

	http.Handle("/voice/", http.StripPrefix("/voice/", http.FileServer(http.Dir("./static"))))

	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func speak(text, user string) error {
	cli, err := googlehome.NewClientWithConfig(googlehome.Config{
		Hostname: os.Getenv("GOOGLE_HOME_IP"),
		Lang:     "ja",
	})

	if err != nil {
		return err
	}

	if user != "" {
		text = fmt.Sprintf("%sからのメッセージが届きました。。。%s", user, text)
	}

	cli.Play(fmt.Sprintf("http://127.0.0.1:%s/voice/line.wav", os.Getenv("PORT")))

	return nil
	/*
		    ctx := context.Background()
		    devices := homecast.LookupAndConnect(ctx)

		    if user != "" {
		        text = fmt.Sprintf("%sからのメッセージが届きました。。。%s", user, text)
		    }

		    for _, device := range devices {
		        //return device.Speak(ctx, text, "ja")
		        u, _ := url.Parse(fmt.Sprintf("http://%s/voice/line.wav", os.Getenv("MEDIA_DOMAIN")))
		        return device.Play(ctx, u)
		    }

			return errors.New("not found devices")
	*/
}

func saveWav(text, user string) error {
	c := voicetext.NewClient(os.Getenv("VOICETEXT_API_KEY"), nil)
	result, err := c.TTS(fmt.Sprintf("%sからのメッセージが届きました。%s", user, text), &voicetext.TTSOptions{
		Speaker: voicetext.SpeakerHikari,
		Volume:  200,
	})
	if err != nil {
		return err
	}

	if result.ErrMsg != nil {
		return fmt.Errorf("%s", result.ErrMsg)
	}

	f, err := os.Create("./static/line.wav")
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(result.Sound)
	return err
}
