package meltysynth

import (
	"encoding/binary"
	"errors"
	"io"
)

func readFourCC(r io.Reader) (string, error) {
	var data [4]byte
	n, err := r.Read(data[:])
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

func readFixedLengthString(r io.Reader, length int32) (string, error) {
	var data []byte = make([]byte, length)
	n, err := r.Read(data[:])
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

func readIntVariableLength(r io.Reader) (int32, error) {
	var acc int32
	count := 0
	for {
		var value byte
		err := binary.Read(r, binary.LittleEndian, &value)
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
