package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

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
	game     *Game
	leftBuf  []float32
	rightBuf []float32
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

	speed := s.game.GetSpeed()
	s.game.sequencer.SetSpeed(speed)
	s.game.sequencer.Render(s.leftBuf[0:sampleCount], s.rightBuf[0:sampleCount])

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
	mplusBigFont font.Face

	audioContext *audio.Context
	player       *audio.Player
	sequencer    *meltysynth.MidiFileSequencer

	mutex sync.Mutex
	speed int32
}

func (g *Game) Update() error {

	if g.audioContext == nil {
		g.audioContext = audio.NewContext(sampleRate)
	}

	if g.player == nil {
		// Pass the (infinite) stream to NewPlayer.
		// After calling Play, the stream never ends as long as the player object lives.
		var err error
		g.player, err = g.audioContext.NewPlayer(&stream{game: g})
		if err != nil {
			return err
		}
		g.player.Play()
	}

	newSpeed := g.speed
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		newSpeed -= 10
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		newSpeed += 10
	}
	if newSpeed < 10 {
		newSpeed = 10
	}
	if newSpeed > 1000 {
		newSpeed = 1000
	}
	if g.speed != newSpeed {
		g.mutex.Lock()
		g.speed = newSpeed
		g.mutex.Unlock()
	}

	return nil
}

func (g *Game) GetSpeed() float64 {
	g.mutex.Lock()
	result := float64(g.speed) / 100.0
	g.mutex.Unlock()
	return result
}

func (g *Game) Draw(screen *ebiten.Image) {
	msg := fmt.Sprintf("TPS: %0.2f\nThis is an example using infinite audio stream.", ebiten.CurrentFPS())
	text.Draw(screen, fmt.Sprintf("Playback speed: x%.1f", float64(g.speed)/100.0), g.mplusBigFont, 10, 100, color.White)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	tt, _ := opentype.Parse(fonts.MPlus1pRegular_ttf)

	mplusBigFont, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingFull,
	})

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
	game.mplusBigFont = mplusBigFont
	game.sequencer = sequencer
	game.speed = 100

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
