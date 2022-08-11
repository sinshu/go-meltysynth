package meltysynth

import (
	"encoding/binary"
	"errors"
	"io"
)

type SoundFont struct {
	Info          *SoundFontInfo
	BitsPerSample int32
	WaveData      []int16
	SampleHeaders []*SampleHeader
	Presets       []*Preset
	Instruments   []*Instrument
}

func NewSoundFont(reader io.Reader) (*SoundFont, error) {

	var err error

	var chunkId string
	chunkId, err = readFourCC(reader)
	if err != nil {
		return nil, err
	}
	if chunkId != "RIFF" {
		return nil, errors.New("the riff chunk was not found")
	}

	var size int32
	err = binary.Read(reader, binary.LittleEndian, &size)
	if err != nil {
		return nil, err
	}

	var formType string
	formType, err = readFourCC(reader)
	if err != nil {
		return nil, err
	}
	if formType != "sfbk" {
		return nil, errors.New("the type of the riff chunk must be 'sfbk', but was '" + formType + "'")
	}

	result := new(SoundFont)

	result.Info, err = newSoundFontInfo(reader)
	if err != nil {
		return nil, err
	}

	var sampleData *soundFontSampleData
	sampleData, err = newSoundFontSampleData(reader)
	if err != nil {
		return nil, err
	}
	result.BitsPerSample = sampleData.bitsPerSample
	result.WaveData = sampleData.samples

	var parameters *soundFontParameters
	parameters, err = newSoundFontParameters(reader)
	if err != nil {
		return nil, err
	}

	result.SampleHeaders = parameters.sampleHeaders
	result.Presets = parameters.presets
	result.Instruments = parameters.instruments

	return result, nil
}
