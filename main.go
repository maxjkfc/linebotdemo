package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/spf13/viper"
)

var (
	bot *linebot.Client
	err error
)

var (
	textSlice = []string{
		"你說什麼聽不懂",
		"不知道",
		"不明白",
		"你好帥",
		"你好酷",
		"不好玩",
		"謝謝再聯絡",
		"要你何用",
	}
)

func main() {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	secret := viper.GetString("secret")
	token := viper.GetString("token")

	// 設定 linebot 配置
	bot, err = linebot.New(secret, token)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("config line bot success")
	}

	httphandler()

}

func httphandler() {

	http.HandleFunc("/callback", callback)

	fmt.Println("start http server in :1314")
	http.ListenAndServe(":1314", nil)

}

func callback(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	if err != nil {
		log.Printf("body parse failed with error : %v", err)
		w.WriteHeader(500)
		return
	}

	for _, event := range events {
		// 判斷是什麼事件
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessate(message.Text))).Do()
				if err != nil {
					fmt.Println(err)
				}
			case *linebot.StickerMessage:
				_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewStickerMessage("11537", "52002763")).Do()
				if err != nil {
					fmt.Println(err)
				}
			}
		default:
			spew.Dump(event)
		}
	}
}

func replyMessate(message string) string {
	switch message {
	case "你好":
		return "Max 歡迎你"

	default:
		return randMessage()
	}
}

func randMessage() string {
	rand.Seed(time.Now().UnixNano())
	return textSlice[rand.Intn(len(textSlice))]
}
