package meltysynth

func arrayMultiplyAdd(a float32, x []float32, destination []float32) {

	destinationLength := len(destination)
	for i := 0; i < destinationLength; i++ {
		destination[i] += a * x[i]
	}
}

func arrayMultiplyAddSlope(a float32, step float32, x []float32, destination []float32) {

	destinationLength := len(destination)
	for i := 0; i < destinationLength; i++ {
		destination[i] += a * x[i]
		a += step
	}
}
