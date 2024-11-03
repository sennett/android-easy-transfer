package adbwrapper

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Device struct {
	Name     string
	WriteDir string
}

func FetchDevices() ([]Device, error) {
	out, err := exec.Command("adb", "devices").Output()
	if err != nil {
		return nil, err
	}

	devices := DevicesStdOutToDevices(out)

	return devices, nil
}

func CheckFolderExists(device Device) bool {
	_, err := exec.Command("adb", "-s", device.Name, "shell", "ls", device.WriteDir).Output()
	if err != nil {
		return false
	}
	return true
}

func DevicesStdOutToDevices(output []byte) []Device {
	scanner := bufio.NewScanner(bytes.NewReader(output))
	devices := []Device{}
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		if text == "List of devices attached" {
			continue
		}

		deviceName := strings.Fields(text)[0]

		devices = append(devices, Device{Name: deviceName, WriteDir: "/storage/self/primary/Download"})
	}

	for _, device := range devices {
		//thing := &devices[i]
		if device.Name == "" {
			panic("Found empty device name - exiting")
		}
	}

	return devices
}

func CopyFileToDevice(device Device, src string) error {
	pathComponents := strings.Split(src, "/")
	fileName := pathComponents[len(pathComponents)-1]
	fmt.Printf("---> %v -> %v:%v...", fileName, device.Name, device.WriteDir)
	res, err := exec.Command("adb", "-s", device.Name, "push", src, device.WriteDir).CombinedOutput()
	if err != nil {
		return errors.New(string(res))
	}
	fmt.Println("done")
	return nil
}
