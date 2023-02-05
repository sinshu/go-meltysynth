package meltysynth

import "math"

const (
	fixedGain    = 0.015
	scaleWet     = 3.0
	scaleDamp    = 0.4
	scaleRoom    = 0.28
	offsetRoom   = 0.7
	initialRoom  = 0.5
	initialDamp  = 0.5
	initialWet   = 1.0 / scaleWet
	initialWidth = 1.0
	stereoSpread = 23

	cfTuningL1  = 1116
	cfTuningR1  = 1116 + stereoSpread
	cfTuningL2  = 1188
	cfTuningR2  = 1188 + stereoSpread
	cfTuningL3  = 1277
	cfTuningR3  = 1277 + stereoSpread
	cfTuningL4  = 1356
	cfTuningR4  = 1356 + stereoSpread
	cfTuningL5  = 1422
	cfTuningR5  = 1422 + stereoSpread
	cfTuningL6  = 1491
	cfTuningR6  = 1491 + stereoSpread
	cfTuningL7  = 1557
	cfTuningR7  = 1557 + stereoSpread
	cfTuningL8  = 1617
	cfTuningR8  = 1617 + stereoSpread
	apfTuningL1 = 556
	apfTuningR1 = 556 + stereoSpread
	apfTuningL2 = 441
	apfTuningR2 = 441 + stereoSpread
	apfTuningL3 = 341
	apfTuningR3 = 341 + stereoSpread
	apfTuningL4 = 225
	apfTuningR4 = 225 + stereoSpread
)

type reverb struct {
	cfsL  []*combFilter
	cfsR  []*combFilter
	apfsL []*allPassFilter
	apfsR []*allPassFilter

	gain      float32
	roomSize  float32
	roomSize1 float32
	damp      float32
	damp1     float32
	wet       float32
	wet1      float32
	wet2      float32
	width     float32
}

func newReverb(sampleRate int32) *reverb {
	cfsL := make([]*combFilter, 8)
	cfsL[0] = newCombFilter(scaleTuning(sampleRate, cfTuningL1))
	cfsL[1] = newCombFilter(scaleTuning(sampleRate, cfTuningL2))
	cfsL[2] = newCombFilter(scaleTuning(sampleRate, cfTuningL3))
	cfsL[3] = newCombFilter(scaleTuning(sampleRate, cfTuningL4))
	cfsL[4] = newCombFilter(scaleTuning(sampleRate, cfTuningL5))
	cfsL[5] = newCombFilter(scaleTuning(sampleRate, cfTuningL6))
	cfsL[6] = newCombFilter(scaleTuning(sampleRate, cfTuningL7))
	cfsL[7] = newCombFilter(scaleTuning(sampleRate, cfTuningL8))

	cfsR := make([]*combFilter, 8)
	cfsR[0] = newCombFilter(scaleTuning(sampleRate, cfTuningR1))
	cfsR[1] = newCombFilter(scaleTuning(sampleRate, cfTuningR2))
	cfsR[2] = newCombFilter(scaleTuning(sampleRate, cfTuningR3))
	cfsR[3] = newCombFilter(scaleTuning(sampleRate, cfTuningR4))
	cfsR[4] = newCombFilter(scaleTuning(sampleRate, cfTuningR5))
	cfsR[5] = newCombFilter(scaleTuning(sampleRate, cfTuningR6))
	cfsR[6] = newCombFilter(scaleTuning(sampleRate, cfTuningR7))
	cfsR[7] = newCombFilter(scaleTuning(sampleRate, cfTuningR8))

	apfsL := make([]*allPassFilter, 4)
	apfsL[0] = newAllPassFilter(scaleTuning(sampleRate, apfTuningL1))
	apfsL[1] = newAllPassFilter(scaleTuning(sampleRate, apfTuningL2))
	apfsL[2] = newAllPassFilter(scaleTuning(sampleRate, apfTuningL3))
	apfsL[3] = newAllPassFilter(scaleTuning(sampleRate, apfTuningL4))

	apfsR := make([]*allPassFilter, 4)
	apfsR[0] = newAllPassFilter(scaleTuning(sampleRate, apfTuningR1))
	apfsR[1] = newAllPassFilter(scaleTuning(sampleRate, apfTuningR2))
	apfsR[2] = newAllPassFilter(scaleTuning(sampleRate, apfTuningR3))
	apfsR[3] = newAllPassFilter(scaleTuning(sampleRate, apfTuningR4))

	for i := 0; i < len(apfsL); i++ {
		apfsL[i].setFeedback(0.5)
	}

	for i := 0; i < len(apfsR); i++ {
		apfsR[i].setFeedback(0.5)
	}

	result := new(reverb)
	result.cfsL = cfsL
	result.cfsR = cfsR
	result.apfsL = apfsL
	result.apfsR = apfsR
	result.setWet(initialWet)
	result.setRoomSize(initialRoom)
	result.setDamp(initialDamp)
	result.setWidth(initialWidth)
	return result
}

