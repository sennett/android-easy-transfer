package adbwrapper

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
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
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	wrapper := path.Join(wd, "redirectstdout.sh")

	cmd := exec.Command(wrapper, "adb", "-s", device.Name, "push", "-p", src, device.WriteDir)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	stopSignal := make(chan bool)

	go func(p io.ReadCloser, stop chan bool) {
		reader := bufio.NewReader(pipe)
		for {
			select {
			case <-stop:
				return
			default:
				line, _ := reader.ReadString('\r')
				fmt.Println(line)
				fmt.Println(len(line))
				fmt.Println("here")
				fmt.Println(err)
			}
		}
	}(pipe, stopSignal)

	err = cmd.Wait()
	fmt.Println("stopping")
	stopSignal <- true
	if err != nil {
		return err
	}

	fmt.Println("done")
	return nil
}
