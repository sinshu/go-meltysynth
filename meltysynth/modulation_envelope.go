package meltysynth

import "math"

type modulationEnvelope struct {
	synthesizer          *Synthesizer
	attackSlope          float64
	decaySlope           float64
	releaseSlope         float64
	attackStartTime      float64
	holdStartTime        float64
	decayStartTime       float64
	decayEndTime         float64
	releaseEndTime       float64
	sustainLevel         float32
	releaseLevel         float32
	processedSampleCount int32
	stage                int32
	value                float32
}

func newModulationEnvelope(synthesizer *Synthesizer) *modulationEnvelope {
	result := new(modulationEnvelope)
	result.synthesizer = synthesizer
	return result
}

func (envelope *modulationEnvelope) start(delay float32, attack float32, hold float32, decay float32, sustain float32, release float32) {

	envelope.attackSlope = 1 / float64(attack)
	envelope.decaySlope = 1 / float64(decay)
	envelope.releaseSlope = 1 / float64(release)

	envelope.attackStartTime = float64(delay)
	envelope.holdStartTime = envelope.attackStartTime + float64(attack)
	envelope.decayStartTime = envelope.holdStartTime + float64(hold)

	envelope.decayEndTime = envelope.decayStartTime + float64(decay)
	envelope.releaseEndTime = float64(release)

	envelope.sustainLevel = calcClamp(sustain, 0, 1)
	envelope.releaseLevel = 0

	envelope.processedSampleCount = 0
	envelope.stage = env_Delay
	envelope.value = 0

	envelope.process(0)
}

func (envelope *modulationEnvelope) release() {

	envelope.stage = env_Release
	envelope.releaseEndTime += float64(envelope.processedSampleCount) / float64(envelope.synthesizer.SampleRate)
	envelope.releaseLevel = envelope.value
}

func (envelope *modulationEnvelope) process(sampleCount int32) bool {

	envelope.processedSampleCount += sampleCount

	currentTime := float64(envelope.processedSampleCount) / float64(envelope.synthesizer.SampleRate)

	for envelope.stage <= env_Hold {

		var endTime float64
		switch envelope.stage {

		case env_Delay:
			endTime = envelope.attackStartTime

		case env_Attack:
			endTime = envelope.holdStartTime

		case env_Hold:
			endTime = envelope.decayStartTime

		default:
			panic("invalid envelope stage")
		}

		if currentTime < endTime {
			break
		} else {
			envelope.stage++
		}
	}

	switch envelope.stage {

	case env_Delay:
		envelope.value = 0
		return true

	case env_Attack:
		envelope.value = float32(envelope.attackSlope * (currentTime - envelope.attackStartTime))
		return true

	case env_Hold:
		envelope.value = 1
		return true

	case env_Decay:
		envelope.value = float32(math.Max(envelope.decaySlope*(envelope.decayEndTime-currentTime), float64(envelope.sustainLevel)))
		return envelope.value > nonAudible

	case env_Release:
		envelope.value = float32(math.Max(float64(envelope.releaseLevel)*float64(envelope.releaseSlope)*(envelope.releaseEndTime-currentTime), 0))
		return envelope.value > nonAudible

	default:
		panic("invalid envelope stage.")
	}
}
