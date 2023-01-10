package meltysynth

import (
	"fmt"
)

type InstrumentRegion struct {
	Sample *SampleHeader
	gs     [61]int16
}

func createInstrumentRegion(inst *Instrument, global *zone, local *zone, samples []*SampleHeader) (*InstrumentRegion, error) {
	result := new(InstrumentRegion)

	result.gs[gen_InitialFilterCutoffFrequency] = 13500
	result.gs[gen_DelayModulationLfo] = -12000
	result.gs[gen_DelayVibratoLfo] = -12000
	result.gs[gen_DelayModulationEnvelope] = -12000
	result.gs[gen_AttackModulationEnvelope] = -12000
	result.gs[gen_HoldModulationEnvelope] = -12000
	result.gs[gen_DecayModulationEnvelope] = -12000
	result.gs[gen_ReleaseModulationEnvelope] = -12000
	result.gs[gen_DelayVolumeEnvelope] = -12000
	result.gs[gen_AttackVolumeEnvelope] = -12000
	result.gs[gen_HoldVolumeEnvelope] = -12000
	result.gs[gen_DecayVolumeEnvelope] = -12000
	result.gs[gen_ReleaseVolumeEnvelope] = -12000
	result.gs[gen_KeyRange] = 0x7F00
	result.gs[gen_VelocityRange] = 0x7F00
	result.gs[gen_KeyNumber] = -1
	result.gs[gen_Velocity] = -1
	result.gs[gen_ScaleTuning] = 100
	result.gs[gen_OverridingRootKey] = -1

	for i := 0; i < len(global.generators); i++ {
		result.setParameter(global.generators[i])
	}

	for i := 0; i < len(local.generators); i++ {
		result.setParameter(local.generators[i])
	}

	id := result.gs[gen_SampleID]
	if !(0 <= id && int(id) < len(samples)) {
		return nil, fmt.Errorf("the instrument %q contains an invalid sample id %d", inst.Name, id)
	}
	result.Sample = samples[id]

	return result, nil
}

func createInstrumentRegions(inst *Instrument, zones []*zone, samples []*SampleHeader) ([]*InstrumentRegion, error) {
	var err error

	// Is the first one the global zone?
	if len(zones[0].generators) == 0 || zones[0].generators[len(zones[0].generators)-1].generatorType != gen_SampleID {

		// The first one is the global zone.
		global := zones[0]

		// The global zone is regarded as the base setting of subsequent zones.
		count := len(zones) - 1
		regions := make([]*InstrumentRegion, count)
		for i := 0; i < count; i++ {
			regions[i], err = createInstrumentRegion(inst, global, zones[i+1], samples)
			if err != nil {
				return nil, err
			}
		}
		return regions, nil
	} else {
		// No global zone.
		count := len(zones)
		regions := make([]*InstrumentRegion, count)
		for i := 0; i < count; i++ {
			regions[i], err = createInstrumentRegion(inst, createEmptyZone(), zones[i], samples)
			if err != nil {
				return nil, err
			}
		}
		return regions, nil
	}
}

func (region *InstrumentRegion) setParameter(param generator) {
	index := param.generatorType

	// Unknown generators should be ignored.
	if 0 <= int(index) && int(index) < len(region.gs) {
		region.gs[index] = int16(param.value)
	}
}

func (region *InstrumentRegion) contains(key int32, velocity int32) bool {
	containsKey := region.GetKeyRangeStart() <= key && key <= region.GetKeyRangeEnd()
	containsVelocity := region.GetVelocityRangeStart() <= velocity && velocity <= region.GetVelocityRangeEnd()
	return containsKey && containsVelocity
}

func (region *InstrumentRegion) GetSampleStart() int32 {
	return region.Sample.Start + region.GetStartAddressOffset()
}

func (region *InstrumentRegion) GetSampleEnd() int32 {
	return region.Sample.End + region.GetEndAddressOffset()
}

