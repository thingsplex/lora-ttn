package model

import (
	"github.com/TheThingsNetwork/ttn/core/types"
	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/fimptype"
)

type Device struct {
	DeviceId      string // Human readable device ID in TTN
	Address       string // Serial number
	DriverName    string
	lastUplinkMsg types.UplinkMessage
	driver        Driver
}

func (dev *Device) LastUplinkMsg() types.UplinkMessage {
	return dev.lastUplinkMsg
}

func NewDevice(address string, driver Driver) *Device {
	return &Device{Address: address, driver: driver}
}

func (dev *Device) Driver() Driver {
	return dev.driver
}

func (dev *Device) SetDriver(driver Driver) {
	dev.driver = driver
}
func (dev *Device) DecodeAndConvertToFimp( msg *types.UplinkMessage) []fimpgo.FimpMessage {
	dev.lastUplinkMsg = *msg
	return dev.driver.DecodeAndConvertToFimp(dev.Address,msg.PayloadRaw)
}

func (dev *Device) GetInclusionReport() fimptype.ThingInclusionReport  {
	prodInfo := dev.driver.GetProductInfo()
	inclReport := fimptype.ThingInclusionReport{
		IntegrationId:     "",
		Address:           dev.Address,
		Type:              "",
		ProductHash:       prodInfo.ProductHash ,
		Alias:             prodInfo.ProductAlias,
		CommTechnology:    "ttn-lora",
		ProductId:         prodInfo.ProductId,
		ProductName:       prodInfo.ProductName,
		ManufacturerId:    prodInfo.ManufacturerId,
		DeviceId:          dev.Address,
		HwVersion:         "1",
		SwVersion:         prodInfo.SwVersion,
		PowerSource:       prodInfo.PowerSource,
		WakeUpInterval:    "-1",
		Security:          "",
		Tags:              nil,
		Groups:            []string{"ch_0"},
		PropSets:          nil,
		TechSpecificProps: nil,
		Services:          dev.driver.GetSupportedServices(dev.Address),
	}
	return inclReport
}