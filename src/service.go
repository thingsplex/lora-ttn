package main

import (
	"flag"
	"fmt"
	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/discovery"
	"github.com/futurehomeno/fimpgo/edgeapp"
	log "github.com/sirupsen/logrus"
	"github.com/thingsplex/lora-ttn/model"
	"github.com/thingsplex/lora-ttn/router"
	"github.com/thingsplex/lora-ttn/ttn"
	"time"
)

func main() {
	var workDir string
	flag.StringVar(&workDir, "c", "", "Work dir")
	flag.Parse()
	if workDir == "" {
		workDir = "./"
	} else {
		fmt.Println("Work dir ", workDir)
	}
	appLifecycle := edgeapp.NewAppLifecycle()
	configs := model.NewConfigs(workDir)
	err := configs.LoadFromFile()
	if err != nil {
		fmt.Print(err)
		panic("Can't load config file.")
	}

	edgeapp.SetupLog(configs.LogFile, configs.LogLevel, configs.LogFormat)
	log.Info("--------------Starting lora-ttn----------------")
	log.Info("Work directory : ", configs.WorkDir)

	appLifecycle.SetAppState(edgeapp.AppStateNotConfigured, nil)

	ttnClient := ttn.NewLoraTtnClient()

	devDb := model.NewDevDB(configs.GetDataDir())
	devDb.LoadFromFile()
	//devDb.AddDeviceByModel("A840418C71825AE0","dragino_lds01")

	mqtt := fimpgo.NewMqttTransport(configs.MqttServerURI, configs.MqttClientIdPrefix, configs.MqttUsername, configs.MqttPassword, true, 1, 1)
	err = mqtt.Start()
	responder := discovery.NewServiceDiscoveryResponder(mqtt)
	responder.RegisterResource(model.GetDiscoveryResource())
	responder.Start()

	fimpRouter := router.NewFromFimpRouter(mqtt,ttnClient,devDb, appLifecycle, configs)
	fimpRouter.Start()
	//------------------ Remote API check -- !!!!IMPORTANT!!!!-------------
	// The app MUST perform remote API availability check.
	// During gateway boot process the app might be started before network is initialized or another local app booted.
	// Remove that codee if the app is not dependent from local network internet availability.
	//------------------ Sample code --------------------------------------
	sys := edgeapp.NewSystemCheck()
	sys.WaitForInternet(time.Minute * 60)
	//---------------------------------------------------------------------

	ttnRouter := router.NewTtnToFimpRouter(mqtt,ttnClient,devDb,appLifecycle)

	appLifecycle.SetAppState(edgeapp.AppStateRunning, nil)

	if configs.IsConfigured() {
		appLifecycle.SetConfigState(edgeapp.ConfigStateConfigured)
	}

	//------------------ Sample code --------------------------------------

	for {
		appLifecycle.WaitForState("main", edgeapp.SystemEventTypeConfigState, edgeapp.ConfigStateConfigured)

		ttnClient.Init(configs.AppID,configs.AppAccessKey)
		ttnRouter.Init()

		appLifecycle.WaitForState("main", edgeapp.SystemEventTypeConfigState, edgeapp.ConfigStateNotConfigured)
	}

	mqtt.Stop()
	time.Sleep(5 * time.Second)
}
