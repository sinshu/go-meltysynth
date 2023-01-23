package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"time"

	"github.com/sinshu/go-meltysynth/meltysynth"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	sf2Path := flag.String("sf2", "", "Sound font file location")
	outPath := flag.String("o", "out.pcm", "file to write pcm file to")
	midiPath := flag.String("midi", "", `midi file to synth, or "chord" for example`)
	flag.Parse()

	if len(*sf2Path) == 0 {
		flag.PrintDefaults()
		return fmt.Errorf("missing sf2 path")
	}
	if len(*outPath) == 0 {
		flag.PrintDefaults()
		return fmt.Errorf("missing output path")
	}
	if len(*midiPath) == 0 {
		flag.PrintDefaults()
		return fmt.Errorf("missing midi path")
	}

	sf2, err := os.Open(*sf2Path)
	if err != nil {
		return err
	}
	soundFont, err := meltysynth.NewSoundFont(sf2)
	sf2.Close()
	if err != nil {
		return err
	}

	switch *midiPath {
	default:
		err = midi(soundFont, *midiPath, *outPath)
	case "chord":
		err = simpleChord(soundFont, *outPath)
	}
	if err != nil {
		return err
	}
	return nil
}

func simpleChord(soundFont *meltysynth.SoundFont, outputFile string) error {
	// Create the synthesizer.
	settings := meltysynth.NewSynthesizerSettings(44100)
	synthesizer, err := meltysynth.NewSynthesizer(soundFont, settings)
	if err != nil {
		return err
	}

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

	return writeFile(left, right, outputFile)
}

func midi(soundFont *meltysynth.SoundFont, midiFilePath string, outputFile string) error {
	// Create the synthesizer.
	settings := meltysynth.NewSynthesizerSettings(44100)
	synthesizer, err := meltysynth.NewSynthesizer(soundFont, settings)
	if err != nil {
		return err
	}

	// Load the MIDI file.
	mid, err := os.Open(midiFilePath)
	if err != nil {
		return err
	}
	midiFile, err := meltysynth.NewMidiFile(mid)
	mid.Close()
	if err != nil {
		return err
	}

	// Create the MIDI sequencer.
	sequencer := meltysynth.NewMidiFileSequencer(synthesizer)
	sequencer.Play(midiFile, true)

	// The output buffer.
	length := int(float64(settings.SampleRate) * float64(midiFile.GetLength()) / float64(time.Second))
	left := make([]float32, length)
	right := make([]float32, length)

	// Render the waveform.
	sequencer.Render(left, right)

	return writeFile(left, right, outputFile)
}

func writeFile(left []float32, right []float32, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return writePCMInterleavedInt16(left, right, f)
}

func writePCMInterleavedInt16(left []float32, right []float32, pcm io.Writer) error {
	length := len(left)
	var max float64

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

	return binary.Write(pcm, binary.LittleEndian, data)
}
