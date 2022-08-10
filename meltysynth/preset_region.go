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