func (region *InstrumentRegion) GetSampleStartLoop() int32 {
	return region.Sample.StartLoop + region.GetStartLoopAddressOffset()
}

func (region *InstrumentRegion) GetSampleEndLoop() int32 {
	return region.Sample.EndLoop + region.GetEndLoopAddressOffset()
}

func (region *InstrumentRegion) GetStartAddressOffset() int32 {
	return 32768*int32(region.gs[gen_StartAddressCoarseOffset]) + int32(region.gs[gen_StartAddressOffset])
}

func (region *InstrumentRegion) GetEndAddressOffset() int32 {
	return 32768*int32(region.gs[gen_EndAddressCoarseOffset]) + int32(region.gs[gen_EndAddressOffset])
}

func (region *InstrumentRegion) GetStartLoopAddressOffset() int32 {
	return 32768*int32(region.gs[gen_StartLoopAddressCoarseOffset]) + int32(region.gs[gen_StartLoopAddressOffset])
}

func (region *InstrumentRegion) GetEndLoopAddressOffset() int32 {
	return 32768*int32(region.gs[gen_EndLoopAddressCoarseOffset]) + int32(region.gs[gen_EndLoopAddressOffset])
}

func (region *InstrumentRegion) GetModulationLfoToPitch() int32 {
	return int32(region.gs[gen_ModulationLfoToPitch])
}

func (region *InstrumentRegion) GetVibratoLfoToPitch() int32 {
	return int32(region.gs[gen_VibratoLfoToPitch])
}

func (region *InstrumentRegion) GetModulationEnvelopeToPitch() int32 {
	return int32(region.gs[gen_ModulationEnvelopeToPitch])
}

func (region *InstrumentRegion) GetInitialFilterCutoffFrequency() float32 {
	return calcCentsToHertz(float32(region.gs[gen_InitialFilterCutoffFrequency]))
}

func (region *InstrumentRegion) GetInitialFilterQ() float32 {
	return float32(0.1) * float32(region.gs[gen_InitialFilterQ])
}

func (region *InstrumentRegion) GetModulationLfoToFilterCutoffFrequency() int32 {
	return int32(region.gs[gen_ModulationLfoToFilterCutoffFrequency])
}

func (region *InstrumentRegion) GetModulationEnvelopeToFilterCutoffFrequency() int32 {
	return int32(region.gs[gen_ModulationEnvelopeToFilterCutoffFrequency])
}

func (region *InstrumentRegion) GetModulationLfoToVolume() float32 {
	return float32(0.1) * float32(region.gs[gen_ModulationLfoToVolume])
}

func (region *InstrumentRegion) GetChorusEffectsSend() float32 {
	return float32(0.1) * float32(region.gs[gen_ChorusEffectsSend])
}

func (region *InstrumentRegion) GetReverbEffectsSend() float32 {
	return float32(0.1) * float32(region.gs[gen_ReverbEffectsSend])
}

func (region *InstrumentRegion) GetPan() float32 {
	return float32(0.1) * float32(region.gs[gen_Pan])
}

func (region *InstrumentRegion) GetDelayModulationLfo() float32 {
	return calcTimecentsToSeconds(float32(region.gs[gen_DelayModulationLfo]))
}

func (region *InstrumentRegion) GetFrequencyModulationLfo() float32 {
	return calcCentsToHertz(float32(region.gs[gen_FrequencyModulationLfo]))
}

func (region *InstrumentRegion) GetDelayVibratoLfo() float32 {
	return calcTimecentsToSeconds(float32(region.gs[gen_DelayVibratoLfo]))
}

func (region *InstrumentRegion) GetFrequencyVibratoLfo() float32 {
	return calcCentsToHertz(float32(region.gs[gen_FrequencyVibratoLfo]))
}

func (region *InstrumentRegion) GetDelayModulationEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.gs[gen_DelayModulationEnvelope]))
}

func (region *InstrumentRegion) GetAttackModulationEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.gs[gen_AttackModulationEnvelope]))
}