func (r *reverb) mute() {
	for i := 0; i < len(r.cfsL); i++ {
		r.cfsL[i].mute()
	}

	for i := 0; i < len(r.cfsR); i++ {
		r.cfsR[i].mute()
	}

	for i := 0; i < len(r.apfsL); i++ {
		r.apfsL[i].mute()
	}

	for i := 0; i < len(r.apfsR); i++ {
		r.apfsR[i].mute()
	}
}

func scaleTuning(sampleRate int32, tuning int) int {
	return int(math.Round(float64(sampleRate) / 44100.0 * float64(tuning)))
}

func (r *reverb) process(input []float32, outputLeft []float32, outputRight []float32) {
	length := len(input)

	for t := 0; t < length; t++ {
		outputLeft[t] = 0
	}
	for t := 0; t < length; t++ {
		outputRight[t] = 0
	}

	for i := 0; i < len(r.cfsL); i++ {
		r.cfsL[i].process(input, outputLeft)
	}

	for i := 0; i < len(r.apfsL); i++ {
		r.apfsL[i].process(outputLeft)
	}

	for i := 0; i < len(r.cfsR); i++ {
		r.cfsR[i].process(input, outputRight)
	}

	for i := 0; i < len(r.apfsR); i++ {
		r.apfsR[i].process(outputRight)
	}

	// With the default settings, we can skip this part.
	if 1.0-r.wet1 > 1.0e-3 || r.wet2 > 1.0e-3 {
		for t := 0; t < length; t++ {
			left := outputLeft[t]
			right := outputRight[t]
			outputLeft[t] = left*r.wet1 + right*r.wet2
			outputRight[t] = right*r.wet1 + left*r.wet2
		}
	}
}

func (r *reverb) update() {
	r.wet1 = r.wet * (r.width/2.0 + 0.5)
	r.wet2 = r.wet * ((1.0 - r.width) / 2.0)

	r.roomSize1 = r.roomSize
	r.damp1 = r.damp
	r.gain = fixedGain

	for i := 0; i < len(r.cfsL); i++ {
		r.cfsL[i].setFeedback(r.roomSize1)
		r.cfsL[i].setDamp(r.damp1)
	}

	for i := 0; i < len(r.cfsR); i++ {
		r.cfsR[i].setFeedback(r.roomSize1)
		r.cfsR[i].setDamp(r.damp1)
	}
}

func (r *reverb) getInputGain() float32 {
	return r.gain
}

func (r *reverb) setRoomSize(value float32) {
	r.roomSize = (value * scaleRoom) + offsetRoom
	r.update()
}

func (r *reverb) setDamp(value float32) {
	r.damp = value * scaleDamp
	r.update()
}

func (r *reverb) setWet(value float32) {
	r.wet = value * scaleWet
	r.update()
}

