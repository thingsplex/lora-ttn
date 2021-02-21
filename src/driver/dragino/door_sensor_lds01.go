package dragino

import (
	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/fimptype"
	"github.com/thingsplex/lora-ttn/model"
	"github.com/thingsplex/lora-ttn/utils"
)

type DoorSensorLds01 struct {
	temp float32
	batteryVoltage float64
	batteryLevel uint16
	binState bool
}

func NewDoorSensorLds01() model.Driver {
	return &DoorSensorLds01{}
}

// Parsing byte array to internal representation
func (d *DoorSensorLds01) DecodeAndConvertToFimp(address string,payload []byte) []fimpgo.FimpMessage {
	var value = (uint16(payload[0])<<8 | uint16(payload[1]))& 0x3FFF

	d.batteryVoltage = float64(value)/1000 //Battery,units:V
	d.batteryLevel = utils.Convert3vBatteryVoltageToLevel(d.batteryVoltage)
	d.binState = false
	if payload[0]&0x80 > 0 {
		d.binState = true
	} //1:open,0:close

	return d.getFimpEvents(address)
}

func (d *DoorSensorLds01) GetBatteryVoltage() float64 {
	return d.batteryVoltage
}

func (d *DoorSensorLds01) GetBatteryLevel() uint16 {
	return d.batteryLevel
}

func (d *DoorSensorLds01) GetHumidity() float32 {
	return 0
}

func (d *DoorSensorLds01) GetTemperature() float32 {
	return 0
}

func (d *DoorSensorLds01) GetBinaryState() bool {
	return d.binState
}

func (d *DoorSensorLds01) getFimpEvents(address string) []fimpgo.FimpMessage {
	// TODO:somehow convert battery voltage to level
	//msg := fimpgo.NewIntMessage("evt.lvl.report", "battery",100, nil, nil, nil)
	msg := fimpgo.NewBoolMessage("evt.open.report", "sensor_contact", d.binState, nil, nil, nil)
	msg.Topic = "pt:j1/mt:evt/rt:dev/rn:ttn-lora/ad:1/sv:sensor_contact/ad:"+address

	msg2 := fimpgo.NewIntMessage("evt.lvl.report", "battery",int64(d.batteryLevel), nil, nil, nil)
	msg2.Topic = "pt:j1/mt:evt/rt:dev/rn:ttn-lora/ad:1/sv:battery/ad:"+address


	return []fimpgo.FimpMessage{*msg,*msg2}
	//msg = fimpgo.NewFloatMessage("evt.sensor.report", "sensor_humid", val,  map[string]string{"unit":"%"}, nil, nil)
}

func (d *DoorSensorLds01) GetProductInfo() model.ProductInfo {
	return model.ProductInfo{
		ProductHash:    "lora_dragino_lds01",
		ProductAlias:   "LoRaWAN Door Sensor",
		ProductId:      "LDS01",
		ProductName:    "Door sensor",
		ManufacturerId: "Dragino",
		SwVersion:      "",
		PowerSource:    "battery",
	}
}

func (d *DoorSensorLds01) GetSupportedServices(address string) []fimptype.Service {
	batteryInterfaces := []fimptype.Interface{{
		Type:      "in",
		MsgType:   "cmd.lvl.get_report",
		ValueType: "null",
		Version:   "1",
	}, {
		Type:      "out",
		MsgType:   "evt.lvl.report",
		ValueType: "int",
		Version:   "1",
	}, {
		Type:      "out",
		MsgType:   "evt.alarm.report",
		ValueType: "str_map",
		Version:   "1",
	}}

	batteryService := fimptype.Service{
		Name:    "battery",
		Alias:   "battery",
		Address: "/rt:dev/rn:ttn-lora/ad:1/sv:battery/ad:"+address,
		Enabled: true,
		Groups:  []string{"ch_0"},
		Props: map[string]interface{}{},
		Tags:             nil,
		PropSetReference: "",
		Interfaces:       batteryInterfaces,
	}

	contactSensorInterfaces := []fimptype.Interface{{
		Type:      "in",
		MsgType:   "cmd.open.get_report",
		ValueType: "null",
		Version:   "1",
	}, {
		Type:      "out",
		MsgType:   "evt.open.report",
		ValueType: "bool",
		Version:   "1",
	}}

	contactService := fimptype.Service{
		Name:    "sensor_contact",
		Alias:   "Door/window contact",
		Address: "/rt:dev/rn:ttn-lora/ad:1/sv:sensor_contact/ad:"+address,
		Enabled: true,
		Groups:  []string{"ch_0"},
		Props: map[string]interface{}{},
		Tags:             nil,
		PropSetReference: "",
		Interfaces:       contactSensorInterfaces,
	}

	return []fimptype.Service{batteryService,contactService}
}


