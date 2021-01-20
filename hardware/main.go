package main

import (
	"log"

	"github.com/gotmc/libusb"
)

func main() {
	ctx, _ := libusb.Init()
	defer ctx.Exit()
	devices, _ := ctx.GetDeviceList()
	for _, device := range devices {
		usbDeviceDescriptor, _ := device.GetDeviceDescriptor()
		handle, _ := device.Open()
		defer handle.Close()
		snIndex := usbDeviceDescriptor.SerialNumberIndex
		serialNumber, _ := handle.GetStringDescriptorASCII(snIndex)
		log.Printf("Found S/N: %s", serialNumber)
	}
}
