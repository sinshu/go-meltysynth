package meltysynth

import (
	"math"
)

var resonancePeakOffset = float32(1 - 1/math.Sqrt(2))

type biQuadFilter struct {
	synthesizer *Synthesizer
	active      bool
	a0          float32
	a1          float32
	a2          float32
	a3          float32
	a4          float32
	x1          float32
	x2          float32
	y1          float32
	y2          float32
}

func newBiQuadFilter(s *Synthesizer) *biQuadFilter {
	result := new(biQuadFilter)
	result.synthesizer = s
	return result
}

func (bf *biQuadFilter) clearBuffer() {
	bf.x1 = 0
	bf.x2 = 0
	bf.y1 = 0
	bf.y2 = 0
}

func (bf *biQuadFilter) setLowPassFilter(cutoffFrequency float32, resonance float32) {
	if cutoffFrequency >= 0.499*float32(bf.synthesizer.SampleRate) {
		bf.active = false
		return
	}
	bf.active = true

	// This equation gives the Q value which makes the desired resonance peak.
	// The error of the resultant peak height is less than 3%.
	q := resonance - resonancePeakOffset/(1+6*(resonance-1))

	w := 2 * math.Pi * float64(cutoffFrequency) / float64(bf.synthesizer.SampleRate)
	cosw := math.Cos(w)
	alpha := math.Sin(w) / float64(2*q)

	b0 := (1 - cosw) / 2
	b1 := 1 - cosw
	b2 := (1 - cosw) / 2
	a0 := 1 + alpha
	a1 := -2 * cosw
	a2 := 1 - alpha

	bf.setCoefficients(float32(a0), float32(a1), float32(a2), float32(b0), float32(b1), float32(b2))
}

func (bf *biQuadFilter) process(block []float32) {
	blockLength := len(block)

	if bf.active {
		for t := 0; t < blockLength; t++ {
			input := block[t]
			output := bf.a0*input + bf.a1*bf.x1 + bf.a2*bf.x2 - bf.a3*bf.y1 - bf.a4*bf.y2

			bf.x2 = bf.x1
			bf.x1 = input
			bf.y2 = bf.y1
			bf.y1 = output

			block[t] = output
		}
	} else {
		bf.x2 = block[blockLength-2]
		bf.x1 = block[blockLength-1]
		bf.y2 = bf.x2
		bf.y1 = bf.x1
	}
}

func (bf *biQuadFilter) setCoefficients(a0 float32, a1 float32, a2 float32, b0 float32, b1 float32, b2 float32) {
	bf.a0 = b0 / a0
	bf.a1 = b1 / a0
	bf.a2 = b2 / a0
	bf.a3 = a1 / a0
	bf.a4 = a2 / a0
}
