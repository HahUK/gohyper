## Summary

*gohyper* is a simple tool to enable battery low notifications and mute indication for the HyperX Cloud Flight S.

It features:

* Low battery notifications
* System Tray Icon:
  * Power off indicator
  * When powered on, microphone mute indicator
  * Menu item showing battery level, microphone and USB charging status
* Notifications configurable through command line arguments

## Usage

```
Usage of gohyper:
  -interval duration
        Interval between battery checks (default 2m0s)
  -notifications
        Enable low battery notifications (default true)
  -threshold value
        Percentage low battery threshold (default 40%)

```
## Dependencies

* [go-hid](https://github.com/sstallion/go-hid) for cross platform USB HID communication
* [beeep](https://github.com/gen2brain/beeep) for cross platform notifications
* [systray](github.com/getlantern/systray) for cross platform system tray icon and menu


## Notes

It is written in Go and uses USB HID to communicate with the headset dongle.

The protocol has been reverse engineered from usb captures so mistakes are likely however it does seem to work. Only tested on Linux.

HyperX and the HyperX logo are registered trademarks or trademarks of HP Inc. and/or Kingston Technology Corporation in the U.S. and/or other countries. All registered trademarks and trademarks are property of their respective owners. 