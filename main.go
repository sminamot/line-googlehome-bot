package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/evalphobia/google-home-client-go/googlehome"
	"github.com/line/line-bot-sdk-go/linebot"
)

const voinceTextWebAPIURL = "https://api.voicetext.jp/v1/tts"

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
					if err := speak(message.Text, userName); err != nil {
						log.Print(err)
					}
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

	url, err := getVoiceURL(text)
	if err != nil {
		return err
	}

	return cli.Play(url)
}

func getVoiceURL(text string) (string, error) {
	v := url.Values{}
	v.Add("text", text)
	v.Add("speaker", "hikari")
	v.Add("format", "mp3")
	v.Add("volume", "200")

	res, err := requestVoiceTextWebAPI(os.Getenv("VOICETEXT_API_KEY"), v)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", errors.New("failed to request VoiceText Web API")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")},
	)
	if err != nil {
		return "", err
	}

	bucket := aws.String(os.Getenv("AWS_S3_BUCKET"))
	filename := aws.String(fmt.Sprint(time.Now().UnixNano()) + ".mp3")
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: bucket,
		Key:    filename,
		Body:   res.Body,
	})
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: bucket,
		Key:    filename,
	})
	url, err := req.Presign(30 * time.Second) // 有効期限を指定して署名付きURLを取得
	if err != nil {
		fmt.Println("failed to S3 sign request", err)
	}

	return url, nil
}

func requestVoiceTextWebAPI(key string, values url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, voinceTextWebAPIURL, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = values.Encode()
	req.SetBasicAuth(key, "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return http.DefaultClient.Do(req)
}
