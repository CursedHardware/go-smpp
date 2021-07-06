package templates

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var tpl *template.Template

//go:embed *.gohtml
var fs embed.FS

func init() {
	tpl = template.Must(template.New("").ParseFS(fs, "*.gohtml"))
}

func Render(name string, data interface{}) (rendered string, err error) {
	var buf bytes.Buffer
	err = tpl.ExecuteTemplate(&buf, name+".gohtml", data)
	rendered = buf.String()
	return
}

type ReceiverData struct {
	Sender     string
	Receiver   string
	ReceivedAt time.Time
	Content    string
	Tags       []string
}

type ReplyData struct {
	Sender   string
	Receiver string
	Segments string
}

func Receiver(data *ReceiverData) (string, error)      { return Render("receiver", data) }
func Reply(data *ReplyData) (string, error)            { return Render("reply", data) }
func Submitted() (string, error)                       { return Render("submitted", nil) }
func Unauthorized(user *tgbotapi.User) (string, error) { return Render("unauthorized", user) }

func Must(rendered string, err error) string {
	if err != nil {
		panic(err)
	}
	return rendered
}
