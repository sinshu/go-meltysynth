package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/sinshu/go-meltysynth/meltysynth"
)

type audioStream struct {
	seq      *meltysynth.MidiFileSequencer
	leftBuf  []float32
	rightBuf []float32
}

func NewAudioStream(seq *meltysynth.MidiFileSequencer) *audioStream {
	result := new(audioStream)
	result.seq = seq
	return result
}

func (s *audioStream) Stream(samples [][2]float64) (n int, ok bool) {
	sampleCount := len(samples)

	if s.leftBuf == nil {
		s.leftBuf = make([]float32, sampleCount)
		s.rightBuf = make([]float32, sampleCount)
	} else if len(s.leftBuf) < sampleCount {
		s.leftBuf = make([]float32, sampleCount)
		s.rightBuf = make([]float32, sampleCount)
	}

	s.seq.Render(s.leftBuf[0:sampleCount], s.rightBuf[0:sampleCount])

	for i := range samples {
		samples[i][0] = float64(s.leftBuf[i])
		samples[i][1] = float64(s.rightBuf[i])
	}

	return sampleCount, true
}

func (s *audioStream) Err() error {
	return nil
}

func run() {
	var sampleRate int32 = 44100

	sf2, _ := os.Open("TimGM6mb.sf2")
	soundFont, _ := meltysynth.NewSoundFont(sf2)
	sf2.Close()

	settings := meltysynth.NewSynthesizerSettings(sampleRate)
	synthesizer, _ := meltysynth.NewSynthesizer(soundFont, settings)

	mid, _ := os.Open("C:\\Windows\\Media\\flourish.mid")
	midiFile, _ := meltysynth.NewMidiFile(mid)
	mid.Close()

	sequencer := meltysynth.NewMidiFileSequencer(synthesizer)
	sequencer.Play(midiFile, true)

	cfg := pixelgl.WindowConfig{
		Title:  "MIDI music playback!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Skyblue)

	sr := beep.SampleRate(sampleRate)
	speaker.Init(sr, sr.N(time.Second/10))
	speaker.Play(NewAudioStream(sequencer))

	for !win.Closed() {
		win.Update()
	}

	speaker.Close()
}

func main() {
	pixelgl.Run(run)
}
