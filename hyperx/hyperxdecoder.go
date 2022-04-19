package hyperx

import (
	"bytes"
	"encoding/binary"
)

var audioCallbacks []AudioCallback
var statusCallbacks []StatusCallback
var inputCallbacks []InputCallback

func RegisterAudioCallback(audioCallback AudioCallback) {
	audioCallbacks = append(audioCallbacks, audioCallback)
}

func RegisterStatusCallback(statusCallback StatusCallback) {
	statusCallbacks = append(statusCallbacks, statusCallback)
}

func RegisterInputCallback(inputCallback InputCallback) {
	inputCallbacks = append(inputCallbacks, inputCallback)
}

func DecodePacket(packet []byte) error {
	switch packet[0] {
	case byte(INPUT_PACKET):
		var InputPacket Input
		packetReader := bytes.NewReader(packet)
		err := binary.Read(packetReader, binary.LittleEndian, &InputPacket)
		if err != nil {
			return err
		}
		for _, inputCallback := range inputCallbacks {
			go inputCallback(InputPacket)
		}
	case byte(STATUS_PACKET):
		var StatusPacket StatusUpdate
		packetReader := bytes.NewReader(packet)
		err := binary.Read(packetReader, binary.LittleEndian, &StatusPacket)
		if err != nil {
			return err
		}
		for _, statusCallback := range statusCallbacks {
			go statusCallback(StatusPacket)
		}
	case byte(AUDIO_PACKET):
		var AudioPacket AudioUpdate
		packetReader := bytes.NewReader(packet)
		err := binary.Read(packetReader, binary.LittleEndian, &AudioPacket)
		if err != nil {
			return err
		}
		for _, audioCallback := range audioCallbacks {
			go audioCallback(AudioPacket)
		}
	}
	return nil
}
