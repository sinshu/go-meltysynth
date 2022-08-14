package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/sinshu/go-meltysynth/meltysynth"
)

const (
	screenWidth  = 640
	screenHeight = 480
	sampleRate   = 48000
	frequency    = 440
)

// stream is an infinite stream of 440 Hz sine wave.
type stream struct {
	sequencer *meltysynth.MidiFileSequencer
	leftBuf   []float32
	rightBuf  []float32
}

// Read is io.Reader's Read.
//
// Read fills the data with sine wave samples.
func (s *stream) Read(buf []byte) (int, error) {

	sampleCount := len(buf) / 4

	if s.leftBuf == nil {
		s.leftBuf = make([]float32, sampleCount)
		s.rightBuf = make([]float32, sampleCount)
	} else if len(s.leftBuf) < sampleCount {
		s.leftBuf = make([]float32, sampleCount)
		s.rightBuf = make([]float32, sampleCount)
	}

	s.sequencer.Render(s.leftBuf[0:sampleCount], s.rightBuf[0:sampleCount])

	for i := 0; i < sampleCount; i++ {

		b1 := int(32768 * s.leftBuf[i])
		if b1 < math.MinInt16 {
			b1 = math.MinInt16
		} else if b1 > math.MaxInt16 {
			b1 = math.MaxInt16
		}

		b2 := int(32768 * s.rightBuf[i])
		if b2 < math.MinInt16 {
			b2 = math.MinInt16
		} else if b2 > math.MaxInt16 {
			b2 = math.MaxInt16
		}

		buf[4*i] = byte(b1)
		buf[4*i+1] = byte(b1 >> 8)
		buf[4*i+2] = byte(b2)
		buf[4*i+3] = byte(b2 >> 8)
	}

	return len(buf), nil
}

// Close is io.Closer's Close.
func (s *stream) Close() error {
	return nil
}

type Game struct {
	audioContext *audio.Context
	player       *audio.Player
	sequencer    *meltysynth.MidiFileSequencer
}

func (g *Game) Update() error {
	if g.audioContext == nil {
		g.audioContext = audio.NewContext(sampleRate)
	}
	if g.player == nil {
		// Pass the (infinite) stream to NewPlayer.
		// After calling Play, the stream never ends as long as the player object lives.
		var err error
		g.player, err = g.audioContext.NewPlayer(&stream{sequencer: g.sequencer})
		if err != nil {
			return err
		}
		g.player.Play()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	msg := fmt.Sprintf("TPS: %0.2f\nThis is an example using infinite audio stream.", ebiten.CurrentFPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

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

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("SoundFont MIDI synthesis!!!")

	game := new(Game)
	game.sequencer = sequencer
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
