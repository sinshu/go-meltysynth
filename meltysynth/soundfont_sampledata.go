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
		return nil, errors.New("The LIST chunk was not found.")
	}

	var pos int32 = 0
	var end int32
	err = binary.Read(reader, binary.LittleEndian, &end)
	if err != nil {
		return nil, err
	}

	var listType string
	listType, err = readFourCC(reader)
	if listType != "sdta" {
		return nil, errors.New("The type of the LIST chunk must be 'sdta', but was '" + listType + "'.")
	}
	pos += 4

	var result soundFontSampleData

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
			result.samples = make([]int16, size/2, size/2)
			err = binary.Read(reader, binary.LittleEndian, result.samples)

		case "sm24":
			// 24 bit audio is not supported.
			n, err = reader.Read(make([]byte, size))
			if n != int(size) {
				return nil, errors.New("Failed to read the 24 bit audio data.")
			}

		default:
			return nil, errors.New("The INFO list contains an unknown ID '" + id + "'.")
		}

		if err != nil {
			return nil, err
		}

		pos += size
	}

	if result.samples == nil {
		return nil, errors.New("No valid sample data was found.")
	}

	return &result, nil
}
