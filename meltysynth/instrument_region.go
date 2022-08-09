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

	id := result.gs[GEN_Instrument]
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

	index := parameter.value

	// Unknown generators should be ignored.
	if 0 <= index && int(index) < len(region.gs) {
		region.gs[index] = int16(parameter.value)
	}
}
