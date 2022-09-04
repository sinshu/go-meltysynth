package main

import (
	"encoding/binary"
	"math"
	"os"
	"time"

	"github.com/sinshu/go-meltysynth/meltysynth"
)

func main() {

	simpleChord()
	flourish()
}

func simpleChord() {

	// Load the SoundFont.
	sf2, _ := os.Open("TimGM6mb.sf2")
	soundFont, _ := meltysynth.NewSoundFont(sf2)
	sf2.Close()

	// Create the synthesizer.
	settings := meltysynth.NewSynthesizerSettings(44100)
	synthesizer, _ := meltysynth.NewSynthesizer(soundFont, settings)

	// Play some notes (middle C, E, G).
	synthesizer.NoteOn(0, 60, 100)
	synthesizer.NoteOn(0, 64, 100)
	synthesizer.NoteOn(0, 67, 100)

	// The output buffer (3 seconds).
	length := 3 * settings.SampleRate
	left := make([]float32, length)
	right := make([]float32, length)

	// Render the waveform.
	synthesizer.Render(left, right)

	writePcmInterleavedInt16(left, right, "simpleChord.pcm")
}

func flourish() {

	// Load the SoundFont.
	sf2, _ := os.Open("TimGM6mb.sf2")
	soundFont, _ := meltysynth.NewSoundFont(sf2)
	sf2.Close()

	// Create the synthesizer.
	settings := meltysynth.NewSynthesizerSettings(44100)
	synthesizer, _ := meltysynth.NewSynthesizer(soundFont, settings)

	// Load the MIDI file.
	mid, _ := os.Open("C:\\Windows\\Media\\flourish.mid")
	midiFile, _ := meltysynth.NewMidiFile(mid)
	mid.Close()

	// Create the MIDI sequencer.
	sequencer := meltysynth.NewMidiFileSequencer(synthesizer)
	sequencer.Play(midiFile, true)

	// The output buffer.
	length := int(float64(settings.SampleRate) * float64(midiFile.GetLength()) / float64(time.Second))
	left := make([]float32, length)
	right := make([]float32, length)

	// Render the waveform.
	sequencer.Render(left, right)

	writePcmInterleavedInt16(left, right, "flourish.pcm")
}

func writePcmInterleavedInt16(left []float32, right []float32, path string) {

	length := len(left)

	max := 0.0

	for i := 0; i < length; i++ {
		absLeft := math.Abs(float64(left[i]))
		absRight := math.Abs(float64(right[i]))
		if max < absLeft {
			max = absLeft
		}
		if max < absRight {
			max = absRight
		}
	}

	a := 32768 * float32(0.99/max)

	data := make([]int16, 2*length)

	for i := 0; i < length; i++ {
		data[2*i] = int16(a * left[i])
		data[2*i+1] = int16(a * right[i])
	}

	pcm, _ := os.Create(path)
	binary.Write(pcm, binary.LittleEndian, data)
	pcm.Close()
}
