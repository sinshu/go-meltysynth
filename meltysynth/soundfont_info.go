package meltysynth

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type SoundFontInfo struct {
	Version           SoundFontVersion
	TargetSoundEngine string
	BankName          string
	RomName           string
	RomVersion        SoundFontVersion
	CreationDate      string
	Auther            string
	TargetProduct     string
	Copyright         string
	Comments          string
	Tools             string
}

func newSoundFontInfo(r io.Reader) (*SoundFontInfo, error) {
	var err error

	var chunkId string
	chunkId, err = readFourCC(r)
	if err != nil {
		return nil, err
	}
	if chunkId != "LIST" {
		return nil, errors.New("the list chunk was not found")
	}

	var pos int32 = 0
	var end int32
	err = binary.Read(r, binary.LittleEndian, &end)
	if err != nil {
		return nil, err
	}

	var listType string
	listType, err = readFourCC(r)
	if err != nil {
		return nil, err
	}
	if listType != "INFO" {
		return nil, fmt.Errorf(`the type of the list chunk must be "INFO", but was %q`, listType)
	}
	pos += 4
	result := new(SoundFontInfo)

	for pos < end {
		var id string
		id, err = readFourCC(r)
		if err != nil {
			return nil, err
		}
		pos += 4

		var size int32
		err = binary.Read(r, binary.LittleEndian, &size)
		if err != nil {
			return nil, err
		}
		pos += 4

		switch id {
		case "ifil":
			result.Version, err = newSoundFontVersion(r)
		case "isng":
			result.TargetSoundEngine, err = readFixedLengthString(r, size)
		case "INAM":
			result.BankName, err = readFixedLengthString(r, size)
		case "irom":
			result.RomName, err = readFixedLengthString(r, size)
		case "iver":
			result.RomVersion, err = newSoundFontVersion(r)
		case "ICRD":
			result.CreationDate, err = readFixedLengthString(r, size)
		case "IENG":
			result.Auther, err = readFixedLengthString(r, size)
		case "IPRD":
			result.TargetProduct, err = readFixedLengthString(r, size)
		case "ICOP":
			result.Copyright, err = readFixedLengthString(r, size)
		case "ICMT":
			result.Comments, err = readFixedLengthString(r, size)
		case "ISFT":
			result.Tools, err = readFixedLengthString(r, size)
		default:
			return nil, fmt.Errorf("the info list contains an unknown id %q", id)
		}

		if err != nil {
			return nil, err
		}

		pos += size
	}

	return result, nil
}
