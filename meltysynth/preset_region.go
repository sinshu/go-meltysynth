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

	result.gs[GEN_KeyRange] = 0x7F00
	result.gs[GEN_VelocityRange] = 0x7F00

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

	id := result.gs[GEN_Instrument]
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
	if len(zones[0].generators) == 0 || zones[0].generators[len(zones[0].generators)-1].generatorType != GEN_Instrument {
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
	return int32(region.gs[GEN_ModulationLfoToPitch])
}

func GetPresetVibratoLfoToPitch(region *PresetRegion) int32 {
	return int32(region.gs[GEN_VibratoLfoToPitch])
}

func GetPresetModulationEnvelopeToPitch(region *PresetRegion) int32 {
	return int32(region.gs[GEN_ModulationEnvelopeToPitch])
}

func GetPresetInitialFilterCutoffFrequency(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_InitialFilterCutoffFrequency]))
}

func GetPresetInitialFilterQ(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_InitialFilterQ])
}

func GetPresetModulationLfoToFilterCutoffFrequency(region *PresetRegion) int32 {
	return int32(region.gs[GEN_ModulationLfoToFilterCutoffFrequency])
}

func GetPresetModulationEnvelopeToFilterCutoffFrequency(region *PresetRegion) int32 {
	return int32(region.gs[GEN_ModulationEnvelopeToFilterCutoffFrequency])
}

func GetPresetModulationLfoToVolume(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_ModulationLfoToVolume])
}

func GetPresetChorusEffectsSend(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_ChorusEffectsSend])
}

func GetPresetReverbEffectsSend(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_ReverbEffectsSend])
}

func GetPresetPan(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_Pan])
}

func GetPresetDelayModulationLfo(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_DelayModulationLfo]))
}

func GetPresetFrequencyModulationLfo(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_FrequencyModulationLfo]))
}

func GetPresetDelayVibratoLfo(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_DelayVibratoLfo]))
}

func GetPresetFrequencyVibratoLfo(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_FrequencyVibratoLfo]))
}

func GetPresetDelayModulationEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_DelayModulationEnvelope]))
}

func GetPresetAttackModulationEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_AttackModulationEnvelope]))
}

func GetPresetHoldModulationEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_HoldModulationEnvelope]))
}

func GetPresetDecayModulationEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_DecayModulationEnvelope]))
}

func GetPresetSustainModulationEnvelope(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_SustainModulationEnvelope])
}

func GetPresetReleaseModulationEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_ReleaseModulationEnvelope]))
}

func GetPresetKeyNumberToModulationEnvelopeHold(region *PresetRegion) int32 {
	return int32(region.gs[GEN_KeyNumberToModulationEnvelopeHold])
}

func GetPresetKeyNumberToModulationEnvelopeDecay(region *PresetRegion) int32 {
	return int32(region.gs[GEN_KeyNumberToModulationEnvelopeDecay])
}

func GetPresetDelayVolumeEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_DelayVolumeEnvelope]))
}

func GetPresetAttackVolumeEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_AttackVolumeEnvelope]))
}

func GetPresetHoldVolumeEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_HoldVolumeEnvelope]))
}

func GetPresetDecayVolumeEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_DecayVolumeEnvelope]))
}

func GetPresetSustainVolumeEnvelope(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_SustainVolumeEnvelope])
}

func GetPresetReleaseVolumeEnvelope(region *PresetRegion) float32 {
	return calcCentsToMultiplyingFactor(float32(region.gs[GEN_ReleaseVolumeEnvelope]))
}

func GetPresetKeyNumberToVolumeEnvelopeHold(region *PresetRegion) int32 {
	return int32(region.gs[GEN_KeyNumberToVolumeEnvelopeHold])
}

func GetPresetKeyNumberToVolumeEnvelopeDecay(region *PresetRegion) int32 {
	return int32(region.gs[GEN_KeyNumberToVolumeEnvelopeDecay])
}

func GetPresetKeyRangeStart(region *PresetRegion) int32 {
	return int32(region.gs[GEN_KeyRange]) & 0xFF
}

func GetPresetKeyRangeEnd(region *PresetRegion) int32 {
	return (int32(region.gs[GEN_KeyRange]) >> 8) & 0xFF
}

func GetPresetVelocityRangeStart(region *PresetRegion) int32 {
	return int32(region.gs[GEN_VelocityRange]) & 0xFF
}

func GetPresetVelocityRangeEnd(region *PresetRegion) int32 {
	return (int32(region.gs[GEN_VelocityRange]) >> 8) & 0xFF
}

func GetPresetInitialAttenuation(region *PresetRegion) float32 {
	return float32(0.1) * float32(region.gs[GEN_InitialAttenuation])
}

func GetPresetCoarseTune(region *PresetRegion) int32 {
	return int32(region.gs[GEN_CoarseTune])
}

func GetPresetFineTune(region *PresetRegion) int32 {
	return int32(region.gs[GEN_FineTune])
}

func GetPresetScaleTuning(region *PresetRegion) int32 {
	return int32(region.gs[GEN_ScaleTuning])
}
