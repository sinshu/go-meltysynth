package meltysynth

import (
	"errors"
	"fmt"
)

type Preset struct {
	Name        string
	PatchNumber int32
	BankNumber  int32
	Library     int32
	Genre       int32
	Morphology  int32
	Regions     []*PresetRegion
}

func createPreset(info *presetInfo, zones []*zone, instruments []*Instrument) (*Preset, error) {
	var err error

	result := new(Preset)

	result.Name = info.name
	result.PatchNumber = info.patchNumber
	result.BankNumber = info.bankNumber
	result.Library = info.library
	result.Genre = info.genre
	result.Morphology = info.morphology

	zoneCount := info.zoneEndIndex - info.zoneStartIndex + 1
	if zoneCount <= 0 {
		return nil, fmt.Errorf("the preset %q has no zone", info.name)
	}

	zoneSpan := zones[info.zoneStartIndex : info.zoneStartIndex+zoneCount]

	result.Regions, err = createPresetRegions(result, zoneSpan, instruments)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func createPresets(infos []*presetInfo, zones []*zone, instruments []*Instrument) ([]*Preset, error) {

	var err error

	if len(infos) <= 1 {
		return nil, errors.New("no valid preset was found")
	}

	// The last one is the terminator.
	count := len(infos) - 1

	presets := make([]*Preset, count)

	for i := 0; i < count; i++ {
		presets[i], err = createPreset(infos[i], zones, instruments)
		if err != nil {
			return nil, err
		}
	}

	return presets, nil
}
