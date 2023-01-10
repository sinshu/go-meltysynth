package meltysynth

import (
	"fmt"
)

type PresetRegion struct {
	Instrument *Instrument
	gs         [61]int16
}

func createPresetRegion(preset *Preset, global *zone, local *zone, instruments []*Instrument) (*PresetRegion, error) {
	result := new(PresetRegion)

	result.gs[gen_KeyRange] = 0x7F00
	result.gs[gen_VelocityRange] = 0x7F00

	for i := 0; i < len(global.generators); i++ {
		result.setParameter(global.generators[i])
	}

	for i := 0; i < len(local.generators); i++ {
		result.setParameter(local.generators[i])
	}

	id := result.gs[gen_Instrument]
	if !(0 <= id && int(id) < len(instruments)) {
		return nil, fmt.Errorf("the preset %q contains an invalid instrument id %d", preset.Name, id)
	}
	result.Instrument = instruments[id]

	return result, nil
}

func createPresetRegions(preset *Preset, zones []*zone, instruments []*Instrument) ([]*PresetRegion, error) {

	var global *zone = nil
	var err error

	// Is the first one the global zone?
	if len(zones[0].generators) == 0 || zones[0].generators[len(zones[0].generators)-1].generatorType != gen_Instrument {

		// The first one is the global zone.
		global = zones[0]

		// The global zone is regarded as the base setting of subsequent zones.
		count := len(zones) - 1
		regions := make([]*PresetRegion, count)
		for i := 0; i < count; i++ {
			regions[i], err = createPresetRegion(preset, global, zones[i+1], instruments)
			if err != nil {
				return nil, err
			}
		}
		return regions, nil

	} else {

		// No global zone.
		count := len(zones)
		regions := make([]*PresetRegion, count)
		for i := 0; i < count; i++ {
			regions[i], err = createPresetRegion(preset, createEmptyZone(), zones[i], instruments)
			if err != nil {
				return nil, err
			}
		}
		return regions, nil
	}
}

func (region *PresetRegion) setParameter(param generator) {

	index := param.generatorType

	// Unknown generators should be ignored.
	if 0 <= int(index) && int(index) < len(region.gs) {
		region.gs[index] = int16(param.value)
	}
}

func (region *PresetRegion) contains(key int32, velocity int32) bool {
	containsKey := region.GetKeyRangeStart() <= key && key <= region.GetKeyRangeEnd()
	containsVelocity := region.GetVelocityRangeStart() <= velocity && velocity <= region.GetVelocityRangeEnd()
	return containsKey && containsVelocity
}

func (region *PresetRegion) GetModulationLfoToPitch() int32 {
	return int32(region.gs[gen_ModulationLfoToPitch])
}

func (region *PresetRegion) GetVibratoLfoToPitch() int32 {
	return int32(region.gs[gen_VibratoLfoToPitch])
}

func (region *PresetRegion) GetModulationEnvelopeToPitch() int32 {
	return int32(region.gs[gen_ModulationEnvelopeToPitch])
}

func (region *PresetRegion) GetInitialFilterCutoffFrequency() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_InitialFilterCutoffFrequency]))
}

func (region *PresetRegion) GetInitialFilterQ() float32 {
	return float32(0.1) * float32(region.gs[gen_InitialFilterQ])
}

func (region *PresetRegion) GetModulationLfoToFilterCutoffFrequency() int32 {
	return int32(region.gs[gen_ModulationLfoToFilterCutoffFrequency])
}

func (region *PresetRegion) GetModulationEnvelopeToFilterCutoffFrequency() int32 {
	return int32(region.gs[gen_ModulationEnvelopeToFilterCutoffFrequency])
}

func (region *PresetRegion) GetModulationLfoToVolume() float32 {
	return float32(0.1) * float32(region.gs[gen_ModulationLfoToVolume])
}

func (region *PresetRegion) GetChorusEffectsSend() float32 {
	return float32(0.1) * float32(region.gs[gen_ChorusEffectsSend])
}

