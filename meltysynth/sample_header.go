package meltysynth

import (
	"encoding/binary"
	"errors"
	"io"
)

type SampleHeader struct {
	Name            string
	Start           int32
	End             int32
	StartLoop       int32
	EndLoop         int32
	SampleRate      int32
	OriginalPitch   uint8
	PitchCorrection int8
	Link            uint16
	SampleType      uint16
}

func readSampleHeadersFromChunk(r io.Reader, size int32) ([]*SampleHeader, error) {
	var n int
	var err error

	if size%46 != 0 {
		return nil, errors.New("the sample header list is invalid")
	}

	count := size/46 - 1
	headers := make([]*SampleHeader, count)

	for i := int32(0); i < count; i++ {
		header := new(SampleHeader)

		header.Name, err = readFixedLengthString(r, 20)
		if err != nil {
			return nil, err
		}

		var start int32
		err = binary.Read(r, binary.LittleEndian, &start)
		if err != nil {
			return nil, err
		}
		header.Start = start

		var end int32
		err = binary.Read(r, binary.LittleEndian, &end)
		if err != nil {
			return nil, err
		}
		header.End = end

		var startLoop int32
		err = binary.Read(r, binary.LittleEndian, &startLoop)
		if err != nil {
			return nil, err
		}
		header.StartLoop = startLoop

		var endLoop int32
		err = binary.Read(r, binary.LittleEndian, &endLoop)
		if err != nil {
			return nil, err
		}
		header.EndLoop = endLoop

		var sampleRate int32
		err = binary.Read(r, binary.LittleEndian, &sampleRate)
		if err != nil {
			return nil, err
		}
		header.SampleRate = sampleRate

		var originalPitch uint8
		err = binary.Read(r, binary.LittleEndian, &originalPitch)
		if err != nil {
			return nil, err
		}
		header.OriginalPitch = originalPitch

		var pitchCorrection int8
		err = binary.Read(r, binary.LittleEndian, &pitchCorrection)
		if err != nil {
			return nil, err
		}
		header.PitchCorrection = pitchCorrection

		var link uint16
		err = binary.Read(r, binary.LittleEndian, &link)
		if err != nil {
			return nil, err
		}
		header.Link = link

		var sampleType uint16
		err = binary.Read(r, binary.LittleEndian, &sampleType)
		if err != nil {
			return nil, err
		}
		header.SampleType = sampleType

		headers[i] = header
	}

	// The last one is the terminator.
	n, err = r.Read(make([]byte, 46))
	if err != nil {
		return nil, err
	}
	if n != 46 {
		return nil, errors.New("failed to read the sample header list")
	}

	return headers, nil
}
