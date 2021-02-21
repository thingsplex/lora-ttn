package model

import (
	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/fimptype"
)

type ProductInfo struct {
	ProductHash string
	ProductAlias string
	ProductId    string
	ProductName  string
	ManufacturerId  string
	SwVersion    string
	PowerSource  string
}

type Driver interface {
	DecodeAndConvertToFimp(address string,payload []byte) []fimpgo.FimpMessage
	GetProductInfo() ProductInfo                            // for inclusion report
	GetSupportedServices(address string) []fimptype.Service // for inclusion report
}
