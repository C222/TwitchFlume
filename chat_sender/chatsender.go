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

type Tags struct {
	Badges      string
	Color       string
	DisplayName string
	Emotes      string
	Id          string
	Mod         bool
	RoomId      int
	Subscriber  bool
	Turbo       bool
	UserId      int
	UserType    string
}

type ChatLine struct {
	Tags
	Channel  string
	Message  string
	Username string
}

func handleError(e error, die bool) bool {
	if e != nil {
		if die {
			log.WithFields(logrus.Fields{
				"err": e,
			}).Fatal("ChatSender Fatal")
		} else {
			log.WithFields(logrus.Fields{
				"err": e,
			}).Error("ChatSender Error")
		}
		return true
	}
	return false
}

func ParseLine(msg string) (*ChatLine, error) {
	var chat_line ChatLine
	split_line := strings.SplitN(msg, " :", 3)
	if len(split_line) != 3 || !strings.HasSuffix(msg, "\r\n") {
		return nil, fmt.Errorf("Malformed Chat Line")
	}
	tags := strings.TrimLeft(split_line[0], "@")
	for _, tag := range strings.Split(tags, ";") {
		tag_split := strings.SplitN(tag, "=", 2)
		switch tag_split[0] {
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

	chat_line.Username = strings.SplitN(split_line[1], "!", 2)[0]
	chat_line.Channel = strings.SplitN(split_line[1], "#", 2)[1]

	chat_line.Message = strings.TrimRight(split_line[2], "\n\r")

	return &chat_line, nil
}

func (self *ChatSender) SendLine(msg string) bool {
	if strings.Contains(msg, " PRIVMSG ") {
		chat_line, e := ParseLine(msg)
		if e == nil {
			log.WithFields(logrus.Fields{
				"chat": chat_line,
			}).Debug("Message Parsed")

			uri := fmt.Sprintf("http://logs-01.loggly.com/inputs/%s/tag/%s/", self.Key, chat_line.Channel)

			json_s, e := json.Marshal(chat_line)
			handleError(e, false)

			resp, e := http.Post(uri, "application/json", strings.NewReader(string(json_s)))
			if handleError(e, false) == false {
				defer resp.Body.Close()
				log.WithFields(logrus.Fields{
					"resp": resp,
					}).Debug("Message Sent")

				return true
			} else {
				return false
			}
		}
	}
	return false
}