func (region *InstrumentRegion) GetHoldModulationEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.gs[gen_HoldModulationEnvelope]))
}

func (region *InstrumentRegion) GetDecayModulationEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.gs[gen_DecayModulationEnvelope]))
}

func (region *InstrumentRegion) GetSustainModulationEnvelope() float32 {
	return float32(0.1) * float32(region.gs[gen_SustainModulationEnvelope])
}

func (region *InstrumentRegion) GetReleaseModulationEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.gs[gen_ReleaseModulationEnvelope]))
}

func (region *InstrumentRegion) GetKeyNumberToModulationEnvelopeHold() int32 {
	return int32(region.gs[gen_KeyNumberToModulationEnvelopeHold])
}

func (region *InstrumentRegion) GetKeyNumberToModulationEnvelopeDecay() int32 {
	return int32(region.gs[gen_KeyNumberToModulationEnvelopeDecay])
}

func (region *InstrumentRegion) GetDelayVolumeEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.gs[gen_DelayVolumeEnvelope]))
}

func (region *InstrumentRegion) GetAttackVolumeEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.gs[gen_AttackVolumeEnvelope]))
}

func (region *InstrumentRegion) GetHoldVolumeEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.gs[gen_HoldVolumeEnvelope]))
}

func (region *InstrumentRegion) GetDecayVolumeEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.gs[gen_DecayVolumeEnvelope]))
}

func (region *InstrumentRegion) GetSustainVolumeEnvelope() float32 {
	return float32(0.1) * float32(region.gs[gen_SustainVolumeEnvelope])
}

func (region *InstrumentRegion) GetReleaseVolumeEnvelope() float32 {
	return calcTimecentsToSeconds(float32(region.gs[gen_ReleaseVolumeEnvelope]))
}

func (region *InstrumentRegion) GetKeyNumberToVolumeEnvelopeHold() int32 {
	return int32(region.gs[gen_KeyNumberToVolumeEnvelopeHold])
}

func (region *InstrumentRegion) GetKeyNumberToVolumeEnvelopeDecay() int32 {
	return int32(region.gs[gen_KeyNumberToVolumeEnvelopeDecay])
}

func (region *InstrumentRegion) GetKeyRangeStart() int32 {
	return int32(region.gs[gen_KeyRange]) & 0xFF
}

func (region *InstrumentRegion) GetKeyRangeEnd() int32 {
	return (int32(region.gs[gen_KeyRange]) >> 8) & 0xFF
}

func (region *InstrumentRegion) GetVelocityRangeStart() int32 {
	return int32(region.gs[gen_VelocityRange]) & 0xFF
}

func (region *InstrumentRegion) GetVelocityRangeEnd() int32 {
	return (int32(region.gs[gen_VelocityRange]) >> 8) & 0xFF
}

func (region *InstrumentRegion) GetInitialAttenuation() float32 {
	return float32(0.1) * float32(region.gs[gen_InitialAttenuation])
}

func (region *InstrumentRegion) GetCoarseTune() int32 {
	return int32(region.gs[gen_CoarseTune])
}

func (region *InstrumentRegion) GetFineTune() int32 {
	return int32(region.gs[gen_FineTune]) + int32(region.Sample.PitchCorrection)
}

func (region *InstrumentRegion) GetSampleModes() int32 {
	if region.gs[gen_SampleModes] != 2 {
		return int32(region.gs[gen_SampleModes])
	} else {
		return loop_NoLoop
	}
}

func (region *InstrumentRegion) GetScaleTuning() int32 {
	return int32(region.gs[gen_ScaleTuning])
}

func (region *InstrumentRegion) GetExclusiveClass() int32 {
	return int32(region.gs[gen_ExclusiveClass])
}

func (region *InstrumentRegion) GetRootKey() int32 {
	if region.gs[gen_OverridingRootKey] != -1 {
		return int32(region.gs[gen_OverridingRootKey])
	} else {
		return int32(region.Sample.OriginalPitch)
	}
}
