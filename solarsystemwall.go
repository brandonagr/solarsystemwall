package main

import (
	"fmt"

	solar "github.com/brandonagr/solarsystemwall/solar"
)

func main() {

	system := solar.DefaultSystem()

	//t := solar.NewCircle(system)
	//fmt.Printf("Hello, world.\n %v", t)

	solar.NewWebDisplay(system, 136, 91)

	// keep running
	fmt.Scanln()
}
