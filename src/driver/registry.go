package driver

import (
	"github.com/thingsplex/lora-ttn/driver/dragino"
	"github.com/thingsplex/lora-ttn/model"
)

type Constructor func() model.Driver

var Registry = map[string]Constructor{
	"dragino_lds01":      dragino.NewDoorSensorLds01,
}