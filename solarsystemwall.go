package main

import (
	"time"

	solar "github.com/brandonagr/solarsystemwall/solar"
)

func main() {
	system := solar.DefaultSystem()
	display := solar.NewWebDisplay(system, 136, 91)

	runAnimationLoopForever(system, display)
}

func runAnimationLoopForever(system *solar.System, display *solar.WebDisplay) {
	curTime := time.Now()
	prevTime := curTime

	ticks := time.NewTicker(time.Duration(100.0) * time.Millisecond) // 10 hz
	defer ticks.Stop()

	for _ = range ticks.C {
		prevTime, curTime = curTime, time.Now()
		dt := curTime.Sub(prevTime).Seconds()

		system.Animate(dt)
		display.Render(system)
	}
}
