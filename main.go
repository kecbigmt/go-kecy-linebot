package main
import (
	"log"
	"net/http"
	"os"
	"strings"
	"fmt"
	"regexp"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/kecbigmt/go-white-and-black-doors/automata/oldLulu_001"
	"github.com/kecbigmt/go-white-and-black-doors/automata/oldLulu_008"
	"github.com/kecbigmt/go-white-and-black-doors/automata/oldLulu_047"
)

func makeInput(t string) []byte{
	b := make([]byte, len(t))
	for i, l := range t {
		switch l{
		case '0':
			b[i] = uint8(0)
		case '1':
			b[i] = uint8(1)
		default:
			b[i] = uint8(255)
		}
	}
  return b
}

func main() {
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
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
			var text string
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					switch {
					case message.Text == "へい":
						text = "ほー"
					case regexp.MustCompile(`(僕|私|俺|ぼく|わたし|おれ)は(誰|だれ)`).MatchString(message.Text):
						userId := event.Source.UserID
						res, err := bot.GetProfile(userId).Do()
						if err != nil {
							log.Print(err)
	          }
						text = fmt.Sprintf("Display Name:%v\nPicture URL:%v\nStatus Message:%v\nUser ID:%v", res.DisplayName, res.PictureURL, res.StatusMessage, userId)
					case strings.HasPrefix(message.Text, "L1:"):
						t := strings.Replace(message.Text, "L1:", "", 1)
						b := makeInput(t)
						if err := oldLulu_001.Validate(b); err != nil {
							text = fmt.Sprintf("拒否\n%v", err)
						} else {
							text = "受理"
						}
					case strings.HasPrefix(message.Text, "L8:"):
						t := strings.Replace(message.Text, "L8:", "", 1)
						b := makeInput(t)
						if err := oldLulu_008.Validate(b); err != nil {
							text = fmt.Sprintf("拒否\n%v", err)
						} else {
							text = "受理"
						}
					case strings.HasPrefix(message.Text, "L47:"):
						t := strings.Replace(message.Text, "L47:", "", 1)
						b := makeInput(t)
						if err := oldLulu_047.Validate(b); err != nil {
							text = fmt.Sprintf("拒否\n%v", err)
						} else {
							text = "受理"
						}
					default:
						text = message.Text
					}
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
						log.Print(err)
          }
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
