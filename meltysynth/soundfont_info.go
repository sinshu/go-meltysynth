package meltysynth

import (
	"encoding/binary"
	"errors"
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

func NewSoundFontInfo(reader io.Reader) (*SoundFontInfo, error) {

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
	if listType != "INFO" {
		return nil, errors.New("The type of the LIST chunk must be 'INFO', but was '" + listType + "'.")
	}
	pos += 4

	result := new(SoundFontInfo)

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

		case "ifil":
			result.Version, err = NewSoundFontVersion(reader)
		case "isng":
			result.TargetSoundEngine, err = readFixedLengthString(reader, size)
		case "INAM":
			result.BankName, err = readFixedLengthString(reader, size)
		case "irom":
			result.RomName, err = readFixedLengthString(reader, size)
		case "iver":
			result.RomVersion, err = NewSoundFontVersion(reader)
		case "ICRD":
			result.CreationDate, err = readFixedLengthString(reader, size)
		case "IENG":
			result.Auther, err = readFixedLengthString(reader, size)
		case "IPRD":
			result.TargetProduct, err = readFixedLengthString(reader, size)
		case "ICOP":
			result.Copyright, err = readFixedLengthString(reader, size)
		case "ICMT":
			result.Comments, err = readFixedLengthString(reader, size)
		case "ISFT":
			result.Tools, err = readFixedLengthString(reader, size)
		default:
			return nil, errors.New("The INFO list contains an unknown ID '" + id + "'.")
		}

		if err != nil {
			return nil, err
		}

		pos += size
	}

	return result, nil
}
