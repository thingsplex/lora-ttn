package ttn

import (
	log "github.com/sirupsen/logrus"
	"testing"
)


func TestLoraTtnClient_Init(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	client := NewLoraTtnClient()
	client.Init("test-lab","")
	client.ListDevices()
	//client.Sub()
}

func TestLoraTtnClient_Parse(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	//client := NewLoraTtnClient("test-lab","")
	//payload,_ := hex.DecodeString("8bca0100006b000001")
	//client.ParseDeviceMsg("",payload)
}
