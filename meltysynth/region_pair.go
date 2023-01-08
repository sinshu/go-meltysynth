package meltysynth

type regionPair struct {
	preset     *PresetRegion
	instrument *InstrumentRegion
}

func newRegionPair(preset *PresetRegion, inst *InstrumentRegion) regionPair {
	var result regionPair

	result.preset = preset
	result.instrument = inst

	return result
}

func (region regionPair) getGeneratorValue(generatorType uint16) int32 {
	return int32(region.preset.gs[generatorType]) + int32(region.instrument.gs[generatorType])
}

func (region regionPair) GetSampleStart() int32 {
	return region.instrument.GetSampleStart()
}

func (region regionPair) GetSampleEnd() int32 {
	return region.instrument.GetSampleEnd()
}

func (region regionPair) GetSampleStartLoop() int32 {
	return region.instrument.GetSampleStartLoop()
}

func (region regionPair) GetSampleEndLoop() int32 {
	return region.instrument.GetSampleEndLoop()
}

func (region regionPair) GetStartAddressOffset() int32 {
	return region.instrument.GetStartAddressOffset()
}

func (region regionPair) GetEndAddressOffset() int32 {
	return region.instrument.GetEndAddressOffset()
}

func (region regionPair) GetStartLoopAddressOffset() int32 {
	return region.instrument.GetStartLoopAddressOffset()
}

func (region regionPair) GetEndLoopAddressOffset() int32 {
	return region.instrument.GetEndLoopAddressOffset()
}

func (region regionPair) GetModulationLfoToPitch() int32 {
	return region.getGeneratorValue(gen_ModulationLfoToPitch)
}

func (region regionPair) GetVibratoLfoToPitch() int32 {
	return region.getGeneratorValue(gen_VibratoLfoToPitch)
}

func (region regionPair) GetModulationEnvelopeToPitch() int32 {
	return region.getGeneratorValue(gen_ModulationEnvelopeToPitch)
}

func (region regionPair) GetInitialFilterCutoffFrequency() float32 {
	return calcCentsToHertz(float32(region.getGeneratorValue(gen_InitialFilterCutoffFrequency)))
}

func (region regionPair) GetInitialFilterQ() float32 {
	return float32(0.1) * float32(region.getGeneratorValue(gen_InitialFilterQ))
}

func (region regionPair) GetModulationLfoToFilterCutoffFrequency() int32 {
	return region.getGeneratorValue(gen_ModulationLfoToFilterCutoffFrequency)
}

func (region regionPair) GetModulationEnvelopeToFilterCutoffFrequency() int32 {
	return region.getGeneratorValue(gen_ModulationEnvelopeToFilterCutoffFrequency)
}

func (region regionPair) GetModulationLfoToVolume() float32 {
	return float32(0.1) * float32(region.getGeneratorValue(gen_ModulationLfoToVolume))
}

func (region regionPair) GetChorusEffectsSend() float32 {
	return float32(0.1) * float32(region.getGeneratorValue(gen_ChorusEffectsSend))
}

func (region regionPair) GetReverbEffectsSend() float32 {
	return float32(0.1) * float32(region.getGeneratorValue(gen_ReverbEffectsSend))
}

func (region regionPair) GetPan() float32 {
	return float32(0.1) * float32(region.getGeneratorValue(gen_Pan))
}

func (region regionPair) GetDelayModulationLfo() float32 {
	return calcTimecentsToSeconds(float32(region.getGeneratorValue(gen_DelayModulationLfo)))
}

func (region regionPair) GetFrequencyModulationLfo() float32 {
	return calcCentsToHertz(float32(region.getGeneratorValue(gen_FrequencyModulationLfo)))
}

func (region regionPair) GetDelayVibratoLfo() float32 {
	return calcTimecentsToSeconds(float32(region.getGeneratorValue(gen_DelayVibratoLfo)))
}

func (region regionPair) GetFrequencyVibratoLfo() float32 {
	return calcCentsToHertz(float32(region.getGeneratorValue(gen_FrequencyVibratoLfo)))
}

func (region regionPair) GetDelayModulationEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.getGeneratorValue(gen_DelayModulationEnvelope)))
}

func (region regionPair) GetAttackModulationEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.getGeneratorValue(gen_AttackModulationEnvelope)))
}

func (region regionPair) GetHoldModulationEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.getGeneratorValue(gen_HoldModulationEnvelope)))
}

func (region regionPair) GetDecayModulationEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.getGeneratorValue(gen_DecayModulationEnvelope)))
}

func (region regionPair) GetSustainModulationEnvelope() float32 {
	return float32(0.1) * float32(region.getGeneratorValue(gen_SustainModulationEnvelope))
}

func (region regionPair) GetReleaseModulationEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.getGeneratorValue(gen_ReleaseModulationEnvelope)))
}

func (region regionPair) GetKeyNumberToModulationEnvelopeHold() int32 {
	return region.getGeneratorValue(gen_KeyNumberToModulationEnvelopeHold)
}

func (region regionPair) GetKeyNumberToModulationEnvelopeDecay() int32 {
	return region.getGeneratorValue(gen_KeyNumberToModulationEnvelopeDecay)
}

func (region regionPair) GetDelayVolumeEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.getGeneratorValue(gen_DelayVolumeEnvelope)))
}

func (region regionPair) GetAttackVolumeEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.getGeneratorValue(gen_AttackVolumeEnvelope)))
}

func (region regionPair) GetHoldVolumeEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.getGeneratorValue(gen_HoldVolumeEnvelope)))
}

func (region regionPair) GetDecayVolumeEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.getGeneratorValue(gen_DecayVolumeEnvelope)))
}

func (region regionPair) GetSustainVolumeEnvelope() float32 {
	return float32(0.1) * float32(region.getGeneratorValue(gen_SustainVolumeEnvelope))
}

func (region regionPair) GetReleaseVolumeEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.getGeneratorValue(gen_ReleaseVolumeEnvelope)))
}

func (region regionPair) GetKeyNumberToVolumeEnvelopeHold() int32 {
	return region.getGeneratorValue(gen_KeyNumberToVolumeEnvelopeHold)
}

func (region regionPair) GetKeyNumberToVolumeEnvelopeDecay() int32 {
	return region.getGeneratorValue(gen_KeyNumberToVolumeEnvelopeDecay)
}

func (region regionPair) GetInitialAttenuation() float32 {
	return float32(0.1) * float32(region.getGeneratorValue(gen_InitialAttenuation))
}

func (region regionPair) GetCoarseTune() int32 {
	return region.getGeneratorValue(gen_CoarseTune)
}

func (region regionPair) GetFineTune() int32 {
	return region.getGeneratorValue(gen_FineTune) + int32(region.instrument.Sample.PitchCorrection)
}

func (region regionPair) GetSampleModes() int32 {
	return region.instrument.GetSampleModes()
}

func (region regionPair) GetScaleTuning() int32 {
	return region.getGeneratorValue(gen_ScaleTuning)
}

func (region regionPair) GetExclusiveClass() int32 {
	return region.instrument.GetExclusiveClass()
}

func (region regionPair) GetRootKey() int32 {
	return region.instrument.GetRootKey()
}