func (r *reverb) setWidth(value float32) {
	r.width = value
	r.update()
}

type combFilter struct {
	buffer []float32

	bufferIndex int
	filterStore float32

	feedback float32
	damp1    float32
	damp2    float32
}

func newCombFilter(bufferSize int) *combFilter {
	result := new(combFilter)
	result.buffer = make([]float32, bufferSize)
	result.bufferIndex = 0
	result.filterStore = 0
	result.feedback = 0
	result.damp1 = 0
	result.damp2 = 0
	return result
}

func (cf *combFilter) mute() {
	bufLen := len(cf.buffer)
	for i := 0; i < bufLen; i++ {
		cf.buffer[i] = 0
	}
}

func (cf *combFilter) process(inputBlock []float32, outputBlock []float32) {
	bufferLength := len(cf.buffer)
	outputBlockLength := len(outputBlock)

	blockIndex := 0
	for blockIndex < outputBlockLength {
		if cf.bufferIndex == bufferLength {
			cf.bufferIndex = 0
		}

		srcRem := bufferLength - cf.bufferIndex
		dstRem := outputBlockLength - blockIndex
		rem := int(math.Min(float64(srcRem), float64(dstRem)))

		for t := 0; t < rem; t++ {
			blockPos := blockIndex + t
			bufferPos := cf.bufferIndex + t

			input := inputBlock[blockPos]

			// The following ifs are to avoid performance problem due to denormalized number.
			// These ifs are equivalent to "if abs(value) < 1.0E-6".

			output := cf.buffer[bufferPos]
			if (math.Float32bits(output) & 0x7FFFFFFF) < 897988541 {
				output = 0
			}

			cf.filterStore = (output * cf.damp2) + (cf.filterStore * cf.damp1)
			if (math.Float32bits(cf.filterStore) & 0x7FFFFFFF) < 897988541 {
				cf.filterStore = 0
			}

			cf.buffer[bufferPos] = input + (cf.filterStore * cf.feedback)
			outputBlock[blockPos] += output
		}

		cf.bufferIndex += rem
		blockIndex += rem
	}
}

func (cf *combFilter) setFeedback(value float32) {
	cf.feedback = value
}

func (cf *combFilter) setDamp(value float32) {
	cf.damp1 = value
	cf.damp2 = 1 - value
}

type allPassFilter struct {
	buffer []float32

	bufferIndex int

	feedback float32
}

func newAllPassFilter(bufferSize int) *allPassFilter {
	result := new(allPassFilter)
	result.buffer = make([]float32, bufferSize)
	result.bufferIndex = 0
	result.feedback = 0
	return result
}

func (apf *allPassFilter) mute() {
	bufLen := len(apf.buffer)
	for i := 0; i < bufLen; i++ {
		apf.buffer[i] = 0
	}
}

func (apf *allPassFilter) process(block []float32) {
	bufferLength := len(apf.buffer)
	blockLength := len(block)

	var blockIndex = 0
	for blockIndex < blockLength {
		if apf.bufferIndex == bufferLength {
			apf.bufferIndex = 0
		}

		srcRem := bufferLength - apf.bufferIndex
		dstRem := blockLength - blockIndex
		rem := int(math.Min(float64(srcRem), float64(dstRem)))

		for t := 0; t < rem; t++ {
			blockPos := blockIndex + t
			bufferPos := apf.bufferIndex + t

			input := block[blockPos]

			bufout := apf.buffer[bufferPos]
			if (math.Float32bits(bufout) & 0x7FFFFFFF) < 897988541 {
				bufout = 0
			}

			block[blockPos] = bufout - input
			apf.buffer[bufferPos] = input + (bufout * apf.feedback)
		}

		apf.bufferIndex += rem
		blockIndex += rem
	}
}

func (apf *allPassFilter) setFeedback(value float32) {
	apf.feedback = value
}