func (region *PresetRegion) GetReverbEffectsSend() float32 {
	return float32(0.1) * float32(region.gs[gen_ReverbEffectsSend])
}

func (region *PresetRegion) GetPan() float32 {
	return float32(0.1) * float32(region.gs[gen_Pan])
}

func (region *PresetRegion) GetDelayModulationLfo() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_DelayModulationLfo]))
}

func (region *PresetRegion) GetFrequencyModulationLfo() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_FrequencyModulationLfo]))
}

func (region *PresetRegion) GetDelayVibratoLfo() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_DelayVibratoLfo]))
}

func (region *PresetRegion) GetFrequencyVibratoLfo() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_FrequencyVibratoLfo]))
}

func (region *PresetRegion) GetDelayModulationEnvelope() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_DelayModulationEnvelope]))
}

func (region *PresetRegion) GetAttackModulationEnvelope() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_AttackModulationEnvelope]))
}

func (region *PresetRegion) GetHoldModulationEnvelope() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_HoldModulationEnvelope]))
}

func (region *PresetRegion) GetDecayModulationEnvelope() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_DecayModulationEnvelope]))
}

func (region *PresetRegion) GetSustainModulationEnvelope() float32 {
	return float32(0.1) * float32(region.gs[gen_SustainModulationEnvelope])
}

func (region *PresetRegion) GetReleaseModulationEnvelope() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_ReleaseModulationEnvelope]))
}

func (region *PresetRegion) GetKeyNumberToModulationEnvelopeHold() int32 {
	return int32(region.gs[gen_KeyNumberToModulationEnvelopeHold])
}

func (region *PresetRegion) GetKeyNumberToModulationEnvelopeDecay() int32 {
	return int32(region.gs[gen_KeyNumberToModulationEnvelopeDecay])
}

func (region *PresetRegion) GetDelayVolumeEnvelope() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_DelayVolumeEnvelope]))
}

func (region *PresetRegion) GetAttackVolumeEnvelope() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_AttackVolumeEnvelope]))
}

func (region *PresetRegion) GetHoldVolumeEnvelope() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_HoldVolumeEnvelope]))
}

func (region *PresetRegion) GetDecayVolumeEnvelope() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_DecayVolumeEnvelope]))
}

func (region *PresetRegion) GetSustainVolumeEnvelope() float32 {
	return float32(0.1) * float32(region.gs[gen_SustainVolumeEnvelope])
}

func (region *PresetRegion) GetReleaseVolumeEnvelope() float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_ReleaseVolumeEnvelope]))
}

func (region *PresetRegion) GetKeyNumberToVolumeEnvelopeHold() int32 {
	return int32(region.gs[gen_KeyNumberToVolumeEnvelopeHold])
}

func (region *PresetRegion) GetKeyNumberToVolumeEnvelopeDecay() int32 {
	return int32(region.gs[gen_KeyNumberToVolumeEnvelopeDecay])
}

func (region *PresetRegion) GetKeyRangeStart() int32 {
	return int32(region.gs[gen_KeyRange]) & 0xFF
}

func (region *PresetRegion) GetKeyRangeEnd() int32 {
	return (int32(region.gs[gen_KeyRange]) >> 8) & 0xFF
}

func (region *PresetRegion) GetVelocityRangeStart() int32 {
	return int32(region.gs[gen_VelocityRange]) & 0xFF
}

func (region *PresetRegion) GetVelocityRangeEnd() int32 {
	return (int32(region.gs[gen_VelocityRange]) >> 8) & 0xFF
}

func (region *PresetRegion) GetInitialAttenuation() float32 {
	return float32(0.1) * float32(region.gs[gen_InitialAttenuation])
}

func (region *PresetRegion) GetCoarseTune() int32 {
	return int32(region.gs[gen_CoarseTune])
}

func (region *PresetRegion) GetFineTune() int32 {
	return int32(region.gs[gen_FineTune])
}

func (region *PresetRegion) GetScaleTuning() int32 {
	return int32(region.gs[gen_ScaleTuning])
}
