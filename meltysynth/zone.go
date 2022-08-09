package meltysynth

import (
	"errors"
)

type zone struct {
	generators []generator
}

func createZones(infos []zoneInfo, generators []generator) ([]zone, error) {

	if len(infos) <= 1 {
		return nil, errors.New("No valid zone was found.")
	}

	// The last one is the terminator.
	count := len(infos) - 1

	zones := make([]zone, count, count)

	for i := 0; i < count; i++ {

		info := infos[i]

		var zone zone
		zone.generators = make([]generator, info.generatorCount, info.generatorCount)
		for j := int32(0); j < info.generatorCount; j++ {
			zone.generators[j] = generators[info.generatorIndex+j]
		}

		zones[i] = zone
	}

	return zones, nil
}
