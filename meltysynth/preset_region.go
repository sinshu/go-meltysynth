package meltysynth

import (
	"errors"
	"strconv"
)

type PresetRegion struct {
	Instrument *Instrument
	gs         [61]int16
}

func createPresetRegion(preset *Preset, global []generator, local []generator, instruments []*Instrument) (*PresetRegion, error) {

	result := new(PresetRegion)

	result.gs[gen_KeyRange] = 0x7F00
	result.gs[gen_VelocityRange] = 0x7F00

	if global != nil {
		for i := 0; i < len(global); i++ {
			setPresetRegionParameter(result, global[i])
		}
	}

	if local != nil {
		for i := 0; i < len(local); i++ {
			setPresetRegionParameter(result, local[i])
		}
	}

	id := result.gs[gen_Instrument]
	if !(0 <= id && int(id) < len(instruments)) {
		return nil, errors.New("The preset '" + preset.Name + "' contains an invalid instrument ID '" + strconv.Itoa(int(id)) + "'.")
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
	}

	if global != nil {
		count := len(zones) - 1
		regions := make([]*PresetRegion, count, count)
		for i := 0; i < count; i++ {
			regions[i], err = createPresetRegion(preset, global.generators, zones[i+1].generators, instruments)
			if err != nil {
				return nil, err
			}
		}
		return regions, nil
	} else {
		// No global zone.
		count := len(zones)
		regions := make([]*PresetRegion, count, count)
		for i := 0; i < count; i++ {
			regions[i], err = createPresetRegion(preset, nil, zones[i].generators, instruments)
		}
		return regions, nil
	}
}

func setPresetRegionParameter(region *PresetRegion, parameter generator) {

	index := parameter.generatorType

	// Unknown generators should be ignored.
	if 0 <= index && int(index) < len(region.gs) {
		region.gs[index] = int16(parameter.value)
	}
}

func GetPresetModulationLfoToPitch(region *PresetRegion) int32 {
	return int32(region.gs[gen_ModulationLfoToPitch])
}

func GetPresetVibratoLfoToPitch(region *PresetRegion) int32 {
	return int32(region.gs[gen_VibratoLfoToPitch])
}

func GetPresetModulationEnvelopeToPitch(region *PresetRegion) int32 {
	return int32(region.gs[gen_ModulationEnvelopeToPitch])
}

func GetPresetInitialFilterCutoffFrequency(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_InitialFilterCutoffFrequency]))
}

func GetPresetInitialFilterQ(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[gen_InitialFilterQ])
}

func GetPresetModulationLfoToFilterCutoffFrequency(region *PresetRegion) int32 {
	return int32(region.gs[gen_ModulationLfoToFilterCutoffFrequency])
}

func GetPresetModulationEnvelopeToFilterCutoffFrequency(region *PresetRegion) int32 {
	return int32(region.gs[gen_ModulationEnvelopeToFilterCutoffFrequency])
}

func GetPresetModulationLfoToVolume(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[gen_ModulationLfoToVolume])
}

func GetPresetChorusEffectsSend(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[gen_ChorusEffectsSend])
}

func GetPresetReverbEffectsSend(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[gen_ReverbEffectsSend])
}

func GetPresetPan(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[gen_Pan])
}

func GetPresetDelayModulationLfo(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_DelayModulationLfo]))
}

func GetPresetFrequencyModulationLfo(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_FrequencyModulationLfo]))
}

func GetPresetDelayVibratoLfo(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_DelayVibratoLfo]))
}

func GetPresetFrequencyVibratoLfo(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_FrequencyVibratoLfo]))
}

func GetPresetDelayModulationEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_DelayModulationEnvelope]))
}

func GetPresetAttackModulationEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_AttackModulationEnvelope]))
}

func GetPresetHoldModulationEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_HoldModulationEnvelope]))
}

func GetPresetDecayModulationEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_DecayModulationEnvelope]))
}

func GetPresetSustainModulationEnvelope(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[gen_SustainModulationEnvelope])
}

func GetPresetReleaseModulationEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_ReleaseModulationEnvelope]))
}

func GetPresetKeyNumberToModulationEnvelopeHold(region *PresetRegion) int32 {
	return int32(region.gs[gen_KeyNumberToModulationEnvelopeHold])
}

func GetPresetKeyNumberToModulationEnvelopeDecay(region *PresetRegion) int32 {
	return int32(region.gs[gen_KeyNumberToModulationEnvelopeDecay])
}

func GetPresetDelayVolumeEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_DelayVolumeEnvelope]))
}

func GetPresetAttackVolumeEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_AttackVolumeEnvelope]))
}

func GetPresetHoldVolumeEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_HoldVolumeEnvelope]))
}

func GetPresetDecayVolumeEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_DecayVolumeEnvelope]))
}

func GetPresetSustainVolumeEnvelope(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[gen_SustainVolumeEnvelope])
}

func GetPresetReleaseVolumeEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[gen_ReleaseVolumeEnvelope]))
}

func GetPresetKeyNumberToVolumeEnvelopeHold(region *PresetRegion) int32 {
	return int32(region.gs[gen_KeyNumberToVolumeEnvelopeHold])
}

func GetPresetKeyNumberToVolumeEnvelopeDecay(region *PresetRegion) int32 {
	return int32(region.gs[gen_KeyNumberToVolumeEnvelopeDecay])
}

func GetPresetKeyRangeStart(region *PresetRegion) int32 {
	return int32(region.gs[gen_KeyRange]) & 0xFF
}

func GetPresetKeyRangeEnd(region *PresetRegion) int32 {
	return (int32(region.gs[gen_KeyRange]) >> 8) & 0xFF
}

func GetPresetVelocityRangeStart(region *PresetRegion) int32 {
	return int32(region.gs[gen_VelocityRange]) & 0xFF
}

func GetPresetVelocityRangeEnd(region *PresetRegion) int32 {
	return (int32(region.gs[gen_VelocityRange]) >> 8) & 0xFF
}

func GetPresetInitialAttenuation(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[gen_InitialAttenuation])
}

func GetPresetCoarseTune(region *PresetRegion) int32 {
	return int32(region.gs[gen_CoarseTune])
}

func GetPresetFineTune(region *PresetRegion) int32 {
	return int32(region.gs[gen_FineTune])
}

func GetPresetScaleTuning(region *PresetRegion) int32 {
	return int32(region.gs[gen_ScaleTuning])
}
