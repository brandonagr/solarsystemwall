package main

import (
	"fmt"

	solar "github.com/brandonagr/solarsystemwall/solar"
)

func main() {

	system := &solar.System{}
	t := solar.NewCircle(system)

	fmt.Printf("Hello, world.\n %v", t)

	solar.NewWebDisplay(system, 640, 240)

	// keep running
	fmt.Scanln()
}
