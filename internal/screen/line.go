package screen

import (
	"fmt"
	"github.com/boz/go-throttle"
	"github.com/fatih/color"
	"strings"
	"time"
)

type Screen struct {
}

type Line struct {
	filename string
	device   string
	progress int
	complete bool
}

var lines = make([]*Line, 0)

func NewLine(device string, filenameAndPath string) *Line {
	pathComponents := strings.Split(filenameAndPath, "/")
	fileName := pathComponents[len(pathComponents)-1]
	line := &Line{progress: 0, complete: false, device: device, filename: fileName}
	lines = append(lines, line)
	return line
}

func (l *Line) SetProgress(progress int) {
	l.progress = progress
	throttleRefresh.Trigger()
}

func (l *Line) SetComplete() {
	l.complete = true
	throttleRefresh.Trigger()
}

func getDuration() time.Duration {
	duration, _ := time.ParseDuration("500ms")
	return duration
}

var throttleRefresh = throttle.ThrottleFunc(getDuration(), true, refresh)

func refresh() {
	clearScreen()
	for i, line := range lines {
		progressBarWidth := 20
		progressChars := ""
		for i := 0; i < progressBarWidth; i++ {
			if i < line.progress*progressBarWidth/100 {
				progressChars += "█"
			} else {
				progressChars += "_"
			}
		}
		tick := fmt.Sprintf("%v%%", line.progress)
		if line.complete {
			tick = "✅ "
		}

		colorIndex := i % len(colorArray)
		_, _ = colorArray[colorIndex].Printf("%v %v %v -> %v\n", progressChars, tick, line.filename, line.device)
		//fmt.Printf("%v %v %v -> %v\n", progressChars, tick, line.filename, line.device)
	}
}

var colorArray = []*color.Color{
	color.New(color.FgHiRed),
	color.RGB(255, 165, 0), // orange
	color.RGB(200, 200, 0), // yellow
	color.New(color.FgHiGreen),
	color.New(color.FgHiCyan), // blue
	color.RGB(64, 0, 255),     // indigo
	color.RGB(143, 0, 255),    // violet
	color.New(color.FgBlack),  // black because Macky
}
