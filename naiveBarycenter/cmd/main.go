package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		os.Exit(1)
	}

	numBodies, errStr := strconv.Atoi(args[1])
	if errStr != nil {
		os.Exit(1)
	}

	rand.Seed(time.Now().Unix())

	posMax := 100
	massMax := 5

	for i := 0; i < numBodies; i++ {
		posX := rand.Intn(posMax*2) - posMax // it is not possible to generate negative numbers
		posY := rand.Intn(posMax*2) - posMax
		posZ := rand.Intn(posMax*2) - posMax

		// if our mass is zero our model does not make any sense
		// and we have to ensure the mass never be less tha one or more than our mass max
		mass := rand.Intn(massMax-1) + 1
		fmt.Printf("%d:%d:%d:%d\n", posX, posY, posZ, mass)
	}
}

