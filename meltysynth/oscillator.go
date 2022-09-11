package meltysynth

import (
	"math"
)

// In this class, fixed-point numbers are used for speed-up.
// A fixed-point number is expressed by Int64, whose lower 24 bits represent the fraction part,
// and the rest represent the integer part.
// For clarity, fixed-point number variables have a suffix "_fp".

const fracBits int32 = 24
const fracUnit int64 = 1 << fracBits
const fpToSample float32 = float32(1) / float32(32768*fracUnit)

type oscillator struct {
	synthesizer      *Synthesizer
	data             []int16
	loopMode         int32
	sampleRate       int32
	sampleStart      int32
	sampleEnd        int32
	startLoop        int32
	endLoop          int32
	rootKey          int32
	tune             float32
	pitchChangeScale float32
	sampleRateRatio  float32
	looping          bool
	position_fp      int64
}

func newOscillator(synthesizer *Synthesizer) *oscillator {
	result := new(oscillator)
	result.synthesizer = synthesizer
	return result
}

func (oscillator *oscillator) start(data []int16, loopMode int32, sampleRate int32, start int32, end int32, startLoop int32, endLoop int32, rootKey int32, coarseTune int32, fineTune int32, scaleTuning int32) {

	oscillator.data = data
	oscillator.loopMode = loopMode
	oscillator.sampleRate = sampleRate
	oscillator.sampleStart = start
	oscillator.sampleEnd = end
	oscillator.startLoop = startLoop
	oscillator.endLoop = endLoop
	oscillator.rootKey = rootKey

	oscillator.tune = float32(coarseTune) + float32(0.01)*float32(fineTune)
	oscillator.pitchChangeScale = float32(0.01) * float32(scaleTuning)
	oscillator.sampleRateRatio = float32(sampleRate) / float32(oscillator.synthesizer.SampleRate)

	if loopMode == loop_NoLoop {
		oscillator.looping = false
	} else {
		oscillator.looping = true
	}

	oscillator.position_fp = int64(start) << fracBits
}

func (oscillator *oscillator) release() {

	if oscillator.loopMode == loop_LoopUntilNoteOff {
		oscillator.looping = false
	}
}

func (oscillator *oscillator) process(block []float32, pitch float32) bool {

	pitchChange := oscillator.pitchChangeScale*(pitch-float32(oscillator.rootKey)) + oscillator.tune
	pitchRatio := float64(oscillator.sampleRateRatio) * math.Pow(float64(2), float64(pitchChange)/float64(12))
	return oscillator.fillBlock(block, pitchRatio)
}

func (oscillator *oscillator) fillBlock(block []float32, pitchRatio float64) bool {

	pitchRatio_fp := int64(float64(fracUnit) * pitchRatio)

	if oscillator.looping {
		return oscillator.fillBlock_Continuous(block, pitchRatio_fp)
	} else {
		return oscillator.fillBlock_NoLoop(block, pitchRatio_fp)
	}
}

func (oscillator *oscillator) fillBlock_NoLoop(block []float32, pitchRatio_fp int64) bool {

	blockLength := len(block)

	for t := 0; t < blockLength; t++ {

		index := int32(oscillator.position_fp >> fracBits)

		if index >= oscillator.sampleEnd {
			if t > 0 {
				for i := t; i < blockLength; i++ {
					block[i] = 0
				}
				return true
			} else {
				return false
			}
		}

		x1 := oscillator.data[index]
		x2 := oscillator.data[index+1]
		a_fp := oscillator.position_fp & (fracUnit - 1)
		block[t] = fpToSample * float32((int64(x1)<<fracBits)+a_fp*int64(x2-x1))

		oscillator.position_fp += pitchRatio_fp
	}

	return true
}

func (oscillator *oscillator) fillBlock_Continuous(block []float32, pitchRatio_fp int64) bool {

	blockLength := len(block)

	endLoop_fp := int64(oscillator.endLoop) << fracBits

	loopLength := int32(oscillator.endLoop - oscillator.startLoop)
	loopLength_fp := int64(loopLength) << fracBits

	for t := 0; t < blockLength; t++ {

		if oscillator.position_fp >= endLoop_fp {
			oscillator.position_fp -= loopLength_fp
		}

		index1 := int32(oscillator.position_fp >> fracBits)
		index2 := index1 + 1

		if index2 >= oscillator.endLoop {
			index2 -= loopLength
		}

		x1 := oscillator.data[index1]
		x2 := oscillator.data[index2]
		a_fp := oscillator.position_fp & (fracUnit - 1)
		block[t] = fpToSample * float32((int64(x1)<<fracBits)+a_fp*int64(x2-x1))

		oscillator.position_fp += pitchRatio_fp
	}

	return true
}
