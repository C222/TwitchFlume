package chatsender

import "net/http"
import "github.com/Sirupsen/logrus"
import logger "../logger"
import "strings"
import "fmt"
import "encoding/json"
import "strconv"

var log = logger.GetLogger()

type ChatSender struct {
  Key string
}

type Tags struct{
  Badges string
  Color string
  DisplayName string
  Emotes string
  Id string
  Mod bool
  RoomId int
  Subscriber bool
  Turbo bool
  UserId int
  UserType string
}

type ChatLine struct {
  Tags
  Channel string
  Message string
}

func ParseLine(msg string) *ChatLine {
  var chat_line ChatLine
  split := strings.Split(msg, " PRIVMSG ")
  latter := strings.SplitN(split[1], " :", 2)
  tags := strings.TrimLeft(strings.SplitN(split[0], " :", 2)[0], "@")

  for _, tag := range strings.Split(tags, ";"){
    tag_split := strings.SplitN(tag, "=", 2)
    switch tag_split[0]{
      case "badges":
        chat_line.Badges = tag_split[1]
      case "color":
        chat_line.Color = tag_split[1]
      case "display-name":
        chat_line.DisplayName = tag_split[1]
      case "emotes":
        chat_line.Emotes = tag_split[1]
      case "id":
        chat_line.Id = tag_split[1]
      case "mod":
        chat_line.Mod = (tag_split[1] == "1")
      case "room-id":
        chat_line.RoomId, _ = strconv.Atoi(tag_split[1])
      case "subscriber":
        chat_line.Subscriber = (tag_split[1] == "1")
      case "turbo":
        chat_line.Turbo = (tag_split[1] == "1")
      case "user-id":
        chat_line.UserId, _ = strconv.Atoi(tag_split[1])
      case "user-type":
        chat_line.UserType = tag_split[1]
    }
  }

  chat_line.Channel = strings.TrimLeft(latter[0], "#")
  chat_line.Message = strings.TrimRight(latter[1], "\n\r")
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
