package meltysynth

import (
	"encoding/binary"
	"errors"
	"io"
)

type soundFontSampleData struct {
	bitsPerSample int32
	samples       []int16
}

func newSoundFontSampleData(reader io.Reader) (*soundFontSampleData, error) {

	var n int
	var err error

	var chunkId string
	chunkId, err = readFourCC(reader)
	if err != nil {
		return nil, err
	}
	if chunkId != "LIST" {
		return nil, errors.New("the list chunk was not found")
	}

	var pos int32 = 0
	var end int32
	err = binary.Read(reader, binary.LittleEndian, &end)
	if err != nil {
		return nil, err
	}

	var listType string
	listType, err = readFourCC(reader)
	if err != nil {
		return nil, err
	}
	if listType != "sdta" {
		return nil, errors.New("the type of the list chunk must be 'sdta', but was '" + listType + "'")
	}
	pos += 4

	result := new(soundFontSampleData)

	for pos < end {

		var id string
		id, err = readFourCC(reader)
		if err != nil {
			return nil, err
		}
		pos += 4

		var size int32
		err = binary.Read(reader, binary.LittleEndian, &size)
		if err != nil {
			return nil, err
		}
		pos += 4

		switch id {

		case "smpl":
			result.bitsPerSample = 16
			result.samples = make([]int16, size/2)
			err = binary.Read(reader, binary.LittleEndian, result.samples)

		case "sm24":
			// 24 bit audio is not supported.
			n, err = reader.Read(make([]byte, size))
			if n != int(size) {
				return nil, errors.New("failed to read the 24 bit audio data")
			}

		default:
			return nil, errors.New("the info list contains an unknown id '" + id + "'")
		}

		if err != nil {
			return nil, err
		}

		pos += size
	}

	if result.samples == nil {
		return nil, errors.New("no valid sample data was found")
	}

	return result, nil
}
