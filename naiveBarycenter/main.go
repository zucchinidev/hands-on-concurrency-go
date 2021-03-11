package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

// MassPoint represents a body
type MassPoint struct {
	x, y, z, mass float64
}

func addMassPoints(a, b MassPoint) MassPoint {
	return MassPoint{
		x:    a.x + b.x,
		y:    a.y + b.y,
		z:    a.z + b.z,
		mass: a.mass + b.mass,
	}
}

// maps bodies into the mass sensitive subspace
func avgMassPoints(a, b MassPoint) MassPoint {
	sum := addMassPoints(a, b)
	return MassPoint{
		x:    sum.x / 2,
		y:    sum.y / 2,
		z:    sum.z / 2,
		mass: sum.mass,
	}
}

// Now we need functions that can map points to and from the math sensitive subspace
// We'll call those to weighted subspace and from weighted subspace
func toWeightedSubspace(a MassPoint) MassPoint {
	return MassPoint{
		x:    a.x * a.mass,
		y:    a.y * a.mass,
		z:    a.z * a.mass,
		mass: a.mass,
	}
}

func fromWeightedSubspace(a MassPoint) MassPoint {
	return MassPoint{
		x:    a.x / a.mass,
		y:    a.y / a.mass,
		z:    a.z / a.mass,
		mass: a.mass,
	}
}

func avgMassPointsWeighted(a, b MassPoint) MassPoint {
	aWeighted := toWeightedSubspace(a)
	bWeighted := toWeightedSubspace(b)
	return fromWeightedSubspace(avgMassPoints(aWeighted, bWeighted))
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Incorrect number of arguments")
		os.Exit(1)
	}

	f, errOpeningFile := os.Open(args[1])
	if errOpeningFile != nil {
		fmt.Printf("Error opening file %s\n", errOpeningFile)
		os.Exit(1)
	}
	defer f.Close()

	var massPoints []MassPoint
	startLoading := time.Now()
	for {
		var newMasPoint MassPoint
		_, err := fmt.Fscanf(f, "%f:%f:%f:%f", &newMasPoint.x, &newMasPoint.y, &newMasPoint.z, &newMasPoint.mass)

		if err != nil && err == io.EOF {
			break
		}

		if err != nil {
			continue // malformed line
		}

		massPoints = append(massPoints, newMasPoint)
	}

	// Now we'll report how many points we loaded
	fmt.Printf("Loaded %d values from files in %s.\n", len(massPoints), time.Since(startLoading))

	// So if Len MassPoints is less than or equal to 1,
	// of course we can't find the barycenter of one MassPoint because it's just that point,
	// so if it's less than or equal to one then we will simply create a new error
	if len(massPoints) <= 1 {
		fmt.Printf("Insufficient values found %d\n", len(massPoints))
		os.Exit(1)
	}

	startCalculation := time.Now()

	for len(massPoints) != 1 {
		var newMassPoints []MassPoint
		// we do not want run off the end
		for i := 0; i < len(massPoints)-1; i += 2 {
			// dealing with pairs of MassPoints each time
			newMassPoints = append(newMassPoints, avgMassPointsWeighted(massPoints[i], massPoints[i+1]))
		}

		// Now because we only check that we haven't run off the end,
		// we don't actually check that we handle all of them, we need to do that.
		// So if there are an odd number, we'll simply take the last one and put it into the new array
		if len(massPoints)%2 != 0 {
			newMassPoints = append(newMassPoints, massPoints[len(massPoints)-1])
		}

		// Now we only need to switch out the set of MassPoints and the loop will run again.
		// So MassPoints equals new MassPoints, we're swapping this out and this will reduce
		// by half every single time.
		massPoints = newMassPoints
	}

	// And once the loop actually finishes, we need to get the system average,
	// which is going to equal the zeroth element, the first element.
	// So that will be the only one remaining, because remember our loop
	// condition here is just checking that there are still more than one left.
	// So as soon as there's only one left it will go down here and we'll get the only element of the array.
	systemAverage := massPoints[0]

	fmt.Printf("System barycenter is at (%f, %f, %f) and the system's mass is %f.\n",
		systemAverage.x,
		systemAverage.y,
		systemAverage.z,
		systemAverage.mass)
	fmt.Printf("Calculation took %s", time.Since(startCalculation))
}
