package utils

import "gonum.org/v1/gonum/interp"

func Convert3vBatteryVoltageToLevel(voltage float64) uint16 {
	voltageTable := []float64{2,2.4,2.5,2.6,2.7,2.8,2.9,3}
	levelTable   := []float64{0,10,20,30,40,60,80,100}
	predictor := interp.PiecewiseLinear{}
	predictor.Fit(voltageTable,levelTable)
	return uint16(predictor.Predict(voltage))
}
