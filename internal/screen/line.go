package screen

import (
	"fmt"
	"github.com/yudppp/throttle"
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

func NewLine(device string, filename string) *Line {
	line := &Line{progress: 0, complete: false, device: device, filename: filename}
	lines = append(lines, line)
	return line
}

func (l *Line) SetProgress(progress int) {
	l.progress = progress
	throttleRefresh()
}

func (l *Line) SetComplete() {
	l.complete = true
	throttleRefresh()
}

var throttler throttle.Throttler

func getThrottler() *throttle.Throttler {
	if throttler == nil {
		var throttleDuration, err = time.ParseDuration("500ms")
		if err != nil {
			panic(err)
		}
		throttler = throttle.New(throttleDuration)
	}
	return &throttler
}

func throttleRefresh() {
	throttler := *getThrottler()
	throttler.Do(refresh)
}

func refresh() {
	clearScreen()
	for _, line := range lines {
		tick := fmt.Sprintf("%v%%", line.progress)
		if line.complete {
			tick = "✅ "
		}
		fmt.Printf("%v %v -> %v\n", tick, line.filename, line.device)
	}
}
