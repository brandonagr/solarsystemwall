package main

import (
	"fmt"

	solar "github.com/brandonagr/solarsystemwall/solar"
)

func main() {

	system := &solar.SolarSystem{}
	t := solar.NewCircle(system)

	fmt.Printf("Hello, world.\n %v", t)
}
