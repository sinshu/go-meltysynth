package meltysynth

import (
	"encoding/binary"
	"errors"
	"io"
)

type zoneInfo struct {
	generatorIndex int32
	modulatorIndex int32
	generatorCount int32
	modulatorCount int32
}

func readZonesFromChunk(r io.Reader, size int32) ([]*zoneInfo, error) {
	var err error

	if size%4 != 0 {
		return nil, errors.New("the zone list is invalid")
	}

	count := size / 4
	zones := make([]*zoneInfo, count)

	for i := int32(0); i < count; i++ {
		zone := new(zoneInfo)

		var generatorIndex uint16
		err = binary.Read(r, binary.LittleEndian, &generatorIndex)
		if err != nil {
			return nil, err
		}
		zone.generatorIndex = int32(generatorIndex)

		var modulatorIndex uint16
		err = binary.Read(r, binary.LittleEndian, &modulatorIndex)
		if err != nil {
			return nil, err
		}
		zone.modulatorIndex = int32(modulatorIndex)

		zones[i] = zone
	}

	for i := int32(0); i < count-1; i++ {
		zones[i].generatorCount = zones[i+1].generatorIndex - zones[i].generatorIndex
		zones[i].modulatorCount = zones[i+1].modulatorIndex - zones[i].modulatorIndex
	}

	return zones, nil
}
