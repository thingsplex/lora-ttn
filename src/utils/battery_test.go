package utils

import "testing"

func TestConvert3vBatteryVoltageToLevel(t *testing.T) {
	t.Log(Convert3vBatteryVoltageToLevel(3.1))
	t.Log(Convert3vBatteryVoltageToLevel(3))
	t.Log(Convert3vBatteryVoltageToLevel(2.9))
	t.Log(Convert3vBatteryVoltageToLevel(2.85))
	t.Log(Convert3vBatteryVoltageToLevel(2.75))
	t.Log(Convert3vBatteryVoltageToLevel(2.4))
	t.Log(Convert3vBatteryVoltageToLevel(2.3))
	t.Log(Convert3vBatteryVoltageToLevel(2.2))
	t.Log(Convert3vBatteryVoltageToLevel(2.1))
	t.Log(Convert3vBatteryVoltageToLevel(2))
	t.Log(Convert3vBatteryVoltageToLevel(1.9))
}
