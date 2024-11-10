package adbwrapper

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
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

type Progress struct {
	PercentComplete int // 0-100
	Done            bool
}

func CopyFileToDevice(device Device, src string, progressout chan<- Progress) error {
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
		lines := make(chan string)
		go readLines(lines, pipe)
		for {
			select {
			case <-stop:
				return
			case line := <-lines:
				//fmt.Println("line1", line)
				if progress, err := extractFirstNumber(line); err == nil {
					//fmt.Println("progress", progress)
					progressout <- Progress{PercentComplete: progress, Done: false}
				}
			}
		}
	}(pipe, stopSignal)

	err = cmd.Wait()
	stopSignal <- true
	progressout <- Progress{PercentComplete: 100, Done: true}
	close(progressout)
	if err != nil {
		return err
	}

	return nil
}

func readLines(lines chan<- string, pipe io.ReadCloser) {
	defer close(lines)
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines <- scanner.Text()
	}
}

func extractFirstNumber(str string) (int, error) {
	re := regexp.MustCompile(`\d+%`)
	match := re.FindString(str)
	if match == "" {
		return 0, fmt.Errorf("no number found in string")
	}
	match = strings.TrimSuffix(match, "%")
	num, err := strconv.Atoi(match)
	if err != nil {
		return 0, err
	}
	return num, nil
}
