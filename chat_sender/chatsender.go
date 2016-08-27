package chatsender

// import "net/http"
import "github.com/Sirupsen/logrus"
import logger "../logger"

var log = logger.GetLogger()

type ChatSender struct {
  Key string
}

func (self *ChatSender) SendLine(msg string) bool {
  log.WithFields(logrus.Fields{
    "message": msg,
  }).Info("Message Recieved")
  return true
}
