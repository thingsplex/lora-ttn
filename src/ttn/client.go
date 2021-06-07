package ttn

import (
	"cloud.google.com/go/pubsub"
	"fmt"
	ttnsdk "github.com/TheThingsNetwork/go-app-sdk"
	"github.com/TheThingsNetwork/ttn/core/types"
	log "github.com/sirupsen/logrus"
	//ttnlog "github.com/TheThingsNetwork/go-utils/log"
	//"github.com/TheThingsNetwork/go-utils/log/apex"
	//"github.com/TheThingsNetwork/go-utils/random"
	//"github.com/TheThingsNetwork/ttn/core/types"
)

const (
	sdkClientName = "lora-ttn-edge-app"
)

type LoraTtnClient struct {
	appID string
	appAccessKey string
	client ttnsdk.Client
	pubsub pubsub.Client
	allDevicesPubSub ttnsdk.DeviceSub
}

func NewLoraTtnClient() *LoraTtnClient {
	return &LoraTtnClient{}
}

func (lc *LoraTtnClient) Init(appID string, appAccessKey string) {

	//ttnlog.Set
	// Create a new SDK configuration for the public community network
	config := ttnsdk.NewCommunityConfig(sdkClientName)
	config.ClientVersion = "2.0.5" // The version of the application
	// Create a new SDK client for the application
	lc.appAccessKey = appAccessKey
	lc.appID = appID
	lc.client = config.NewClient(appID, appAccessKey)
}

func (lc *LoraTtnClient) Subscribe() (<-chan *types.UplinkMessage,error) {
	// Start Publish/Subscribe client (MQTT)
	pubsub, err := lc.client.PubSub()
	if err != nil {
		log.WithError(err).Error("my-amazing-app: could not get application pub/sub")
		return nil, err
	}

	// Get a publish/subscribe client for all devices
	lc.allDevicesPubSub = pubsub.AllDevices()

	// Subscribe to events
	//events, err := lc.allDevicesPubSub.SubscribeEvents()
	//if err != nil {
	//	log.WithError(err).Error("my-amazing-app: could not subscribe to events")
	//}
	//log.Debug("After this point, the program won't show anything until we receive an application event.")

	// Subscribe to uplink messages
	return lc.allDevicesPubSub.SubscribeUplink()
	//uplink, err := lc.allDevicesPubSub.SubscribeUplink()
	//if err != nil {
	//	log.WithError(err).Error("my-amazing-app: could not subscribe to uplink messages")
	//}
	//log.Debug("After this point, the program won't show anything until we receive an uplink message from device my-new-device.")
	//for message := range uplink {
	//	hexPayload := hex.EncodeToString(message.PayloadRaw)
	//	payload := message.PayloadRaw
	//	value:= (uint16(payload[0])<<8 | uint16(payload[1]))& 0x3FFF
	//	bat := float32(value)/1000 //Battery,units:V
	//	isDoorOpen := false
	//	if payload[0]&0x80 > 0 {
	//		isDoorOpen = true
	//	} //1:open,0:close
	//	log.WithField("data", hexPayload).Infof("my-amazing-app: received uplink from device : %+v",message)
	//	log.Infof("Parsed batt = %f , is_door_open = %t",bat,isDoorOpen)
	//	//break // normally you wouldn't do this
	//}

}

//func (lc *LoraTtnClient) ParseDeviceMsg(serial string , payload []byte) {
//
//	var value = (uint16(payload[0])<<8 | uint16(payload[1]))& 0x3FFF
//
//	bat := float32(value)/1000 //Battery,units:V
//	isDoorOpen := false
//	if payload[0]&0x80 > 0 {
//		isDoorOpen = true
//	} //1:open,0:close
//	log.Infof("Parsed batt = %f , is_door_open = %t",bat,isDoorOpen)
//
//}


// AddDevice
// deviceId - TTN internal device ID
// deviceEUI - device global EUI address
// applicationEUI - LoraWan application id
// appKey - LoraWan application key
func (lc *LoraTtnClient) AddDevice(deviceId,description,deviceEUI,applicationEUI,appKey string ) (*ttnsdk.Device, error) {
	devMan , err := lc.client.ManageDevices()

	if err != nil {
		return nil, fmt.Errorf("device manager error .err = %s",err.Error())
	}

	if exDev , _ := lc.GetDeviceByDevEui(deviceEUI);exDev != nil {
		log.Debugf("<ttn-cl> Device already registered , returning existing device")
		return exDev.AsDevice(),nil
	}


	devEUI , err := types.ParseDevEUI(deviceEUI)
	if err != nil {
		return nil, fmt.Errorf("incorrect deviceEUI err = %s",err.Error())
	}

	appEUI , err := types.ParseAppEUI(applicationEUI)
	if err != nil {
		return nil,fmt.Errorf("incorrect appEUI err = %s",err.Error())
	}

	appKeyT,err := types.ParseAppKey(appKey)
	if err != nil {
		return nil,fmt.Errorf("incorrect appKey.err = %s",err.Error())
	}

	dev := ttnsdk.Device{}
	dev.AppID = lc.appID
	dev.DevID = deviceId
	dev.Description = description
	dev.DevEUI = devEUI
	dev.AppEUI = appEUI
	dev.AppKey = &appKeyT


	err = devMan.Set(&dev)
	if err != nil {
		return nil, fmt.Errorf("can't add new device to ttn .err = %s",err.Error())
	}

	return devMan.Get(deviceId)
}

func (lc *LoraTtnClient) Cleanup() {
	lc.allDevicesPubSub.Close()
	lc.pubsub.Close()
	lc.client.Close()
}

func (lc *LoraTtnClient) GetDeviceByDevEui(deviceEUI string)(*ttnsdk.SparseDevice, error) {
	devEUI , err := types.ParseDevEUI(deviceEUI)
	if err != nil {
		return nil, fmt.Errorf("incorrect deviceEUI err = %s",err.Error())
	}
	devices, err := lc.client.ManageDevices()
	if err != nil {
		log.WithError(err).Error("could not get device manager")
	}
	// List the first 10 devices
	deviceList, err := devices.List(100, 0)
	if err != nil {
		log.WithError(err).Error("<ttn-cl> could not get devices")
	}
	log.Info("<ttn-cl> found devices")
	for _, device := range deviceList {
		if device.DevEUI == devEUI {
			return device , nil
		}
	}
	return nil,nil
}

func (lc *LoraTtnClient) ListDevices() {
	// Manage devices for the application.
	devices, err := lc.client.ManageDevices()
	if err != nil {
		log.WithError(err).Error("<ttn-cl> could not get device manager")
	}
	// List the first 10 devices
	deviceList, err := devices.List(100, 0)
	if err != nil {
		log.WithError(err).Error("<ttn-cl> could not get devices")
	}
	log.Info("<ttn-cl> found devices")
	for _, device := range deviceList {
		log.Infof("- %s", device.DevID)
	}
}