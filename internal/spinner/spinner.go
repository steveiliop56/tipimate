package spinner

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

var s *spinner.Spinner

func Init() {
	s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
}

func SetMessage(msg string) {
	s.Suffix = "msg"
}

func Start() {
	s.Start()
}

func Stop() {
	s.Stop()
}

func Succeed(msg string) {
	s.Stop()
	fmt.Printf("%s %s\n", color.GreenString("✔"), msg)
	s.Start()
}

func Fail(msg string) {
	s.Stop()
	fmt.Printf("%s %s\n", color.RedString("✘"), msg)
	s.Start()
}

func Warn(msg string) {
	s.Stop()
	fmt.Printf("%s %s\n", color.YellowString("⚠"), msg)
	s.Start()
}

func Info(msg string) {
	s.Stop()
	fmt.Printf("%s %s\n", color.BlueString("🛈"), msg)
	s.Start()
}

func Custom(msg string, icon string) {
	s.Stop()
	fmt.Printf("%s %s\n", icon, msg)
	s.Start()
}