package meltysynth

import (
	"encoding/binary"
	"errors"
	"io"
)

type soundFontParameters struct {
	sampleHeaders []*SampleHeader
	presets       []*Preset
	instruments   []*Instrument
}

func newSoundFontParameters(reader io.Reader) (*soundFontParameters, error) {

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
	if listType != "pdta" {
		return nil, errors.New("the type of the list chunk must be 'pdta', but was '" + listType + "'")
	}
	pos += 4

	var presetInfos []*presetInfo
	var presetBag []*zoneInfo
	var presetGenerators []generator
	var instrumentInfos []*instrumentInfo
	var instrumentBag []*zoneInfo
	var instrumentGenerators []generator
	var sampleHeaders []*SampleHeader

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

		case "phdr":
			presetInfos, err = readPresetsFromChunk(reader, size)
		case "pbag":
			presetBag, err = readZonesFromChunk(reader, size)
		case "pmod":
			err = discardModulatorData(reader, size)
		case "pgen":
			presetGenerators, err = readGeneratorsFromChunk(reader, size)
		case "inst":
			instrumentInfos, err = readInstrumentsFromChunk(reader, size)
		case "ibag":
			instrumentBag, err = readZonesFromChunk(reader, size)
		case "imod":
			err = discardModulatorData(reader, size)
		case "igen":
			instrumentGenerators, err = readGeneratorsFromChunk(reader, size)
		case "shdr":
			sampleHeaders, err = readSampleHeadersFromChunk(reader, size)
		default:
			return nil, errors.New("the info list contains an unknown id '" + id + "'")
		}

		if err != nil {
			return nil, err
		}

		pos += size
	}

	if presetInfos == nil {
		return nil, errors.New("the phdr sub-chunk was not found")
	}
	if presetBag == nil {
		return nil, errors.New("the pbag sub-chunk was not found")
	}
	if presetGenerators == nil {
		return nil, errors.New("the pgen sub-chunk was not found")
	}
	if instrumentInfos == nil {
		return nil, errors.New("the inst sub-chunk was not found")
	}
	if instrumentBag == nil {
		return nil, errors.New("the ibag sub-chunk was not found")
	}
	if instrumentGenerators == nil {
		return nil, errors.New("the igen sub-chunk was not found")
	}
	if sampleHeaders == nil {
		return nil, errors.New("the shdr sub-chunk was not found")
	}

	parameters := new(soundFontParameters)

	parameters.sampleHeaders = sampleHeaders

	instrumentZones, err := createZones(instrumentBag, instrumentGenerators)
	if err != nil {
		return nil, err
	}

	parameters.instruments, err = createInstruments(instrumentInfos, instrumentZones, sampleHeaders)
	if err != nil {
		return nil, err
	}

	presetZones, err := createZones(presetBag, presetGenerators)
	if err != nil {
		return nil, err
	}

	parameters.presets, err = createPresets(presetInfos, presetZones, parameters.instruments)
	if err != nil {
		return nil, err
	}

	return parameters, nil
}
