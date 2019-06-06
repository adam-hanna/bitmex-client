// +build unit

package websocket

import (
	"testing"
)

func TestConnectMaster(t *testing.T) {
	conn, err := Connect("www.bitmex.com")
	defer conn.Close()
	if err != nil {
		t.Fatalf("err connecting:\n%v", err)
	}
	if conn == nil {
		t.Error("No connect to ws")
	}
}

func TestConnectDev(t *testing.T) {
	conn, err := Connect("testnet.bitmex.com")
	defer conn.Close()
	if err != nil {
		t.Fatalf("err connecting:\n%v", err)
	}
	if conn == nil {
		t.Error("No connect to testnet ws")
	}
}

func TestWorkerReadMessages(t *testing.T) {
	//chReaderMessage := make(chan interface{})
	//conn, err := Connect("testnet.bitmex.com")
	//defer conn.Close()
	//if err != nil {
	//t.Fatalf("err connecting:\n%v", err)
	//}

	//go ReadFromWSToChannel(conn, chReaderMessage)
	//message := <-chReaderMessage
	//if message == nil {
	//t.Error("Empty message")
	//}
	//t.Log(message)
}

func TestWorkerWriteMessages(t *testing.T) {
	//conn, err := Connect("testnet.bitmex.com")
	//defer conn.Close()
	//if err != nil {
	//t.Fatalf("err connecting:\n%v", err)
	//}

	//// Read
	//chReadFromWS := make(chan interface{}, 10)
	//go ReadFromWSToChannel(conn, chReadFromWS)

	//// Write
	//chWriteToWS := make(chan interface{}, 10)
	//go WriteFromChannelToWS(conn, chWriteToWS)

	//// Send ping
	//chWriteToWS <- []byte(`ping`)

	//// Read first response message
	//message := <-chReadFromWS
	//if !strings.Contains(string(message.([]byte)), "Welcome to the BitMEX") {
	//fmt.Println(string(message.([]byte)))
	//t.Error("No welcome message")
	//}

	//// Read second response message
	//message = <-chReadFromWS
	//if !strings.Contains(string(message.([]byte)), "pong") {
	//fmt.Println(string(message.([]byte)))
	//t.Error("No pong message")
	//}

	//time.Sleep(1 * time.Second)
}
