package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"

	"github.com/hahuk/gohyper/hyperx"
)

var NOTIFICATIONS_ENABLED bool = true
var NOTIFICATIONS_INTERVAL time.Duration = 2 * time.Minute
var LOW_BATTERY_THRESHOLD uint8 = 40

var LastBatteryLevel uint8 = 255
var LastWasPoweredOn bool = false
var LastWasUSBCharging bool = false
var LastMuteWasOn bool = false

var MenuItemBatteryLevel *systray.MenuItem
var MenuItemMuted *systray.MenuItem

//go:embed assets/closedmic-white-nospace.png
var IconDataMuted []byte

//go:embed assets/openmic-white-nospace.png
var IconDataOpen []byte

//go:embed assets/off-white.png
var IconDataOff []byte

func myAudioCallback(au hyperx.AudioUpdate) {
	log.Println("audio", au)
	if au.MuteAndMonitor == hyperx.MUTE_ON_AND_MONITOR_OFF || au.MuteAndMonitor == hyperx.MUTE_ON_AND_MONITOR_ON {
		LastMuteWasOn = true
	} else {
		LastMuteWasOn = false
	}

	go updateSysTray()
}

func myStatusCallback(su hyperx.StatusUpdate) {
	log.Println("status", su)
	if su.UpdateType == hyperx.BATTERY_LEVEL {
		LastBatteryLevel = uint8(su.BatteryLevel)
	}
	if su.UpdateType == hyperx.POWER_STATE {
		LastWasPoweredOn = su.PowerUpdate == hyperx.ON
	}
	if su.UpdateType == hyperx.USB_CHARGING {
		LastWasUSBCharging = su.PowerUpdate == hyperx.ON
	}

	go updateSysTray()
}

func myInputCallback(input hyperx.Input) {
	log.Println("input", input)
}

func checkBattery() {
	log.Println("LastBatteryLevel", LastBatteryLevel)
	log.Println("LastWasPoweredOn", LastWasPoweredOn)
	log.Println("LastWasUSBCharging", LastWasUSBCharging)
	if LastBatteryLevel < LOW_BATTERY_THRESHOLD && LastWasPoweredOn && !LastWasUSBCharging {
		batteryErrorString := fmt.Sprint("Headset battery is at ", LastBatteryLevel, "%, please charge it.")
		err := beeep.Alert("Headset Battery Low", batteryErrorString, "assets/warning.png")
		if err != nil {
			log.Println("Cannot create beeep alert for low battery", err)
		}
	}
}

func updateSysTray() {
	if LastWasPoweredOn {
		var micStatus = "open"
		var micIcon = IconDataOpen

		if LastMuteWasOn {
			micStatus = "muted"
			micIcon = IconDataMuted
		}

		var IsBatteryLow = LastBatteryLevel < LOW_BATTERY_THRESHOLD

		generatedIcon, err := getBatteryLevelIcon(micIcon,
			LastBatteryLevel,
			LastWasUSBCharging,
			IsBatteryLow)

		if err != nil {
			log.Println("Error creating battery level icon", err)
		} else {
			micIcon = generatedIcon
		}

		var batteryStatus = ""
		if LastWasUSBCharging {
			batteryStatus = " (Charging)"
		}

		micString := fmt.Sprint("Mic: ", micStatus)
		batteryString := fmt.Sprint("Battery Level: ", LastBatteryLevel, batteryStatus)
		tooltipString := fmt.Sprint(micString, ", ", batteryString)

		MenuItemBatteryLevel.SetTitle(batteryString)
		MenuItemMuted.SetTitle(micString)

		systray.SetTooltip(tooltipString)
		systray.SetIcon(micIcon)
	} else {
		systray.SetTooltip("Powered Off")
		systray.SetIcon(IconDataOff)
	}
}

func onSysTrayReady() {
	hyperx.RegisterAudioCallback(myAudioCallback)
	hyperx.RegisterInputCallback(myInputCallback)
	hyperx.RegisterStatusCallback(myStatusCallback)

	systray.SetIcon(IconDataOff)
	systray.SetTitle("HyperX Headset Monitor")

	MenuItemBatteryLevel = systray.AddMenuItem("Battery Level", "Battery Level")
	MenuItemMuted = systray.AddMenuItem("Mic", "Mic Status")
	MenuItemQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		<-MenuItemQuit.ClickedCh
		log.Println("Quitting from menu item")
		systray.Quit()
	}()

	go hyperx.Run()
}

func onSysTrayExit() {
	hyperx.Stop()
}

func main() {
	flag.BoolVar(&NOTIFICATIONS_ENABLED, "notifications", NOTIFICATIONS_ENABLED, "Enable low battery notifications")
	flag.DurationVar(&NOTIFICATIONS_INTERVAL, "interval", NOTIFICATIONS_INTERVAL, "Interval between battery checks")
	flag.Func("threshold", "Percentage low battery threshold (default 40%)", func(flagValue string) error {
		tempThreshold, err := strconv.ParseUint(flagValue, 10, 8)
		if err != nil {
			return err
		}
		if tempThreshold > 100 {
			return errors.New("battery percentage cannot be more than 100")
		}
		LOW_BATTERY_THRESHOLD = uint8(tempThreshold)
		return nil
	})
	flag.Parse()

	log.Println("LOW_BATTERY_THRESHOLD", LOW_BATTERY_THRESHOLD)
	log.Println("NOTIFICATIONS_ENABLED", NOTIFICATIONS_ENABLED)
	log.Println("NOTIFICATIONS_INTERVAL", NOTIFICATIONS_INTERVAL)

	if NOTIFICATIONS_ENABLED {
		batteryTicker := time.NewTicker(NOTIFICATIONS_INTERVAL)

		go func() {
			checkBattery()
			for range batteryTicker.C {
				checkBattery()
			}
		}()
	}

	systray.Run(onSysTrayReady, onSysTrayExit)
}
