package hyperx

import (
	"log"
	"sync"
	"time"

	"github.com/sstallion/go-hid"
)

var device *hid.Device
var exitReadLoop bool = false
var waitGroup sync.WaitGroup

func CallFeatureReport(buffer [62]byte) (int, error) {
	featureReport := [62]byte{7}

	log.Println("Getting Feature Report", featureReport[0])

	_, err := device.GetFeatureReport(featureReport[:])

	if err != nil {
		return 0, err
	}

	log.Println("Sending Feature Report", buffer[0])

	return device.SendFeatureReport(buffer[:])
}

func PollStatusFeatures() {
	toCheck := [][62]byte{statusFeatureReport1, statusFeatureReport2, statusFeatureReport3, statusFeatureReport4, statusFeatureReport5}
	for index, statusFeatureReport := range toCheck {
		log.Println("Polling status feature", index)
		numBytes, err := CallFeatureReport(statusFeatureReport)
		if err != nil {
			log.Fatal("Oops")
		}
		log.Println("Got", numBytes, "bytes")
		time.Sleep(500 * time.Millisecond)
	}
}

func Run() {
	waitGroup.Add(1)
	err := hid.Init()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Init of hid complete.")

	device, err = hid.OpenFirst(VENDOR_ID, PRODUCT_ID)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		log.Println("Closing device")
		device.Close()
		waitGroup.Done()
	}()

	log.Println("Opened")

	productName, err := device.GetProductStr()
	if err != nil {
		log.Println("Couldn't get device name.", err)
	}

	log.Println("Opened device", productName)

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		PollStatusFeatures()
		for range ticker.C {
			PollStatusFeatures()
		}
	}()

	readBuffer := make([]byte, 64)

	log.Println("Entering read loop")
	for !exitReadLoop {
		numBytes, err := device.ReadWithTimeout(readBuffer, 1*time.Second)
		if err == nil && numBytes > 0 {
			log.Println("Read from HID", readBuffer, "length", numBytes)
			err = DecodePacket(readBuffer)
			if err != nil {
				log.Println("Error decoding the packet.", err)
			}
		}
	}
}

func Stop() {
	exitReadLoop = true
	waitGroup.Wait()
}
