package model

import (
	"encoding/json"
	"fmt"
	"github.com/thingsplex/lora-ttn/driver"
	"io/ioutil"
	"path/filepath"
)

type DevDB struct {
	devices []Device
	dataDirPath string
	dbFullPath string
}

func (db *DevDB) Devices() []Device {
	return db.devices
}

func NewDevDB(dataDirPath string) *DevDB {
	dbFullPath := filepath.Join(dataDirPath, "dev_db.json")
	return &DevDB{dataDirPath: dataDirPath,dbFullPath: dbFullPath}
}

func (db *DevDB) AddDeviceByModel(serial,model string ) error {
	drv,ok := driver.Registry[model]
	if !ok {
		return fmt.Errorf("unknown device model")
	}
	db.devices = append(db.devices,*NewDevice(serial,drv()))
	return nil
}

func (db *DevDB) AddDeviceWithDriver(dev Device)  {
	db.devices = append(db.devices,dev)
}

func (db *DevDB) GetDeviceByAddress(addr string) *Device  {
	for i := range db.devices {
		if db.devices[i].Address == addr {
			return &db.devices[i]
		}
	}
	return nil
}

func (db *DevDB) LoadFromFile() error {
	configFileBody, err := ioutil.ReadFile(db.dbFullPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(configFileBody, db.devices)
	if err != nil {
		return err
	}
	return nil
}

func (db *DevDB) SaveToFile() error {
	bpayload, err := json.Marshal(db.devices)
	err = ioutil.WriteFile(db.dbFullPath, bpayload, 0664)
	if err != nil {
		return err
	}
	return err
}
