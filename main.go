package main

import (
	"adb-push-everywhere/internal/adbwrapper"
	"adb-push-everywhere/internal/watcher"
	"fmt"
	"log"
	"os"
	"path"
)

func main() {
	devices, err := adbwrapper.FetchDevices()

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
		writePossible := adbwrapper.CheckFolderExists(device)
		if !writePossible {
			fmt.Printf("Could not write to dir %v on %v\n", device.WriteDir, device.Name)
			return
		}
	}

	fmt.Printf("Found %v devices:\n", len(devices))
	for _, device := range devices {
		fmt.Println(device.Name)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fileTarget := path.Join(cwd, "filetarget")

	fmt.Printf("Watching for new files in %v...", fileTarget)

	// watch for files
	newFile := make(chan watcher.ChanPayload)
	go func() {
		err := watcher.WatchDir(fileTarget, newFile)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for {
		fileEvent := <-newFile
		for _, device := range devices {
			go copyFile(device, fileEvent.Filepath)
		}
	}
}

func copyFile(device adbwrapper.Device, file string) {
	err := adbwrapper.CopyFileToDevice(device, file)
	if err != nil {
		fmt.Printf("Could not copy file %v to device %v\n", file, device.Name)
		log.Fatal(err)
		return
	}
}
