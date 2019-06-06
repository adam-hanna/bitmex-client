package websocket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// Message ...
type Message struct {
	Op   string        `json:"op,omitempty"`
	Args []interface{} `json:"args,omitempty"`
}

// AddArgument ...
func (m *Message) AddArgument(argument string) {
	m.Args = append(m.Args, argument)
}

// Connect ...
func Connect(host string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: host, Path: "/realtime"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("err dialing %s:\n%v", u.String(), err)
		return nil, err
	}

	return conn, nil
}

// ReadFromWSToChannel ...
func ReadFromWSToChannel(c *websocket.Conn, chRead chan<- interface{}) {
	var (
		message interface{}
		err     error
	)

	for {
		_, message, err = c.ReadMessage()
		if err != nil {
			chRead <- err
		}

		chRead <- message
	}
}

// WriteFromChannelToWS ...
func WriteFromChannelToWS(c *websocket.Conn, chWrite <-chan interface{}) {
	var (
		err     error
		message interface{}
	)

	for {
		switch message = <-chWrite; v := message.(type) {
		case string:
			message, err = json.Marshal(v)
			if err != nil {
				log.Printf("err marshaling %s:\n%v", v, err)
				continue
			}

			if err = c.WriteMessage(websocket.TextMessage, message.([]byte)); err != nil {
				log.Printf("err writing message %s:\n%v", v, err)
				continue
			}

		case error:
			log.Printf("received err on channel:\n%v", err)

		default:
			log.Printf("received unknown message of type %T:\n%v", message, message)
		}
	}
}

// GetAuthMessage ...
func GetAuthMessage(key string, secret string) (*Message, error) {
	nonce := time.Now().Unix() + 412
	req := fmt.Sprintf("GET/realtime%d", nonce)
	sig := hmac.New(sha256.New, []byte(secret))
	if _, err := sig.Write([]byte(req)); err != nil {
		log.Printf("err writing sig:\n%v", err)
		return nil, err
	}
	signature := hex.EncodeToString(sig.Sum(nil))
	var msgKey []interface{}
	msgKey = append(msgKey, key)
	msgKey = append(msgKey, nonce)
	msgKey = append(msgKey, signature)

	return &Message{"authKey", msgKey}, nil
}
