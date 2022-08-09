package meltysynth

import (
	"errors"
	"strconv"
)

type InstrumentRegion struct {
	Sample *SampleHeader
	gs     [61]int16
}

func createInstrumentRegion(instrument *Instrument, global []generator, local []generator, samples []SampleHeader) (*InstrumentRegion, error) {

	var result InstrumentRegion

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
			setInstrumentRegionParameter(&result, global[i])
		}
	}

	if local != nil {
		for i := 0; i < len(local); i++ {
			setInstrumentRegionParameter(&result, local[i])
		}
	}

	id := result.gs[GEN_Instrument]
	if !(0 <= id && int(id) < len(samples)) {
		return nil, errors.New("The instrument '" + instrument.Name + "' contains an invalid sample ID '" + strconv.Itoa(int(id)) + "'.")
	}

	return &result, nil
}

func setInstrumentRegionParameter(region *InstrumentRegion, parameter generator) {

	index := parameter.value

	// Unknown generators should be ignored.
	if 0 <= index && int(index) < len(region.gs) {
		region.gs[index] = int16(parameter.value)
	}
}
