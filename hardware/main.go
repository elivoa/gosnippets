package main

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"rpkg.cc/times"
)

func main() {

	a := base64.StdEncoding.EncodeToString([]byte(strings.Join([]string{"app", "26d2fd1c-93f1-45ce-9411-f39c5b6422d4", "cc651deeeedd7dd5fc150d4aabf1a58d028e36f8"}, ":")))
	b := strings.Join([]string{"Bearer", a}, " ")
	fmt.Println("...", b)

	c := time.Now().In(times.GetDefaultLoc()).UnixNano() / 1e6
	fmt.Println("....", c)

	// ctx, _ := libusb.Init()
	// defer ctx.Exit()
	// devices, _ := ctx.GetDeviceList()
	// for _, device := range devices {
	// 	usbDeviceDescriptor, _ := device.GetDeviceDescriptor()
	// 	handle, _ := device.Open()
	// 	defer handle.Close()
	// 	snIndex := usbDeviceDescriptor.SerialNumberIndex
	// 	serialNumber, _ := handle.GetStringDescriptorASCII(snIndex)
	// 	log.Printf("Found S/N: %s", serialNumber)
	// }
}
