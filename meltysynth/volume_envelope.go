package meltysynth

import "math"

type volumeEnvelope struct {
	synthesizer          *Synthesizer
	attackSlope          float64
	decaySlope           float64
	releaseSlope         float64
	attackStartTime      float64
	holdStartTime        float64
	decayStartTime       float64
	releaseStartTime     float64
	sustainLevel         float32
	releaseLevel         float32
	processedSampleCount int32
	stage                int32
	value                float32
	priority             float32
}

func newVolumeEnvelope(s *Synthesizer) *volumeEnvelope {
	result := new(volumeEnvelope)
	result.synthesizer = s
	return result
}

func (env *volumeEnvelope) start(delay float32, attack float32, hold float32, decay float32, sustain float32, release float32) {
	env.attackSlope = 1 / float64(attack)
	env.decaySlope = -9.226 / float64(decay)
	env.releaseSlope = -9.226 / float64(release)

	env.attackStartTime = float64(delay)
	env.holdStartTime = env.attackStartTime + float64(attack)
	env.decayStartTime = env.holdStartTime + float64(hold)
	env.releaseStartTime = 0

	env.sustainLevel = calcClamp(sustain, 0, 1)
	env.releaseLevel = 0

	env.processedSampleCount = 0
	env.stage = env_Delay
	env.value = 0

	env.process(0)
}

func (env *volumeEnvelope) release() {
	env.stage = env_Release
	env.releaseStartTime = float64(env.processedSampleCount) / float64(env.synthesizer.SampleRate)
	env.releaseLevel = env.value
}

func (env *volumeEnvelope) process(sampleCount int32) bool {
	env.processedSampleCount += sampleCount

	currentTime := float64(env.processedSampleCount) / float64(env.synthesizer.SampleRate)

	for env.stage <= env_Hold {
		var endTime float64
		switch env.stage {
		case env_Delay:
			endTime = env.attackStartTime
		case env_Attack:
			endTime = env.holdStartTime
		case env_Hold:
			endTime = env.decayStartTime
		default:
			panic("invalid envelope stage")
		}

		if currentTime < endTime {
			break
		}
		env.stage++
	}

	switch env.stage {
	case env_Delay:
		env.value = 0
		env.priority = 4 + env.value
		return true
	case env_Attack:
		env.value = float32(env.attackSlope * (currentTime - env.attackStartTime))
		env.priority = 3 + env.value
		return true
	case env_Hold:
		env.value = 1
		env.priority = 2 + env.value
		return true
	case env_Decay:
		env.value = float32(math.Max(calcExpCutoff(env.decaySlope*(currentTime-env.decayStartTime)), float64(env.sustainLevel)))
		env.priority = 1 + env.value
		return env.value > nonAudible
	case env_Release:
		env.value = float32(float64(env.releaseLevel) * calcExpCutoff(env.releaseSlope*(currentTime-env.releaseStartTime)))
		env.priority = env.value
		return env.value > nonAudible
	default:
		panic("invalid envelope stage")
	}
}
