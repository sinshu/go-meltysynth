package meltysynth

import (
	"errors"
	"os"
	"testing"
)

const (
	envGS = "MELTYSYNTH_GS"
	envGM = "MELTYSYNTH_GM"

	defaultPathGS = "GeneralUser GS MuseScore v1.442.sf2"
	defaultPathGM = "TimGM6mb.sf2"
)

func loadGS(t *testing.T) *SoundFont {
	return loadSoundFont(t, envGS, defaultPathGS)
}

func loadGM(t *testing.T) *SoundFont {
	return loadSoundFont(t, envGM, defaultPathGM)
}

func loadSoundFont(t *testing.T, env, defaultPath string) *SoundFont {
	var useDefault bool
	p := os.Getenv(env)
	if len(p) == 0 {
		useDefault = true
		p = defaultPath
		// t.Skipf("missing environment variable %q to load soundfont", env)
	}
	f, err := os.Open(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if useDefault {
				t.Skipf("missing environment variable %q to load soundfont, default path %q not found", env, defaultPath)
			}
			t.Fatalf("envionrment variable %q set to %q, but file does not exist", env, p)
		}
		t.Fatal(err)
	}
	sf, err := NewSoundFont(f)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}

	return sf
}
