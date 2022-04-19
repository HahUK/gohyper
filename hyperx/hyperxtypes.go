package hyperx

type PacketTypes byte

const (
	INPUT_PACKET  PacketTypes = 0x02
	STATUS_PACKET PacketTypes = 0x0b
	AUDIO_PACKET  PacketTypes = 0x0a
)

type InputValueTypes byte

const (
	NOOP     InputValueTypes = 0x00
	VOL_UP   InputValueTypes = 0x01
	VOL_DOWN InputValueTypes = 0x02
)

type Input struct {
	PacketType PacketTypes
	Value      InputValueTypes
}

type StatusUpdateTypes byte

const (
	POWER_STATE   StatusUpdateTypes = 0x01
	BATTERY_LEVEL StatusUpdateTypes = 0x02
	USB_CHARGING  StatusUpdateTypes = 0x03
)

type PowerUpdateTypes byte

const (
	OFF PowerUpdateTypes = 0x00
	ON  PowerUpdateTypes = 0x01
)

type StatusUpdate struct {
	PacketType   PacketTypes
	_            byte
	_            byte
	UpdateType   StatusUpdateTypes
	PowerUpdate  PowerUpdateTypes
	_            byte
	_            byte
	BatteryLevel byte
	_            byte
	_            byte
	_            byte
	_            byte
}

type AudioUpdateTypes byte

const (
	AUDIO_UPDATE AudioUpdateTypes = 0x03
)

type MuteAndMonitorTypes byte

const (
	MUTE_OFF_AND_MONITOR_OFF MuteAndMonitorTypes = 0x00
	MUTE_ON_AND_MONITOR_OFF  MuteAndMonitorTypes = 0x02
	MUTE_OFF_AND_MONITOR_ON  MuteAndMonitorTypes = 0x10
	MUTE_ON_AND_MONITOR_ON   MuteAndMonitorTypes = 0x12
)

type ChannelStatusTypes byte

const (
	CHANNEL_OFF ChannelStatusTypes = 0x00
	CHANNEL_ON  ChannelStatusTypes = 0x01
)

type SurroundStatusTypes byte

const (
	SURROUND_OFF SurroundStatusTypes = 0x00
	SURROUND_ON  SurroundStatusTypes = 0x02
)

type AudioUpdate struct {
	PacketType        PacketTypes
	_                 byte
	Surround          SurroundStatusTypes
	UpdateType        AudioUpdateTypes
	MuteAndMonitor    MuteAndMonitorTypes
	_                 byte
	GameChannelLevel  byte
	GameChannelStatus ChannelStatusTypes
	_                 byte
	_                 byte
	ChatLevel         byte
	ChatStatus        ChannelStatusTypes
}

type AudioCallback func(AudioUpdate)
type InputCallback func(Input)
type StatusCallback func(StatusUpdate)
