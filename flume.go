package main

import "github.com/Sirupsen/logrus"
import logger "./logger"
import wschat "./wschat"
import "encoding/json"
import "io/ioutil"
import "./chat_sender"

var log = logger.GetLogger()

type Config struct {
	Channels   []string  `json:"channels"`
	LogglyKey string  `json:"loggly-key"`
}

func main() {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		log.WithFields(logrus.Fields{
			"err": e,
		}).Fatal("Error opening config")
	}

  var config Config
  e = json.Unmarshal(file, &config)
	if e != nil {
		log.WithFields(logrus.Fields{
			"err": e,
		}).Fatal("Error parsing config")
	}

  log.WithFields(logrus.Fields{
    "channels": config.Channels,
    "loggly_key": config.LogglyKey,
  }).Info("Config loaded.")

  cs := chatsender.ChatSender{config.LogglyKey}

  for _, channel := range config.Channels{
    ws := wschat.WsIrc{channel, nil, cs.SendLine}
    ws.Start()
  }
  for {}
}
