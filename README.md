[![GoDoc](https://godoc.org/github.com/adam-hanna/bitmex-client?status.svg)](https://godoc.org/github.com/adam-hanna/bitmex-client)

# Bitmex API
Packages provides a golang client for the bitmex rest and websocket API's.


## Usage

### Config
```json
{
  "IsDev": false,
  "Master": {
    "Host": "www.bitmex.com",
    "Key": "123",
    "Secret": "abc",
    "Timeout": 20
  },
  "Dev": {
    "Host": "testnet.bitmex.com",
    "Key": "123",
    "Secret": "abc",
    "Timeout": 20
  }
}
```

####  REST
```
// Load config
cfg, err := config.LoadConfig("config.json")
if err != nil {
  log.Fatalf(err)
}
ctx := rest.MakeContext(cfg.Key, cfg.Secret, cfg.Host)

// Get wallet
wallet, response, err := rest.GetWallet(ctx)
if err != nil {
  log.Fatal(err)
}
fmt.Printf("Status: %v, wallet amount: %v\n", response.StatusCode, wallet.Amount)

// Place order
params := map[string]interface{}{
    "side":     "Buy",
    "symbol":   "XBTUSD",
    "ordType":  "Limit",
    "orderQty": 1,
    "price":    9000,
    "clOrdID":  "MyUniqID_123",
    "execInst": "ParticipateDoNotInitiate",
}
order, response, err := rest.NewOrder(ctx, params)
if err != nil {
  log.Fatal(err)
}
fmt.Printf("Order: %+v, Response: %+v\n", order, response)
```

#### Websocket
```
// Load config
cfg, err := config.LoadConfig("config.json")
if err != nil {
  log.Fatal(err)
}

// Connect to WS
conn := websocket.Connect(cfg.Host)
defer conn.Close()

// Listen read WS
chReadFromWS := make(chan interface{}, 100)
go websocket.ReadFromWSToChannel(conn, chReadFromWS)

// Listen write WS
chWriteToWS := make(chan interface{}, 100)
go websocket.WriteFromChannelToWS(conn, chWriteToWS)

// Authorize
authMsg, err := websocket.GetAuthMessage(cfg.Key, cfg.Secret)
if err != nil {
  log.Fatal(err)
}
chWriteToWS <- authMsg

// Listen

go func() {
		var (
			err     error
			message interface{}
			res     *bitmex.Response
		)

    for {
			switch message = <-chReadFromWS; v := message.(type) {
			case []byte:
        res, err = bitmex.DecodeMessage(v)
        if err != nil {
          log.Println(err)
          continue
        }

        // Business logic
        switch res.Table {
        case "orderBookL2":
            if res.Action == "partial" {
                // Update table
            } else {
                // Update row
            }
        case "order":
            if res.Action == "partial" {
                // Update table
            } else {
                // Update row
            }
        case "position":
            if res.Action == "partial" {
                // Update table
            } else {
                // Update row
            }
        }

			case error:
				log.Printf("err reading from ws:\n%v", v)

			default:
				log.Printf("unknown message type %T:\n%v", message, message)
      }
    }
}()

```

## Example
```
package main

import (
	"log"
	"strings"

	"github.com/adam-hanna/bitmex-client/bitmex"
	"github.com/adam-hanna/bitmex-client/config"
	"github.com/adam-hanna/bitmex-client/rest"
	"github.com/adam-hanna/bitmex-client/websocket"
)

// Usage example
func main() {
	// Load config
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("could not load config.json")
	}

	ctx := rest.MakeContext(cfg.Key, cfg.Secret, cfg.Host, cfg.Timeout)

	// Get wallet
	w, response, err := rest.GetWallet(ctx)
	if err != nil {
		log.Fatalf("err getting wallet:\n%v", err)
	}

	log.Printf("Status: %v, wallet amount: %v\n", response.StatusCode, w.Amount)

	// Connect to WS
	conn, err := websocket.Connect(cfg.Host)
	if err != nil {
		log.Fatalf("could not connect ws to %s:\n%v", cfg.Host, err)
	}
	defer func() {
		if err = conn.Close(); err != nil {
			log.Fatalf("err closing ws conn:\n%v", err)
		}
	}()

	// Listen read WS
	chReadFromWS := make(chan interface{}, 100)
	go websocket.ReadFromWSToChannel(conn, chReadFromWS)

	// Listen write WS
	chWriteToWS := make(chan interface{}, 100)
	go websocket.WriteFromChannelToWS(conn, chWriteToWS)

	// Authorize
	authMsg, err := websocket.GetAuthMessage(cfg.Key, cfg.Secret)
	if err != nil {
		log.Fatalf("err getting auth message from ws:\n%v", err)
	}
	chWriteToWS <- authMsg

	// Read first response message
	message := <-chReadFromWS
	if !strings.Contains(string(message.([]byte)), "Welcome to the BitMEX") {
		log.Println(string(message.([]byte)))
		log.Fatal("No welcome message")
	}

	// Read auth response success
	switch message := <-chReadFromWS; v := message.(type) {
	case []byte:
		res, err := bitmex.DecodeMessage(v)
		if err != nil {
			log.Fatalf("err decoding message:\n%v", err)
		}

		if !res.Success || res.Request.(map[string]interface{})["op"] != "authKey" {
			log.Fatal("No auth response success")
		}

	case error:
		log.Printf("err reading from ws:\n%v", v)

	default:
		log.Printf("unknown message type %T:\n%v", message, message)
	}

	// Listen websocket before subscribe
	go func() {
		var (
			err     error
			message interface{}
			res     *bitmex.Response
		)

		for {
			switch message = <-chReadFromWS; v := message.(type) {
			case []byte:
				res, err = bitmex.DecodeMessage(v)
				if err != nil {
					log.Printf("err decoding message:\n%v", err)
					continue
				}

				// Your logic here
				log.Printf("%+v\n", res)

			case error:
				log.Printf("err reading from ws:\n%v", v)

			default:
				log.Printf("unknown message type %T:\n%v", message, message)
			}
		}
	}()

	// Subscribe
	messageWS := websocket.Message{Op: "subscribe"}
	messageWS.AddArgument("orderBookL2:XBTUSD")
	messageWS.AddArgument("order")
	messageWS.AddArgument("position")
	chWriteToWS <- messageWS

	// Loop forever
	select {}
}
```
