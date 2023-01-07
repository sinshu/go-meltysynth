package main

import (
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/sinshu/go-meltysynth/meltysynth"
)

var visWaveLen int = 512
var width float64 = 1024
var height float64 = 768

type audioStream struct {
	seq *meltysynth.MidiFileSequencer

	leftBuf  []float32
	rightBuf []float32

	visWave      []float64
	visWaveStart int
	mu           sync.Mutex
}

func NewAudioStream(seq *meltysynth.MidiFileSequencer) *audioStream {
	result := new(audioStream)
	result.seq = seq
	result.visWave = make([]float64, visWaveLen)
	result.visWaveStart = 0
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

	s.mu.Lock()
	for i := range samples {
		samples[i][0] = float64(s.leftBuf[i])
		samples[i][1] = float64(s.rightBuf[i])

		s.visWave[s.visWaveStart] = float64(s.leftBuf[i] + s.rightBuf[i])
		s.visWaveStart++
		if s.visWaveStart == visWaveLen {
			s.visWaveStart = 0
		}
	}
	s.mu.Unlock()

	return sampleCount, true
}

func (s *audioStream) Err() error {
	return nil
}

func (s *audioStream) getVisWave(data []float64) {
	s.mu.Lock()
	var p int = s.visWaveStart
	for i := 0; i < visWaveLen; i++ {
		data[i] = s.visWave[p]
		p++
		if p == visWaveLen {
			p = 0
		}
	}
	s.mu.Unlock()
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
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	stream := NewAudioStream(sequencer)
	sr := beep.SampleRate(sampleRate)
	speaker.Init(sr, sr.N(time.Second/10))
	speaker.Play(stream)

	imd := imdraw.New(nil)
	data := make([]float64, visWaveLen)

	for !win.Closed() {
		stream.getVisWave(data)
		imd.Reset()
		imd.Clear()
		imd.Color = colornames.Lightblue
		imd.EndShape = imdraw.NoEndShape
		for i := 0; i < visWaveLen; i++ {
			x := float64(i) / float64(visWaveLen) * width
			y := (height/4)*data[i] + (height / 2)
			imd.Push(pixel.V(x, y))
		}
		imd.Line(3)

		win.Clear(colornames.Darkblue)
		imd.Draw(win)
		win.Update()
	}

	speaker.Close()
}

func main() {
	pixelgl.Run(run)
}
