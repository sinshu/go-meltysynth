package meltysynth

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"time"
)

const (
	msg_Normal      byte = 0
	msg_TempoChange byte = 252
	msg_EndOfTrack  byte = 255
)

type message struct {
	channel byte
	command byte
	data1   byte
	data2   byte
}

type MidiFile struct {
	messages []message
	times    []time.Duration
}

func newMessage(channel byte, command byte, data1 byte, data2 byte) message {
	var result message
	result.channel = channel
	result.command = command
	result.data1 = data1
	result.data2 = data2
	return result
}

func common2b(status byte, data1 byte) message {
	channel := status & 0x0F
	command := status & 0xF0
	data2 := byte(0)
	return newMessage(channel, command, data1, data2)
}

func common3b(status byte, data1 byte, data2 byte) message {
	channel := status & 0x0F
	command := status & 0xF0
	return newMessage(channel, command, data1, data2)
}

func tempoChange(tempo int32) message {
	command := byte(tempo >> 16)
	data1 := byte(tempo >> 8)
	data2 := byte(tempo)
	return newMessage(msg_TempoChange, command, data1, data2)
}

func endOfTrack() message {
	return newMessage(msg_EndOfTrack, 0, 0, 0)
}

func (message message) getMessageType() byte {
	switch message.channel {
	case msg_TempoChange:
		return msg_TempoChange
	case msg_EndOfTrack:
		return msg_EndOfTrack
	default:
		return msg_Normal
	}
}

func (message message) getTempo() float64 {
	return 60000000.0 / float64((int32(message.command)<<16)|(int32(message.data1)<<8)|int32(message.data2))
}

func NewMidiFile(r io.Reader) (*MidiFile, error) {
	var err error

	chunkType, err := readFourCC(r)
	if err != nil {
		return nil, err
	}
	if chunkType != "MThd" {
		return nil, fmt.Errorf(`the chunk type must be "MThd", but was %q`, chunkType)
	}

	var size int32
	err = binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}
	if size != 6 {
		return nil, errors.New("the mthd chunk has invalid data")
	}

	var format int16
	err = binary.Read(r, binary.BigEndian, &format)
	if err != nil {
		return nil, err
	}
	if !(format == 0 || format == 1) {
		return nil, fmt.Errorf("the format %d is not supported", format)
	}

	var trackCount int16
	err = binary.Read(r, binary.BigEndian, &trackCount)
	if err != nil {
		return nil, err
	}

	var resolution int16
	err = binary.Read(r, binary.BigEndian, &resolution)
	if err != nil {
		return nil, err
	}

	messageLists := make([][]message, trackCount)
	tickLists := make([][]int32, trackCount)
	for i := int16(0); i < trackCount; i++ {
		messageList, tickList, err := readTrack(r)
		if err != nil {
			return nil, err
		}
		messageLists[i] = messageList
		tickLists[i] = tickList
	}

	messages, times := mergeTracks(messageLists, tickLists, resolution)

	result := new(MidiFile)
	result.messages = messages
	result.times = times

	return result, nil
}

