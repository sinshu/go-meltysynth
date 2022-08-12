package meltysynth

import (
	"encoding/binary"
	"errors"
	"io"
)

func readFourCC(reader io.Reader) (string, error) {

	var data [4]byte
	n, err := reader.Read(data[:])
	if err != nil {
		return "", err
	}
	if n != 4 {
		return "", errors.New("failed to read the four-cc")
	}

	for i := 0; i < 4; i++ {
		value := data[i]
		if !(32 <= value && value <= 126) {
			data[i] = '?'
		}
	}

	return string(data[:]), nil
}

func readFixedLengthString(reader io.Reader, length int32) (string, error) {

	var data []byte = make([]byte, length)
	n, err := reader.Read(data[:])
	if err != nil {
		return "", err
	}
	if n != int(length) {
		return "", errors.New("failed to read the fixed-length string")
	}

	var actualLength int32
	for actualLength = 0; actualLength < length; actualLength++ {
		if data[actualLength] == 0 {
			break
		}
	}

	return string(data[0:actualLength]), nil
}

func readIntVariableLength(reader io.Reader) (int32, error) {

	acc := int32(0)
	count := 0
	for {
		var value byte
		err := binary.Read(reader, binary.LittleEndian, &value)
		if err != nil {
			return 0, err
		}
		acc = (acc << 7) | (int32(value) & 127)
		if (value & 128) == 0 {
			break
		}
		count++
		if count == 4 {
			return 0, errors.New("the length of the value must be equal to or less than 4")
		}
	}
	return acc, nil
}
