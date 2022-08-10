package meltysynth

type regionPair struct {
	preset     *PresetRegion
	instrument *InstrumentRegion
}

func newRegionPair(preset *PresetRegion, instrument *InstrumentRegion) regionPair {

	var result regionPair

	result.preset = preset
	result.instrument = instrument

	return result
}

func getGeneratorValue(region regionPair, generatorType uint16) int32 {
	return int32(region.preset.gs[generatorType]) + int32(region.instrument.gs[generatorType])
}

func GetRegionPairSampleStart(region regionPair) int32 {
	return GetInstrumentSampleStart(region.instrument)
}

func GetRegionPairSampleEnd(region regionPair) int32 {
	return GetInstrumentSampleEnd(region.instrument)
}

func GetRegionPairSampleStartLoop(region regionPair) int32 {
	return GetInstrumentSampleStartLoop(region.instrument)
}

func GetRegionPairSampleEndLoop(region regionPair) int32 {
	return GetInstrumentSampleEndLoop(region.instrument)
}

func GetRegionPairStartAddressOffset(region regionPair) int32 {
	return GetInstrumentStartAddressOffset(region.instrument)
}

func GetRegionPairEndAddressOffset(region regionPair) int32 {
	return GetInstrumentEndAddressOffset(region.instrument)
}

func GetRegionPairStartLoopAddressOffset(region regionPair) int32 {
	return GetInstrumentStartLoopAddressOffset(region.instrument)
}

func GetRegionPairEndLoopAddressOffset(region regionPair) int32 {
	return GetInstrumentEndLoopAddressOffset(region.instrument)
}

func GetRegionPairModulationLfoToPitch(region regionPair) int32 {
	return getGeneratorValue(region, gen_ModulationLfoToPitch)
}

func GetRegionPairVibratoLfoToPitch(region regionPair) int32 {
	return getGeneratorValue(region, gen_VibratoLfoToPitch)
}

func GetRegionPairModulationEnvelopeToPitch(region regionPair) int32 {
	return getGeneratorValue(region, gen_ModulationEnvelopeToPitch)
}

func GetRegionPairInitialFilterCutoffFrequency(region regionPair) float32 {
	return calcCentsToHertz(float32(getGeneratorValue(region, gen_InitialFilterCutoffFrequency)))
}

func GetRegionPairInitialFilterQ(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_InitialFilterQ))
}

func GetRegionPairModulationLfoToFilterCutoffFrequency(region regionPair) int32 {
	return getGeneratorValue(region, gen_ModulationLfoToFilterCutoffFrequency)
}

func GetRegionPairModulationEnvelopeToFilterCutoffFrequency(region regionPair) int32 {
	return getGeneratorValue(region, gen_ModulationEnvelopeToFilterCutoffFrequency)
}

func GetRegionPairModulationLfoToVolume(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_ModulationLfoToVolume))
}

func GetRegionPairChorusEffectsSend(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_ChorusEffectsSend))
}

func GetRegionPairReverbEffectsSend(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_ReverbEffectsSend))
}

func GetRegionPairPan(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_Pan))
}

func GetRegionPairDelayModulationLfo(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_DelayModulationLfo)))
}

func GetRegionPairFrequencyModulationLfo(region regionPair) float32 {
	return calcCentsToHertz(float32(getGeneratorValue(region, gen_FrequencyModulationLfo)))
}

func GetRegionPairDelayVibratoLfo(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_DelayVibratoLfo)))
}

func GetRegionPairFrequencyVibratoLfo(region regionPair) float32 {
	return calcCentsToHertz(float32(getGeneratorValue(region, gen_FrequencyVibratoLfo)))
}

func GetRegionPairDelayModulationEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_DelayModulationEnvelope)))
}

func GetRegionPairAttackModulationEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_AttackModulationEnvelope)))
}

func GetRegionPairHoldModulationEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_HoldModulationEnvelope)))
}

func GetRegionPairDecayModulationEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_DecayModulationEnvelope)))
}

func GetRegionPairSustainModulationEnvelope(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_SustainModulationEnvelope))
}

func GetRegionPairReleaseModulationEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_ReleaseModulationEnvelope)))
}

func GetRegionPairKeyNumberToModulationEnvelopeHold(region regionPair) int32 {
	return getGeneratorValue(region, gen_KeyNumberToModulationEnvelopeHold)
}

func GetRegionPairKeyNumberToModulationEnvelopeDecay(region regionPair) int32 {
	return getGeneratorValue(region, gen_KeyNumberToModulationEnvelopeDecay)
}

func GetRegionPairDelayVolumeEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_DelayVolumeEnvelope)))
}

func GetRegionPairAttackVolumeEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_AttackVolumeEnvelope)))
}

func GetRegionPairHoldVolumeEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_HoldVolumeEnvelope)))
}

func GetRegionPairDecayVolumeEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_DecayVolumeEnvelope)))
}

func GetRegionPairSustainVolumeEnvelope(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_SustainVolumeEnvelope))
}

func GetRegionPairReleaseVolumeEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_ReleaseVolumeEnvelope)))
}

func GetRegionPairKeyNumberToVolumeEnvelopeHold(region regionPair) int32 {
	return getGeneratorValue(region, gen_KeyNumberToVolumeEnvelopeHold)
}

func GetRegionPairKeyNumberToVolumeEnvelopeDecay(region regionPair) int32 {
	return getGeneratorValue(region, gen_KeyNumberToVolumeEnvelopeDecay)
}

func GetRegionPairInitialAttenuation(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_InitialAttenuation))
}

func GetRegionPairCoarseTune(region regionPair) int32 {
	return getGeneratorValue(region, gen_CoarseTune)
}

func GetRegionPairFineTune(region regionPair) int32 {
	return getGeneratorValue(region, gen_FineTune) + int32(region.instrument.Sample.PitchCorrection)
}

func GetRegionPairSampleModes(region regionPair) int32 {
	return GetInstrumentSampleModes(region.instrument)
}

func GetRegionPairScaleTuning(region regionPair) int32 {
	return getGeneratorValue(region, gen_ScaleTuning)
}

func GetRegionPairExclusiveClass(region regionPair) int32 {
	return GetInstrumentExclusiveClass(region.instrument)
}

func GetRegionPairRootKey(region regionPair) int32 {
	return GetInstrumentRootKey(region.instrument)
}
