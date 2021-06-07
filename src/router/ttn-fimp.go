package router

import (
	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/edgeapp"
	log "github.com/sirupsen/logrus"
	"github.com/thingsplex/lora-ttn/model"
	"github.com/thingsplex/lora-ttn/ttn"
)

type TtnToFimpRouter struct {
	ttnClient    *ttn.LoraTtnClient
	mqt          *fimpgo.MqttTransport
	devDb        *model.DevDB
	appLifecycle *edgeapp.Lifecycle
	appID        string
	appAccessKey string
}

func NewTtnToFimpRouter(mqt *fimpgo.MqttTransport,ttnClient *ttn.LoraTtnClient, devDb *model.DevDB, appLifecycle *edgeapp.Lifecycle) *TtnToFimpRouter {
	ttnr := &TtnToFimpRouter{ttnClient: ttnClient, mqt: mqt, devDb: devDb, appLifecycle: appLifecycle}

	return ttnr
}

func (tr *TtnToFimpRouter) Init() error {
	log.Infof("<ttn-r> Initializing msg router.")
	ttnChannel,err := tr.ttnClient.Subscribe()
	if err != nil {
		return err
	}
	go func() {
		for msg := range ttnChannel {
			log.Debugf("<ttn-r> New msg from dev %s",msg.HardwareSerial)
			dev := tr.devDb.GetDeviceByAddress(msg.HardwareSerial)
			if dev == nil {
				log.Debugf("<ttn-r> Device with id = %s not registered in db ",msg.HardwareSerial)
				continue
			}
			events := dev.DecodeAndConvertToFimp(msg)
			for ei := range events {
				if events[ei].Topic != "" {
					log.Debugf("<ttn-r> Publishing to topic %s",events[ei].Topic)
					tr.mqt.PublishToTopic(events[ei].Topic,&events[ei])
				}else {
					log.Debug("<ttn-r> Msg skipped , empty topic.")
				}

			}
		}
	}()
	log.Infof("<ttn-r> Router init process has completed")
	return nil
}