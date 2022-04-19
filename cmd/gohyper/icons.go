package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"sync"
)

const IMAGE_WIDTH = 48
const IMAGE_HEIGHT = 48

var chargingColour color.RGBA = color.RGBA{29, 172, 214, 255}
var lowColour color.RGBA = color.RGBA{255, 0, 0, 255}
var normalColour color.RGBA = color.RGBA{0, 255, 0, 255}
var remainingColour color.RGBA = color.RGBA{255, 255, 255, 255}

func getColourToUse(batteryLevel uint8, IsUSBCharging bool, ISBatteryLow bool) color.RGBA {
	if IsUSBCharging {
		return chargingColour
	}

	if ISBatteryLow {
		return lowColour
	}

	return normalColour
}

func generateMaskImage(srcIcon []byte) (image.Image, error) {
	tmpImage, _, tmpError := image.Decode(bytes.NewReader(srcIcon))
	return tmpImage, tmpError
}

func generateColourImage(colourToUse color.RGBA, percentage uint8) image.Image {
	// Create a blank, remainingColour image
	tmpImage := image.NewRGBA(image.Rect(0, 0, IMAGE_WIDTH, IMAGE_HEIGHT))
	draw.Draw(tmpImage, tmpImage.Bounds(), &image.Uniform{remainingColour}, image.Point{}, draw.Src)

	// Figure out how much should be coloured based on battery percentage
	percentageHeight := (float64(IMAGE_HEIGHT) / float64(100)) * float64(percentage)
	percentageBounds := image.Rect(0, 48, IMAGE_WIDTH, 48-int(percentageHeight))

	// Draw the coloured part
	draw.Draw(tmpImage, percentageBounds.Bounds(), &image.Uniform{colourToUse}, image.Point{}, draw.Src)
	return tmpImage
}

func generateBatteryLevelIcon(srcIcon []byte, batteryLevel uint8, IsUSBCharging bool, ISBatteryLow bool) ([]byte, error) {
	var colourToUse = getColourToUse(batteryLevel, IsUSBCharging, ISBatteryLow)
	var wg sync.WaitGroup

	var maskImage image.Image
	var colourImage image.Image
	var outputImage *image.RGBA

	var maskError error

	wg.Add(2)

	go func() {
		defer wg.Done()

		maskImage, maskError = generateMaskImage(srcIcon)
	}()

	go func() {
		defer wg.Done()

		colourImage = generateColourImage(colourToUse, batteryLevel)
	}()

	wg.Wait()

	if maskError != nil {
		return nil, maskError
	}

	outputImage = image.NewRGBA(image.Rect(0, 0, IMAGE_WIDTH, IMAGE_HEIGHT))
	draw.DrawMask(outputImage, outputImage.Bounds(), colourImage, image.Point{}, maskImage, image.Point{}, draw.Src)

	var tmpBuffer = new(bytes.Buffer)

	err := png.Encode(tmpBuffer, outputImage)
	return tmpBuffer.Bytes(), err
}

func getBatteryLevelIcon(srcIcon []byte, batteryLevel uint8, IsUSBCharging bool, ISBatteryLow bool) ([]byte, error) {
	if batteryLevel > 100 {
		batteryLevel = 100
	}

	return generateBatteryLevelIcon(srcIcon, batteryLevel, IsUSBCharging, ISBatteryLow)
}
