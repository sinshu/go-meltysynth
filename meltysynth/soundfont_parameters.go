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
	if listType != "pdta" {
		return nil, errors.New("The type of the LIST chunk must be 'pdta', but was '" + listType + "'.")
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
			return nil, errors.New("The INFO list contains an unknown ID '" + id + "'.")
		}

		if err != nil {
			return nil, err
		}

		pos += size
	}

	if presetInfos == nil {
		return nil, errors.New("The PHDR sub-chunk was not found.")
	}
	if presetBag == nil {
		return nil, errors.New("The PBAG sub-chunk was not found.")
	}
	if presetGenerators == nil {
		return nil, errors.New("The PGEN sub-chunk was not found.")
	}
	if instrumentInfos == nil {
		return nil, errors.New("The INST sub-chunk was not found.")
	}
	if instrumentBag == nil {
		return nil, errors.New("The IBAG sub-chunk was not found.")
	}
	if instrumentGenerators == nil {
		return nil, errors.New("The IGEN sub-chunk was not found.")
	}
	if sampleHeaders == nil {
		return nil, errors.New("The SHDR sub-chunk was not found.")
	}

	parameters := new(soundFontParameters)

	parameters.sampleHeaders = sampleHeaders

	instrumentZones, err := createZones(instrumentBag, instrumentGenerators)
	parameters.instruments, err = createInstruments(instrumentInfos, instrumentZones, sampleHeaders)

	presetZones, err := createZones(presetBag, presetGenerators)
	parameters.presets, err = createPresets(presetInfos, presetZones, parameters.instruments)

	return parameters, nil
}
