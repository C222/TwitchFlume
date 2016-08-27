package chatsender

import "net/http"
import "github.com/Sirupsen/logrus"
import logger "../logger"
import "strings"
import "fmt"
import "encoding/json"

var log = logger.GetLogger()

type ChatSender struct {
  Key string
}

type ChatLine struct {
  Tags string
  Channel string
  Message string
}

func ParseLine(msg string) *ChatLine {
  var chat_line ChatLine
  split := strings.Split(msg, " PRIVMSG ")
  latter := strings.SplitN(split[1], " :", 2)

  chat_line.Channel = strings.TrimLeft(latter[0], "#")
  chat_line.Message = latter[1]
  chat_line.Tags = strings.SplitN(split[0], " :", 2)[0]
  return &chat_line
}

func (self *ChatSender) SendLine(msg string) bool {
  if strings.Contains(msg, " PRIVMSG "){
    chat_line := ParseLine(msg)
    log.WithFields(logrus.Fields{
      "line": chat_line,
    }).Info("Message Recieved")

    uri := fmt.Sprintf("http://logs-01.loggly.com/inputs/%s/tag/%s/", self.Key, chat_line.Channel)

    json_s, _ := json.Marshal(chat_line)

    _, _ = http.Post(uri, "application/json", strings.NewReader(string(json_s)))

    return true
  }
  return false
}
