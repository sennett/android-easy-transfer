package screen

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var clearFuncs = map[string]func(){
	"linux": func() {
		cmd := exec.Command("clearScreen") //Linux example, its tested
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	},
	"windows": func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	},
	"darwin": func() {
		cmd := exec.Command("clear") // osx
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	},
}

func clearScreen() {
	value, ok := clearFuncs[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                               //if we defined a clearScreen func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic(fmt.Sprintf("Your platform is unsupported! I can't clear terminal screen (platform %v)\n", runtime.GOOS))
	}
}
