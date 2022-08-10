package meltysynth

import (
	"errors"
	"strconv"
)

type InstrumentRegion struct {
	Sample *SampleHeader
	gs     [61]int16
}

func createInstrumentRegion(instrument *Instrument, global []generator, local []generator, samples []*SampleHeader) (*InstrumentRegion, error) {

	result := new(InstrumentRegion)

	result.gs[GEN_InitialFilterCutoffFrequency] = 13500
	result.gs[GEN_DelayModulationLfo] = -12000
	result.gs[GEN_DelayVibratoLfo] = -12000
	result.gs[GEN_DelayModulationEnvelope] = -12000
	result.gs[GEN_AttackModulationEnvelope] = -12000
	result.gs[GEN_HoldModulationEnvelope] = -12000
	result.gs[GEN_DecayModulationEnvelope] = -12000
	result.gs[GEN_ReleaseModulationEnvelope] = -12000
	result.gs[GEN_DelayVolumeEnvelope] = -12000
	result.gs[GEN_AttackVolumeEnvelope] = -12000
	result.gs[GEN_HoldVolumeEnvelope] = -12000
	result.gs[GEN_DecayVolumeEnvelope] = -12000
	result.gs[GEN_ReleaseVolumeEnvelope] = -12000
	result.gs[GEN_KeyRange] = 0x7F00
	result.gs[GEN_VelocityRange] = 0x7F00
	result.gs[GEN_KeyNumber] = -1
	result.gs[GEN_Velocity] = -1
	result.gs[GEN_ScaleTuning] = 100
	result.gs[GEN_OverridingRootKey] = -1

	if global != nil {
		for i := 0; i < len(global); i++ {
			setInstrumentRegionParameter(result, global[i])
		}
	}

	if local != nil {
		for i := 0; i < len(local); i++ {
			setInstrumentRegionParameter(result, local[i])
		}
	}

	id := result.gs[GEN_SampleID]
	if !(0 <= id && int(id) < len(samples)) {
		return nil, errors.New("The instrument '" + instrument.Name + "' contains an invalid sample ID '" + strconv.Itoa(int(id)) + "'.")
	}
	result.Sample = samples[id]

	return result, nil
}

func createInstrumentRegions(instrument *Instrument, zones []*zone, samples []*SampleHeader) ([]*InstrumentRegion, error) {

	var global *zone = nil
	var err error

	// Is the first one the global zone?
	if len(zones[0].generators) == 0 || zones[0].generators[len(zones[0].generators)-1].generatorType != GEN_SampleID {
		// The first one is the global zone.
		global = zones[0]
	}

	if global != nil {
		count := len(zones) - 1
		regions := make([]*InstrumentRegion, count, count)
		for i := 0; i < count; i++ {
			regions[i], err = createInstrumentRegion(instrument, global.generators, zones[i+1].generators, samples)
			if err != nil {
				return nil, err
			}
		}
		return regions, nil
	} else {
		// No global zone.
		count := len(zones)
		regions := make([]*InstrumentRegion, count, count)
		for i := 0; i < count; i++ {
			regions[i], err = createInstrumentRegion(instrument, nil, zones[i].generators, samples)
			if err != nil {
				return nil, err
			}
		}
		return regions, nil
	}
}

func setInstrumentRegionParameter(region *InstrumentRegion, parameter generator) {

	index := parameter.generatorType

	// Unknown generators should be ignored.
	if 0 <= index && int(index) < len(region.gs) {
		region.gs[index] = int16(parameter.value)
	}
}

func GetInstrumentSampleStart(region *InstrumentRegion) int32 {
	return region.Sample.Start + GetInstrumentStartAddressOffset(region)
}

func GetInstrumentSampleEnd(region *InstrumentRegion) int32 {
	return region.Sample.End + GetInstrumentEndAddressOffset(region)
}

func GetInstrumentSampleStartLoop(region *InstrumentRegion) int32 {
	return region.Sample.StartLoop + GetInstrumentStartLoopAddressOffset(region)
}

func GetInstrumentSampleEndLoop(region *InstrumentRegion) int32 {
	return region.Sample.EndLoop + GetInstrumentEndLoopAddressOffset(region)
}

func GetInstrumentStartAddressOffset(region *InstrumentRegion) int32 {
	return 32768*int32(region.gs[GEN_StartAddressCoarseOffset]) + int32(region.gs[GEN_StartAddressOffset])
}

func GetInstrumentEndAddressOffset(region *InstrumentRegion) int32 {
	return 32768*int32(region.gs[GEN_EndAddressCoarseOffset]) + int32(region.gs[GEN_EndAddressOffset])
}

func GetInstrumentStartLoopAddressOffset(region *InstrumentRegion) int32 {
	return 32768*int32(region.gs[GEN_StartLoopAddressCoarseOffset]) + int32(region.gs[GEN_StartLoopAddressOffset])
}

func GetInstrumentEndLoopAddressOffset(region *InstrumentRegion) int32 {
	return 32768*int32(region.gs[GEN_EndLoopAddressCoarseOffset]) + int32(region.gs[GEN_EndLoopAddressOffset])
}

func GetInstrumentModulationLfoToPitch(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_ModulationLfoToPitch])
}

func GetInstrumentVibratoLfoToPitch(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_VibratoLfoToPitch])
}

func GetInstrumentModulationEnvelopeToPitch(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_ModulationEnvelopeToPitch])
}

func GetInstrumentInitialFilterCutoffFrequency(region *InstrumentRegion) float32 {
	return calcCentsToHertz(float32(region.gs[GEN_InitialFilterCutoffFrequency]))
}

func GetInstrumentInitialFilterQ(region *InstrumentRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_InitialFilterQ])
}

