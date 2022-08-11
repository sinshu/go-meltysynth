package meltysynth

import (
	"math"
)

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
	position         float64
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

	if loopMode == 0 {
		oscillator.looping = false
	} else {
		oscillator.looping = true
	}

	oscillator.position = float64(start)
}

func (oscillator *oscillator) release() {

	if oscillator.loopMode == 3 {
		oscillator.looping = false
	}
}

func (oscillator *oscillator) process(block []float32, pitch float32) bool {

	pitchChange := oscillator.pitchChangeScale*(pitch-float32(oscillator.rootKey)) + oscillator.tune
	pitchRatio := float64(oscillator.sampleRateRatio) * math.Pow(float64(2), float64(pitchChange)/float64(12))
	return oscillator.fillBlock(block, pitchRatio)
}

func (oscillator *oscillator) fillBlock(block []float32, pitchRatio float64) bool {

	if oscillator.looping {
		return oscillator.fillBlock_Continuous(block, pitchRatio)
	} else {
		return oscillator.fillBlock_NoLoop(block, pitchRatio)
	}
}

func (oscillator *oscillator) fillBlock_NoLoop(block []float32, pitchRatio float64) bool {

	blockLength := len(block)

	for t := 0; t < blockLength; t++ {

		index := int32(oscillator.position)

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
		a := float32(oscillator.position - float64(index))
		block[t] = (float32(x1) + a*float32(x2-x1)) / 32768

		oscillator.position += pitchRatio
	}

	return true
}

func (oscillator *oscillator) fillBlock_Continuous(block []float32, pitchRatio float64) bool {

	blockLength := len(block)

	endLoopPosition := float64(oscillator.endLoop)

	loopLength := oscillator.endLoop - oscillator.startLoop

	for t := 0; t < blockLength; t++ {

		if oscillator.position >= endLoopPosition {
			oscillator.position -= float64(loopLength)
		}

		index1 := int32(oscillator.position)
		index2 := index1 + 1

		if index2 >= oscillator.endLoop {
			index2 -= loopLength
		}

		x1 := oscillator.data[index1]
		x2 := oscillator.data[index2]
		a := oscillator.position - float64(index1)
		block[t] = float32((float64(x1) + a*float64(x2-x1)) / 32768)

		oscillator.position += pitchRatio
	}

	return true
}
