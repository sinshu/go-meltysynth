package meltysynth

import (
	"encoding/binary"
	"io"
)

type SoundFontVersion struct {
	Major int16
	Minor int16
}

func newSoundFontVersion(reader io.Reader) (SoundFontVersion, error) {

	var result SoundFontVersion
	var err error

	var major int16
	err = binary.Read(reader, binary.LittleEndian, &major)
	if err != nil {
		return result, err
	}

	var minor int16
	err = binary.Read(reader, binary.LittleEndian, &minor)
	if err != nil {
		return result, err
	}

	result.Major = major
	result.Minor = minor

	return result, nil
}
