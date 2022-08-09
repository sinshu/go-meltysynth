package meltysynth

import (
	"encoding/binary"
	"errors"
	"io"
)

type SampleHeader struct {
	name            string
	start           int32
	end             int32
	startLoop       int32
	endLoop         int32
	sampleRate      int32
	originalPitch   uint8
	pitchCorrection int8
	link            uint16
	sampleType      uint16
}

func readSampleHeadersFromChunk(reader io.Reader, size int32) ([]SampleHeader, error) {

	var n int
	var err error

	if size%46 != 0 {
		return nil, errors.New("The sample header list is invalid.")
	}

	count := size/46 - 1

	headers := make([]SampleHeader, count, count)

	for i := int32(0); i < count; i++ {

		var header SampleHeader

		header.name, err = readFixedLengthString(reader, 20)
		if err != nil {
			return nil, err
		}

		var start int32
		err = binary.Read(reader, binary.LittleEndian, &start)
		if err != nil {
			return nil, err
		}
		header.start = start

		var end int32
		err = binary.Read(reader, binary.LittleEndian, &end)
		if err != nil {
			return nil, err
		}
		header.end = end

		var startLoop int32
		err = binary.Read(reader, binary.LittleEndian, &startLoop)
		if err != nil {
			return nil, err
		}
		header.startLoop = startLoop

		var endLoop int32
		err = binary.Read(reader, binary.LittleEndian, &endLoop)
		if err != nil {
			return nil, err
		}
		header.endLoop = endLoop

		var sampleRate int32
		err = binary.Read(reader, binary.LittleEndian, &sampleRate)
		if err != nil {
			return nil, err
		}
		header.sampleRate = sampleRate

		var originalPitch uint8
		err = binary.Read(reader, binary.LittleEndian, &originalPitch)
		if err != nil {
			return nil, err
		}
		header.originalPitch = originalPitch

		var pitchCorrection int8
		err = binary.Read(reader, binary.LittleEndian, &pitchCorrection)
		if err != nil {
			return nil, err
		}
		header.pitchCorrection = pitchCorrection

		var link uint16
		err = binary.Read(reader, binary.LittleEndian, &link)
		if err != nil {
			return nil, err
		}
		header.link = link

		var sampleType uint16
		err = binary.Read(reader, binary.LittleEndian, &sampleType)
		if err != nil {
			return nil, err
		}
		header.sampleType = sampleType

		headers[i] = header
	}

	// The last one is the terminator.
	n, err = reader.Read(make([]byte, 46, 46))
	if err != nil {
		return nil, err
	}
	if n != 46 {
		return nil, errors.New("Failed to read the sample header list.")
	}

	return headers, nil
}