func GetInstrumentModulationLfoToFilterCutoffFrequency(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_ModulationLfoToFilterCutoffFrequency])
}

func GetInstrumentModulationEnvelopeToFilterCutoffFrequency(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_ModulationEnvelopeToFilterCutoffFrequency])
}

func GetInstrumentModulationLfoToVolume(region *InstrumentRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_ModulationLfoToVolume])
}

func GetInstrumentChorusEffectsSend(region *InstrumentRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_ChorusEffectsSend])
}

func GetInstrumentReverbEffectsSend(region *InstrumentRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_ReverbEffectsSend])
}

func GetInstrumentPan(region *InstrumentRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_Pan])
}

func GetInstrumentDelayModulationLfo(region *InstrumentRegion) float32 {
	return calcTimecentsToSeconds(float32(region.gs[GEN_DelayModulationLfo]))
}

func GetInstrumentFrequencyModulationLfo(region *InstrumentRegion) float32 {
	return calcCentsToHertz(float32(region.gs[GEN_FrequencyModulationLfo]))
}

func GetInstrumentDelayVibratoLfo(region *InstrumentRegion) float32 {
	return calcTimecentsToSeconds(float32(region.gs[GEN_DelayVibratoLfo]))
}

func GetInstrumentFrequencyVibratoLfo(region *InstrumentRegion) float32 {
	return calcCentsToHertz(float32(region.gs[GEN_FrequencyVibratoLfo]))
}

func GetInstrumentDelayModulationEnvelope(region *InstrumentRegion) float32 {
	return calcTimecentsToSeconds(float32(region.gs[GEN_DelayModulationEnvelope]))
}

func GetInstrumentAttackModulationEnvelope(region *InstrumentRegion) float32 {
	return calcTimecentsToSeconds(float32(region.gs[GEN_AttackModulationEnvelope]))
}

func GetInstrumentHoldModulationEnvelope(region *InstrumentRegion) float32 {
	return calcTimecentsToSeconds(float32(region.gs[GEN_HoldModulationEnvelope]))
}

func GetInstrumentDecayModulationEnvelope(region *InstrumentRegion) float32 {
	return calcTimecentsToSeconds(float32(region.gs[GEN_DecayModulationEnvelope]))
}

func GetInstrumentSustainModulationEnvelope(region *InstrumentRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_SustainModulationEnvelope])
}

func GetInstrumentReleaseModulationEnvelope(region *InstrumentRegion) float32 {
	return calcTimecentsToSeconds(float32(region.gs[GEN_ReleaseModulationEnvelope]))
}

func GetInstrumentKeyNumberToModulationEnvelopeHold(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_KeyNumberToModulationEnvelopeHold])
}

func GetInstrumentKeyNumberToModulationEnvelopeDecay(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_KeyNumberToModulationEnvelopeDecay])
}

func GetInstrumentDelayVolumeEnvelope(region *InstrumentRegion) float32 {
	return calcTimecentsToSeconds(float32(region.gs[GEN_DelayVolumeEnvelope]))
}

func GetInstrumentAttackVolumeEnvelope(region *InstrumentRegion) float32 {
	return calcTimecentsToSeconds(float32(region.gs[GEN_AttackVolumeEnvelope]))
}

func GetInstrumentHoldVolumeEnvelope(region *InstrumentRegion) float32 {
	return calcTimecentsToSeconds(float32(region.gs[GEN_HoldVolumeEnvelope]))
}

func GetInstrumentDecayVolumeEnvelope(region *InstrumentRegion) float32 {
	return calcTimecentsToSeconds(float32(region.gs[GEN_DecayVolumeEnvelope]))
}

func GetInstrumentSustainVolumeEnvelope(region *InstrumentRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_SustainVolumeEnvelope])
}

func GetInstrumentReleaseVolumeEnvelope(region *InstrumentRegion) float32 {
	return calcTimecentsToSeconds(float32(region.gs[GEN_ReleaseVolumeEnvelope]))
}

func GetInstrumentKeyNumberToVolumeEnvelopeHold(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_KeyNumberToVolumeEnvelopeHold])
}

func GetInstrumentKeyNumberToVolumeEnvelopeDecay(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_KeyNumberToVolumeEnvelopeDecay])
}

func GetInstrumentKeyRangeStart(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_KeyRange]) & 0xFF
}

func GetInstrumentKeyRangeEnd(region *InstrumentRegion) int32 {
	return (int32(region.gs[GEN_KeyRange]) >> 8) & 0xFF
}

func GetInstrumentVelocityRangeStart(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_VelocityRange]) & 0xFF
}

func GetInstrumentVelocityRangeEnd(region *InstrumentRegion) int32 {
	return (int32(region.gs[GEN_VelocityRange]) >> 8) & 0xFF
}

func GetInstrumentInitialAttenuation(region *InstrumentRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_InitialAttenuation])
}

func GetInstrumentCoarseTune(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_CoarseTune])
}

func GetInstrumentFineTune(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_FineTune]) + int32(region.Sample.PitchCorrection)
}

func GetInstrumentSampleModes(region *InstrumentRegion) int32 {
	if region.gs[GEN_SampleModes] != 2 {
		return int32(region.gs[GEN_SampleModes])
	} else {
		return 0
	}
}

func GetInstrumentScaleTuning(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_ScaleTuning])
}

func GetInstrumentExclusiveClass(region *InstrumentRegion) int32 {
	return int32(region.gs[GEN_ExclusiveClass])
}

func GetInstrumentRootKey(region *InstrumentRegion) int32 {
	if region.gs[GEN_OverridingRootKey] != -1 {
		return int32(region.gs[GEN_OverridingRootKey])
	} else {
		return int32(region.Sample.OriginalPitch)
	}
}
