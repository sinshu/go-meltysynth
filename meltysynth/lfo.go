package meltysynth

import "math"

type lfo struct {
	synthesizer          *Synthesizer
	active               bool
	delay                float64
	period               float64
	processedSampleCount int32
	value                float32
}

func newLfo(synthesizer *Synthesizer) *lfo {
	result := new(lfo)
	result.synthesizer = synthesizer
	return result
}

func (lfo *lfo) begin(delay float32, frequency float32) {

	if frequency > 1.0e-3 {

		lfo.active = true

		lfo.delay = float64(delay)
		lfo.period = 1.0 / float64(frequency)

		lfo.processedSampleCount = 0
		lfo.value = 0

	} else {

		lfo.active = false
		lfo.value = 0
	}
}

func (lfo *lfo) process() {

	if !lfo.active {
		return
	}

	lfo.processedSampleCount += lfo.synthesizer.BlockSize

	currentTime := float64(lfo.processedSampleCount) / float64(lfo.synthesizer.SampleRate)

	if currentTime < lfo.delay {

		lfo.value = 0

	} else {

		phase := math.Mod(currentTime-lfo.delay, lfo.period) / lfo.period

		if phase < 0.25 {
			lfo.value = float32(4 * phase)
		} else if phase < 0.75 {
			lfo.value = float32(4 * (0.5 - phase))
		} else {
			lfo.value = float32(4 * (phase - 1.0))
		}
	}
}