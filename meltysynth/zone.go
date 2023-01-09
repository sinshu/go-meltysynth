package meltysynth

import (
	"errors"
)

type zone struct {
	generators []generator
}

func createZones(infos []*zoneInfo, generators []generator) ([]*zone, error) {
	if len(infos) <= 1 {
		return nil, errors.New("no valid zone was found")
	}

	// The last one is the terminator.
	count := len(infos) - 1
	zones := make([]*zone, count)

	for i := 0; i < count; i++ {
		info := infos[i]

		zo := new(zone)
		zo.generators = make([]generator, info.generatorCount)
		for j := int32(0); j < info.generatorCount; j++ {
			zo.generators[j] = generators[info.generatorIndex+j]
		}

		zones[i] = zo
	}

	return zones, nil
}

func createEmptyZone() *zone {
	result := new(zone)
	result.generators = make([]generator, 0)
	return result
}
