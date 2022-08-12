package meltysynth

import (
	"math"
	"time"
)

type MidiFileSequencer struct {
	synthesizer *Synthesizer
	midiFile    *MidiFile
	loop        bool
	blockWrote  int32
	currentTime time.Duration
	msgIndex    int32
	loopIndex   int32
}

func NewMidiFileSequencer(synthesizer *Synthesizer) *MidiFileSequencer {
	result := new(MidiFileSequencer)
	result.synthesizer = synthesizer
	return result
}

func (sequencer *MidiFileSequencer) Play(midiFile *MidiFile, loop bool) {

	sequencer.midiFile = midiFile
	sequencer.loop = loop

	sequencer.blockWrote = sequencer.synthesizer.BlockSize

	sequencer.currentTime = time.Duration(0)
	sequencer.msgIndex = 0
	sequencer.loopIndex = 0

	sequencer.synthesizer.Reset()
}

func (sequencer *MidiFileSequencer) Stop() {

	sequencer.midiFile = nil

	sequencer.synthesizer.Reset()
}

func (sequencer *MidiFileSequencer) Render(left []float32, right []float32) {

	wrote := int32(0)
	length := int32(len(left))
	for wrote < length {
		if sequencer.blockWrote == sequencer.synthesizer.BlockSize {
			sequencer.processEvents()
			sequencer.blockWrote = 0
			sequencer.currentTime += time.Duration(float64(time.Second) * float64(sequencer.synthesizer.BlockSize) / float64(sequencer.synthesizer.SampleRate))
		}

		srcRem := sequencer.synthesizer.BlockSize - sequencer.blockWrote
		dstRem := length - wrote
		rem := int32(math.Min(float64(srcRem), float64(dstRem)))

		sequencer.synthesizer.Render(left[wrote:wrote+rem], right[wrote:wrote+rem])

		sequencer.blockWrote += rem
		wrote += rem
	}
}

func (sequencer *MidiFileSequencer) processEvents() {

	if sequencer.midiFile == nil {
		return
	}

	msgLength := int32(len(sequencer.midiFile.messages))
	for sequencer.msgIndex < msgLength {
		time := sequencer.midiFile.times[sequencer.msgIndex]
		msg := sequencer.midiFile.messages[sequencer.msgIndex]
		if time <= sequencer.currentTime {
			if msg.getMessageType() == msg_Normal {
				sequencer.synthesizer.ProcessMidiMessage(int32(msg.channel), int32(msg.command), int32(msg.data1), int32(msg.data2))
			}
			sequencer.msgIndex++
		} else {
			break
		}
	}

	if sequencer.msgIndex == msgLength && sequencer.loop {
		sequencer.currentTime = 0
		sequencer.msgIndex = 0
		sequencer.synthesizer.NoteOffAll(false)
	}
}
