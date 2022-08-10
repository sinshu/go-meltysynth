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

func GetregionPairSampleStart(region regionPair) int32 {
	return GetInstrumentSampleStart(region.instrument)
}

func GetregionPairSampleEnd(region regionPair) int32 {
	return GetInstrumentSampleEnd(region.instrument)
}

func GetregionPairSampleStartLoop(region regionPair) int32 {
	return GetInstrumentSampleStartLoop(region.instrument)
}

func GetregionPairSampleEndLoop(region regionPair) int32 {
	return GetInstrumentSampleEndLoop(region.instrument)
}

func GetregionPairStartAddressOffset(region regionPair) int32 {
	return GetInstrumentStartAddressOffset(region.instrument)
}

func GetregionPairEndAddressOffset(region regionPair) int32 {
	return GetInstrumentEndAddressOffset(region.instrument)
}

func GetregionPairStartLoopAddressOffset(region regionPair) int32 {
	return GetInstrumentStartLoopAddressOffset(region.instrument)
}

func GetregionPairEndLoopAddressOffset(region regionPair) int32 {
	return GetInstrumentEndLoopAddressOffset(region.instrument)
}

func GetregionPairModulationLfoToPitch(region regionPair) int32 {
	return getGeneratorValue(region, gen_ModulationLfoToPitch)
}

func GetregionPairVibratoLfoToPitch(region regionPair) int32 {
	return getGeneratorValue(region, gen_VibratoLfoToPitch)
}

func GetregionPairModulationEnvelopeToPitch(region regionPair) int32 {
	return getGeneratorValue(region, gen_ModulationEnvelopeToPitch)
}

func GetregionPairInitialFilterCutoffFrequency(region regionPair) float32 {
	return calcCentsToHertz(float32(getGeneratorValue(region, gen_InitialFilterCutoffFrequency)))
}

func GetregionPairInitialFilterQ(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_InitialFilterQ))
}

func GetregionPairModulationLfoToFilterCutoffFrequency(region regionPair) int32 {
	return getGeneratorValue(region, gen_ModulationLfoToFilterCutoffFrequency)
}

func GetregionPairModulationEnvelopeToFilterCutoffFrequency(region regionPair) int32 {
	return getGeneratorValue(region, gen_ModulationEnvelopeToFilterCutoffFrequency)
}

func GetregionPairModulationLfoToVolume(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_ModulationLfoToVolume))
}

func GetregionPairChorusEffectsSend(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_ChorusEffectsSend))
}

func GetregionPairReverbEffectsSend(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_ReverbEffectsSend))
}

func GetregionPairPan(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_Pan))
}

func GetregionPairDelayModulationLfo(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_DelayModulationLfo)))
}

func GetregionPairFrequencyModulationLfo(region regionPair) float32 {
	return calcCentsToHertz(float32(getGeneratorValue(region, gen_FrequencyModulationLfo)))
}

func GetregionPairDelayVibratoLfo(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_DelayVibratoLfo)))
}

func GetregionPairFrequencyVibratoLfo(region regionPair) float32 {
	return calcCentsToHertz(float32(getGeneratorValue(region, gen_FrequencyVibratoLfo)))
}

func GetregionPairDelayModulationEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_DelayModulationEnvelope)))
}

func GetregionPairAttackModulationEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_AttackModulationEnvelope)))
}

func GetregionPairHoldModulationEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_HoldModulationEnvelope)))
}

func GetregionPairDecayModulationEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_DecayModulationEnvelope)))
}

func GetregionPairSustainModulationEnvelope(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_SustainModulationEnvelope))
}

func GetregionPairReleaseModulationEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_ReleaseModulationEnvelope)))
}

func GetregionPairKeyNumberToModulationEnvelopeHold(region regionPair) int32 {
	return getGeneratorValue(region, gen_KeyNumberToModulationEnvelopeHold)
}

func GetregionPairKeyNumberToModulationEnvelopeDecay(region regionPair) int32 {
	return getGeneratorValue(region, gen_KeyNumberToModulationEnvelopeDecay)
}

func GetregionPairDelayVolumeEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_DelayVolumeEnvelope)))
}

func GetregionPairAttackVolumeEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_AttackVolumeEnvelope)))
}

func GetregionPairHoldVolumeEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_HoldVolumeEnvelope)))
}

func GetregionPairDecayVolumeEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_DecayVolumeEnvelope)))
}

func GetregionPairSustainVolumeEnvelope(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_SustainVolumeEnvelope))
}

func GetregionPairReleaseVolumeEnvelope(region regionPair) float32 {
	return calcTimecentsToSeconds(float32(getGeneratorValue(region, gen_ReleaseVolumeEnvelope)))
}

func GetregionPairKeyNumberToVolumeEnvelopeHold(region regionPair) int32 {
	return getGeneratorValue(region, gen_KeyNumberToVolumeEnvelopeHold)
}

func GetregionPairKeyNumberToVolumeEnvelopeDecay(region regionPair) int32 {
	return getGeneratorValue(region, gen_KeyNumberToVolumeEnvelopeDecay)
}

func GetregionPairInitialAttenuation(region regionPair) float32 {
	return float32(0.1) * float32(getGeneratorValue(region, gen_InitialAttenuation))
}

func GetregionPairCoarseTune(region regionPair) int32 {
	return getGeneratorValue(region, gen_CoarseTune)
}

func GetregionPairFineTune(region regionPair) int32 {
	return getGeneratorValue(region, gen_FineTune) + int32(region.instrument.Sample.PitchCorrection)
}

func GetregionPairSampleModes(region regionPair) int32 {
	return GetInstrumentSampleModes(region.instrument)
}

func GetregionPairScaleTuning(region regionPair) int32 {
	return getGeneratorValue(region, gen_ScaleTuning)
}

func GetregionPairExclusiveClass(region regionPair) int32 {
	return GetInstrumentExclusiveClass(region.instrument)
}

func GetregionPairRootKey(region regionPair) int32 {
	return GetInstrumentRootKey(region.instrument)
}
