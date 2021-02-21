package router

import (
	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/fimptype"
	"github.com/thingsplex/lora-ttn/ttn"
)

type ListReportRecord struct {
	Address     string `json:"address"`
	Alias       string `json:"alias"`
	PowerSource string `json:"power_source"`
	Hash        string `json:"hash"`
}

type OpResponse struct {

}

type NetworkService struct {
	mqt          *fimpgo.MqttTransport
	ttnClient    *ttn.LoraTtnClient
}

func NewNetworkService(mqt *fimpgo.MqttTransport, ttnClient    *ttn.LoraTtnClient) *NetworkService {
	return &NetworkService{mqt: mqt, ttnClient :ttnClient}
}

func (ns *NetworkService) OpenNetwork(open bool) error {

	return nil
}

func (ns *NetworkService) DeleteThing(deviceId string) error {

	ns.SendExclusionReport(deviceId)
	return nil
}
func (ns *NetworkService) SendExclusionReport(thingId string){
	val := fimptype.ThingExclusionReport{Address: thingId}
	msg := fimpgo.NewMessage("evt.thing.exclusion_report", "conbee", fimpgo.VTypeObject, val, nil, nil, nil)
	adr := fimpgo.Address{MsgType: fimpgo.MsgTypeEvt, ResourceType: fimpgo.ResourceTypeAdapter, ResourceName: "conbee", ResourceAddress: "1"}
	ns.mqt.Publish(&adr, msg)
}

func (ns *NetworkService) SendAllExclusionReports() error {

	return nil
}

func (ns *NetworkService) SendAllInclusionReports() error {

	return nil
}


func (ns *NetworkService) SendInclusionReport(deviceType, deviceId string) error {
	var report interface{}
	msg := fimpgo.NewMessage("evt.thing.inclusion_report", "conbee", fimpgo.VTypeObject, report, nil, nil, nil)
	adr := fimpgo.Address{MsgType: fimpgo.MsgTypeEvt, ResourceType: fimpgo.ResourceTypeAdapter, ResourceName: "conbee", ResourceAddress: "1"}
	ns.mqt.Publish(&adr, msg)
	return nil

}

func (ns *NetworkService) SendListOfDevices() error {
	var report interface{}
	msg := fimpgo.NewMessage("evt.network.all_nodes_report", "conbee", fimpgo.VTypeObject, report, nil, nil, nil)
	adr := fimpgo.Address{MsgType: fimpgo.MsgTypeEvt, ResourceType: fimpgo.ResourceTypeAdapter, ResourceName: "conbee", ResourceAddress: "1"}
	ns.mqt.Publish(&adr, msg)

	return nil
}
