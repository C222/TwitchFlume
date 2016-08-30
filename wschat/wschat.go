package wschat

import "net/http"
import "fmt"
import "encoding/json"
import "io/ioutil"
import "strings"
import "runtime"

import "github.com/Sirupsen/logrus"
import websocket "golang.org/x/net/websocket"
import logger "../logger"

var log = logger.GetLogger()

type ChatProperties struct {
	_                int      `json:"_id"`
	_                bool     `json:"hide_chat_links"`
	_                int      `json:"chat_delay_duration"`
	_                bool     `json:"devchat"`
	_                string   `json:"game"`
	_                bool     `json:"require_verified_account"`
	_                bool     `json:"subsonly"`
	_                []string `json:"chat_servers"`
	WebSocketServers []string `json:"web_socket_servers"`
	_                float32  `json:"web_socket_pct"`
	_                float32  `json:"darklaunch_pct"`
	_                string   `json:"block_chat_notification_token"`
	_                []string `json:"available_chat_notification_tokens"`
	_                string   `json:"sce_title_preset_text_1"`
	_                string   `json:"sce_title_preset_text_2"`
	_                string   `json:"sce_title_preset_text_3"`
	_                string   `json:"sce_title_preset_text_4"`
	_                string   `json:"sce_title_preset_text_5"`
}

func handleError(e error, die bool) {
	if e != nil {
		if die {
			log.WithFields(logrus.Fields{
				"err": e,
			}).Fatal("WsChat Fatal")
		} else {
			log.WithFields(logrus.Fields{
				"err": e,
			}).Error("WsChat Error")
		}
	}
}

func getWsUri(channel string) (string, error) {
	query_uri := fmt.Sprintf("https://api.twitch.tv/api/channels/%s/chat_properties", channel)

	log.WithFields(logrus.Fields{
		"query_uri": query_uri,
	}).Info("Querying")

	resp, e := http.Get(query_uri)
	if e != nil {
		log.WithFields(logrus.Fields{
			"err": e,
		}).Debug("Error quering channel properties")
		return "", fmt.Errorf("Error quering channel properties")
	}
	defer resp.Body.Close()

	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.WithFields(logrus.Fields{
			"err": e,
		}).Debug("Error reading channel properties")
		return "", fmt.Errorf("Error reading channel properties")
	}

	var properties ChatProperties
	e = json.Unmarshal(body, &properties)
	if e != nil {
		log.WithFields(logrus.Fields{
			"err": e,
		}).Debug("Error parsing channel properties")
		return "", fmt.Errorf("Error parsing channel properties")
	}

	return fmt.Sprintf("ws://%s/", properties.WebSocketServers[0]), nil
}

type WsIrc struct {
	Channel   string
	WS        *websocket.Conn
	OnMessage func(string) bool
}

func (self *WsIrc) Start() {
	log.WithFields(logrus.Fields{
		"channel": self.Channel,
	}).Info("Starting new connection.")

	uri, e := getWsUri(self.Channel)
	handleError(e, true)

	log.WithFields(logrus.Fields{
		"uri": uri,
	}).Info("Got connection URI.")

	origin := "http://localhost/"
	self.WS, e = websocket.Dial(uri, "", origin)
	handleError(e, true)

	_, err := self.WS.Write([]byte("CAP REQ :twitch.tv/tags\n"))
	handleError(err, false)
	_, err = self.WS.Write([]byte("PASS blah\n"))
	handleError(err, false)
	_, err = self.WS.Write([]byte("NICK justinfan47865\n"))
	handleError(err, false)

	_, err = self.WS.Write([]byte(fmt.Sprintf("JOIN #%s\n", strings.ToLower(self.Channel))))
	handleError(err, true)

	go self.messageListener()
}

func (self *WsIrc) messageListener() {
	var msg = make([]byte, 65535)
	var n int
	var e error
	var i int = 0
	var chat_line string
	for {
		// runtime.KeepAlive(self)
		n, e = self.WS.Read(msg)
		handleError(e, false)
		chat_line = string(msg[:n])
		if i%10 == 0 {
			log.WithFields(logrus.Fields{
				"i":       i,
				"channel": self.Channel,
			}).Info("Message Count")
		}
		i += 1
		if strings.HasPrefix(chat_line, "PING") {
			n, e = self.WS.Write([]byte("PONG\n"))
			log.WithFields(logrus.Fields{
				"channel": self.Channel,
			}).Info("PING/PONG")
			handleError(e, false)
		} else {
			go self.OnMessage(chat_line)
		}
		runtime.Gosched()
	}
}
