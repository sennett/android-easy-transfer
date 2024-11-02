package main

import (
	"fmt"
	"log"
)

func main() {
	devices, err := FetchDevices()

	if err != nil {
		fmt.Println("Could not fetch devices")
		fmt.Println(err)
		return
	}

	if len(devices) == 0 {
		fmt.Println("No devices found - exiting")
		return
	}

	for _, device := range devices {
		device.WriteDir = "/storage/self/primary/Download/"
		fmt.Println(device)
		writePossible := CheckFolderExists(device)
		if !writePossible {
			fmt.Printf("Could not write to dir %v on %v\n", device.WriteDir, device.Name)
			return
		}
	}

	fmt.Println(devices)

	for _, device := range devices {
		cpErr := CopyFileToDevice(device, "/Users/sennett/Downloads/vlc.log")
		if cpErr != nil {
			fmt.Println("Could not copy file to device")
			log.Fatal(cpErr)
			return
		}
	}

}
