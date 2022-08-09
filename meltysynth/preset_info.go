package meltysynth

import (
	"encoding/binary"
	"errors"
	"io"
)

type presetInfo struct {
	name           string
	patchNumber    int32
	bankNumber     int32
	zoneStartIndex int32
	zoneEndIndex   int32
	library        int32
	genre          int32
	morphology     int32
}

func readPresetsFromChunk(reader io.Reader, size int32) ([]presetInfo, error) {

	var err error

	if size%38 != 0 {
		return nil, errors.New("The preset list is invalid.")
	}

	count := size / 38

	presets := make([]presetInfo, count, count)

	for i := int32(0); i < count; i++ {

		var preset presetInfo

		preset.name, err = readFixedLengthString(reader, 20)
		if err != nil {
			return nil, err
		}

		var patchNumber uint16
		err = binary.Read(reader, binary.LittleEndian, &patchNumber)
		if err != nil {
			return nil, err
		}
		preset.patchNumber = int32(patchNumber)

		var bankNumber uint16
		err = binary.Read(reader, binary.LittleEndian, &bankNumber)
		if err != nil {
			return nil, err
		}
		preset.bankNumber = int32(bankNumber)

		var zoneStartIndex uint16
		err = binary.Read(reader, binary.LittleEndian, &zoneStartIndex)
		if err != nil {
			return nil, err
		}
		preset.zoneStartIndex = int32(zoneStartIndex)

		var library int32
		err = binary.Read(reader, binary.LittleEndian, &library)
		if err != nil {
			return nil, err
		}
		preset.library = library

		var genre int32
		err = binary.Read(reader, binary.LittleEndian, &genre)
		if err != nil {
			return nil, err
		}
		preset.genre = genre

		var morphology int32
		err = binary.Read(reader, binary.LittleEndian, &morphology)
		if err != nil {
			return nil, err
		}
		preset.morphology = morphology

		presets[i] = preset
	}

	for i := int32(0); i < count-1; i++ {
		presets[i].zoneEndIndex = presets[i+1].zoneStartIndex - 1
	}

	return presets, nil
}