func readTrack(r io.Reader) ([]message, []int32, error) {
	var n int
	var err error

	chunkType, err := readFourCC(r)
	if err != nil {
		return nil, nil, err
	}
	if chunkType != "MTrk" {
		return nil, nil, fmt.Errorf(`the chunk type must be "MTrk", but was %q`, chunkType)
	}

	r.Read(make([]byte, 4))

	messages := make([]message, 0, 300)
	ticks := make([]int32, 0, 300)

	var tick int32
	var lastStatus byte

	for {
		delta, err := readIntVariableLength(r)
		if err != nil {
			return nil, nil, err
		}

		var first byte
		err = binary.Read(r, binary.LittleEndian, &first)
		if err != nil {
			return nil, nil, err
		}

		tick += delta

		if (first & 128) == 0 {
			command := lastStatus & 0xF0
			if command == 0xC0 || command == 0xD0 {
				messages = append(messages, common2b(lastStatus, first))
				ticks = append(ticks, tick)
			} else {
				var data2 byte
				err = binary.Read(r, binary.LittleEndian, &data2)
				if err != nil {
					return nil, nil, err
				}
				messages = append(messages, common3b(lastStatus, first, data2))
				ticks = append(ticks, tick)
			}

			continue
		}

		switch first {
		case 0xF0: // System Exclusive
			err = discardData(r)
			if err != nil {
				return nil, nil, err
			}

		case 0xF7: // System Exclusive
			err = discardData(r)
			if err != nil {
				return nil, nil, err
			}

		case 0xFF: // Meta Event
			var metaEvent byte
			err = binary.Read(r, binary.LittleEndian, &metaEvent)
			if err != nil {
				return nil, nil, err
			}
			switch metaEvent {
			case 0x2F: // End of Track
				n, err = r.Read(make([]byte, 1))
				if err != nil {
					return nil, nil, err
				}
				if n != 1 {
					return nil, nil, err
				}
				messages = append(messages, endOfTrack())
				ticks = append(ticks, tick)
				return messages, ticks, nil

			case 0x51: // Tempo
				var tempo int32
				tempo, err = readTempo(r)
				if err != nil {
					return nil, nil, err
				}
				messages = append(messages, tempoChange(tempo))
				ticks = append(ticks, tick)

			default:
				err = discardData(r)
				if err != nil {
					return nil, nil, err
				}
			}

		default:
			command := first & 0xF0
			if command == 0xC0 || command == 0xD0 {
				var data1 byte
				err = binary.Read(r, binary.LittleEndian, &data1)
				if err != nil {
					return nil, nil, err
				}
				messages = append(messages, common2b(first, data1))
				ticks = append(ticks, tick)
			} else {
				var data1 byte
				err = binary.Read(r, binary.LittleEndian, &data1)
				if err != nil {
					return nil, nil, err
				}
				var data2 byte
				err = binary.Read(r, binary.LittleEndian, &data2)
				if err != nil {
					return nil, nil, err
				}
				messages = append(messages, common3b(first, data1, data2))
				ticks = append(ticks, tick)
			}
		}

		lastStatus = first
	}
}

func mergeTracks(messageLists [][]message, tickLists [][]int32, resolution int16) ([]message, []time.Duration) {
	mergedMessages := make([]message, 0, 1000)
	mergedTimes := make([]time.Duration, 0, 1000)

	indices := make([]int, len(messageLists))

	currentTick := int32(0)
	currentTime := time.Duration(0)

	tempo := float64(120)

	for {
		minTick := int32(math.MaxInt32)
		minIndex := int32(-1)
		tickListsLength := len(tickLists)
		for ch := 0; ch < tickListsLength; ch++ {
			if indices[ch] < len(tickLists[ch]) {
				var tick = tickLists[ch][indices[ch]]
				if tick < minTick {
					minTick = tick
					minIndex = int32(ch)
				}
			}
		}

		if minIndex == -1 {
			break
		}

		nextTick := tickLists[minIndex][indices[minIndex]]
		deltaTick := nextTick - currentTick
		deltaTime := time.Duration(float64(time.Second) * (60.0 / (float64(resolution) * tempo) * float64(deltaTick)))

		currentTick += deltaTick
		currentTime += deltaTime

		var message = messageLists[minIndex][indices[minIndex]]
		if message.getMessageType() == msg_TempoChange {
			tempo = message.getTempo()
		} else {
			mergedMessages = append(mergedMessages, message)
			mergedTimes = append(mergedTimes, currentTime)
		}

		indices[minIndex]++
	}

	return mergedMessages, mergedTimes
}

func readTempo(r io.Reader) (int32, error) {
	size, err := readIntVariableLength(r)
	if err != nil {
		return 0, err
	}
	if size != 3 {
		return 0, errors.New("failed to read the tempo value")
	}

	var bs [3]byte
	n, err := r.Read(bs[:])
	if err != nil {
		return 0, err
	}
	if n != 3 {
		return 0, errors.New("failed to read the tempo value")
	}

	b1 := bs[0]
	b2 := bs[1]
	b3 := bs[2]
	return (int32(b1) << 16) | (int32(b2) << 8) | int32(b3), nil
}

func discardData(r io.Reader) error {
	size, err := readIntVariableLength(r)
	if err != nil {
		return err
	}

	n, err := r.Read(make([]byte, size))
	if err != nil {
		return err
	}
	if n != int(size) {
		return errors.New("failed to read the data")
	}

	return nil
}

func (mf *MidiFile) GetLength() time.Duration {
	return mf.times[len(mf.times)-1]
}
