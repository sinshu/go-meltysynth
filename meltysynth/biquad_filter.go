package meltysynth

type biQuadFilter struct {
	synthesizer *Synthesizer
}

func newBiQuadFilter(synthesizer *Synthesizer) *biQuadFilter {
	result := new(biQuadFilter)
	result.synthesizer = synthesizer
	return result
}

func (filter *biQuadFilter) clearBuffer() {
}

func (filter *biQuadFilter) setLowPassFilter(cutoffFrequency float32, resonance float32) {
}
