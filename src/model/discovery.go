package model

import (
	"github.com/futurehomeno/fimpgo/discovery"
)

func GetDiscoveryResource() discovery.Resource {
	return discovery.Resource{
		ResourceName:           ServiceName,
		ResourceType:           discovery.ResourceTypeAd,
		Author:                 "aleksandrs.livincovs@gmail.com",
		IsInstanceConfigurable: false,
		InstanceId:             "1",
		Version:                "1",
		AdapterInfo: discovery.AdapterInfo{
			Technology:            "lora-ttn",
			FwVersion:             "all",
			NetworkManagementType: "inclusion_exclusion",
		},
	}

}
