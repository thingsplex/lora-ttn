package model

import (
	"fmt"
	"github.com/thingsplex/lora-ttn/driver"
)

type DevDB struct {
	devices []*Device
}

func NewDevDB() *DevDB {
	return &DevDB{devices: []*Device{}}
}

func (db *DevDB) AddDeviceByModel(serial,model string ) error {
	drv,ok := driver.Registry[model]
	if !ok {
		return fmt.Errorf("unknown device model")
	}
	db.devices = append(db.devices,NewDevice(serial,drv()))
	return nil
}

func (db *DevDB) AddDeviceWithDriver(dev *Device)  {
	db.devices = append(db.devices,dev)
}

func (db *DevDB) GetDeviceByAddress(addr string) *Device  {
	for i := range db.devices {
		if db.devices[i].Address == addr {
			return db.devices[i]
		}
	}
	return nil
}
