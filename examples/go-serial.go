package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mag/apps/kerrigan/audiocollector/xiaolufawn"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

// TODO 发现设备，"/dev/cu.usbmodem3202191983011" 这个是怎么来的？
// TODO 在序号后面多了个1？这个是多设备的识别？mac下需要列目录了。
// TODO 1. 找到设备。
// TODO 2. 串口打开设备。
func main2() {
	fmt.Println("-------------------------------------------------")

	options := serial.OpenOptions{
		PortName: "/dev/tty.usbmodem3202191983011",
		// PortName:              "/dev/cu.usbmodem3202191983011",
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 8,
		// InterCharacterTimeout: 100,
		// Rs485RxDuringTx: true,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	// Make sure to close it later.
	defer port.Close()

	n, err := port.Write([]byte(xiaolufawn.MsgUDISKENTER))
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}

	fmt.Println("Wrote", n, "bytes. wait for done")

	var resp bytes.Buffer
	for {
		// fmt.Print(".")
		buf := make([]byte, 32)
		n, err := port.Read(buf)
		if err != nil && err != io.EOF {
			// FAILED
			fmt.Println("Error reading from serial port: ", err)
			panic(err)
		}
		if _, err := resp.Write(buf[:n]); err != nil {
			panic(err)
		}
		if err == io.EOF {
			// fmt.Println("EOF??: ", err)
			break
		}
	}
	fmt.Println(" -> ", resp.String())

	// for {
	// 	fmt.Print(".")
	// 	buf := make([]byte, 8)
	// 	n, err := port.Read(buf)
	// 	if err != nil {
	// 		if err != io.EOF {
	// 			fmt.Println("Error reading from serial port: ", err)
	// 		}
	// 	} else {
	// 		buf = buf[:n]
	// 		fmt.Println("Rx: ", hex.EncodeToString(buf))
	// 	}
	// }

	fmt.Println("Wait 3 second to stop.")
	time.Sleep(3 * time.Second)
}

func read(port io.ReadWriteCloser) {
	// * READ
	for {
		fmt.Print(".")
		buf := make([]byte, 12)
		n, err := port.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from serial port: ", err)
			}
		} else {
			buf = buf[:n]
			fmt.Println("Rx: ", hex.EncodeToString(buf))
		}
	}

}
