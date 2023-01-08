package meltysynth

func arrayMultiplyAdd(a float32, x []float32, dst []float32) {
	dstLen := len(dst)
	for i := 0; i < dstLen; i++ {
		dst[i] += a * x[i]
	}
}

func arrayMultiplyAddSlope(a float32, step float32, x []float32, dst []float32) {
	dstLen := len(dst)
	for i := 0; i < dstLen; i++ {
		dst[i] += a * x[i]
		a += step
	}
}
