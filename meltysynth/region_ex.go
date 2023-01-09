package meltysynth

import "math"

func (o *oscillator) startByRegion(data []int16, region regionPair) {
	sampleRate := region.instrument.Sample.SampleRate
	loopMode := region.GetSampleModes()
	sampleStart := region.GetSampleStart()
	sampleEnd := region.GetSampleEnd()
	startLoop := region.GetSampleStartLoop()
	endLoop := region.GetSampleEndLoop()
	rootKey := region.GetRootKey()
	coarseTune := region.GetCoarseTune()
	fineTune := region.GetFineTune()
	scaleTuning := region.GetScaleTuning()

	o.start(data, loopMode, sampleRate, sampleStart, sampleEnd, startLoop, endLoop, rootKey, coarseTune, fineTune, scaleTuning)
}

func (env *volumeEnvelope) startByRegion(region regionPair, key int32, velocity int32) {
	// If the release time is shorter than 10 ms, it will be clamped to 10 ms to avoid pop noise.

	delay := region.GetDelayVolumeEnvelope()
	attack := region.GetAttackVolumeEnvelope()
	hold := region.GetHoldVolumeEnvelope() * calcKeyNumberToMultiplyingFactor(region.GetKeyNumberToVolumeEnvelopeHold(), key)
	decay := region.GetDecayVolumeEnvelope() * calcKeyNumberToMultiplyingFactor(region.GetKeyNumberToVolumeEnvelopeDecay(), key)
	sustain := calcDecibelsToLinear(-region.GetSustainVolumeEnvelope())
	release := float32(math.Max(float64(region.GetReleaseVolumeEnvelope()), 0.01))

	env.start(delay, attack, hold, decay, sustain, release)
}

func (env *modulationEnvelope) startByRegion(region regionPair, key int32, velocity int32) {
	// According to the implementation of TinySoundFont, the attack time should be adjusted by the velocity.

	delay := region.GetDelayModulationEnvelope()
	attack := region.GetAttackModulationEnvelope() * (float32(145-velocity) / 144)
	hold := region.GetHoldModulationEnvelope() * calcKeyNumberToMultiplyingFactor(region.GetKeyNumberToModulationEnvelopeHold(), key)
	decay := region.GetDecayModulationEnvelope() * calcKeyNumberToMultiplyingFactor(region.GetKeyNumberToModulationEnvelopeDecay(), key)
	sustain := 1 - region.GetSustainModulationEnvelope()/100
	release := region.GetReleaseModulationEnvelope()

	env.start(delay, attack, hold, decay, sustain, release)
}

func (lfo *lfo) startVibrato(region regionPair, key int32, velocity int32) {
	lfo.start(region.GetDelayVibratoLfo(), region.GetFrequencyVibratoLfo())
}

func (lfo *lfo) startModulation(region regionPair, key int32, velocity int32) {
	lfo.start(region.GetDelayModulationLfo(), region.GetFrequencyModulationLfo())
}
