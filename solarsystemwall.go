package main

import (
	"runtime"
	"time"
	"fmt"

	solar "github.com/brandonagr/solarsystemwall/solar"
)

func main() {
	system := solar.DefaultSystem()
	display := solar.NewDisplay(system)
	defer display.Dispose()

	fmt.Println("Beginning Animation")

	runAnimationLoopForever(system, display)
}

func runAnimationLoopForever(system *solar.System, display solar.Display) {
	curTime := time.Now()
	prevTime := curTime

	var ticks *time.Ticker
	if runtime.GOOS == "windows" {
		ticks = time.NewTicker(time.Duration(10.0) * time.Millisecond) // 10 hz
	} else {
		ticks = time.NewTicker(time.Duration(100.0) * time.Millisecond) // 100 hz
	}
	defer ticks.Stop()

	for _ = range ticks.C {
		prevTime, curTime = curTime, time.Now()
		dt := curTime.Sub(prevTime).Seconds()

		system.Animate(dt)
		display.Render(system)
	}
}
