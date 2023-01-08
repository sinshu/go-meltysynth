package meltysynth

import (
	"encoding/binary"
	"errors"
	"io"
)

type generator struct {
	generatorType uint16
	value         uint16
}

func readGeneratorsFromChunk(r io.Reader, size int32) ([]generator, error) {
	var n int
	var err error

	if size%4 != 0 {
		return nil, errors.New("the generator list is invalid")
	}

	count := size/4 - 1
	generators := make([]generator, count)

	for i := int32(0); i < count; i++ {
		var gen generator

		var generatorType uint16
		err = binary.Read(r, binary.LittleEndian, &generatorType)
		if err != nil {
			return nil, err
		}
		gen.generatorType = generatorType

		var value uint16
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return nil, err
		}
		gen.value = value

		generators[i] = gen
	}

	// The last one is the terminator.
	n, err = r.Read(make([]byte, 4))
	if err != nil {
		return nil, err
	}
	if n != 4 {
		return nil, errors.New("failed to read the generator list")
	}

	return generators, nil
}
