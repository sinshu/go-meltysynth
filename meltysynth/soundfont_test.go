package meltysynth

import (
	"testing"
)

func TestTimGM6mb_SoundFont(t *testing.T) {
	soundFont := loadGM(t)

	TimGM6mb_SoundFontInfo(t, soundFont)
	TimGM6mb_SoundFontSampleData(t, soundFont)
}

func TimGM6mb_SoundFontInfo(t *testing.T, soundFont *SoundFont) {
	if soundFont.Info.Version.Major != 2 {
		t.Fail()
	}

	if soundFont.Info.Version.Minor != 1 {
		t.Fail()
	}

	if soundFont.Info.TargetSoundEngine != "EMU8000" {
		t.Fail()
	}

	if soundFont.Info.BankName != "TimGM6mb1.sf2" {
		t.Fail()
	}

	if soundFont.Info.RomName != "" {
		t.Fail()
	}

	if soundFont.Info.RomVersion.Major != 0 {
		t.Fail()
	}

	if soundFont.Info.RomVersion.Minor != 0 {
		t.Fail()
	}

	if soundFont.Info.CreationDate != "" {
		t.Fail()
	}

	if soundFont.Info.Auther != "" {
		t.Fail()
	}

	if soundFont.Info.TargetProduct != "" {
		t.Fail()
	}

	if soundFont.Info.Copyright != "" {
		t.Fail()
	}

	if soundFont.Info.Comments != "" {
		t.Fail()
	}

	if soundFont.Info.Tools != "Awave Studio v8.5" {
		t.Fail()
	}
}

func TimGM6mb_SoundFontSampleData(t *testing.T, soundFont *SoundFont) {
	if soundFont.BitsPerSample != 16 {
		t.Fail()
	}

	if len(soundFont.WaveData) != 2882168 {
		t.Fail()
	}

	first16 := []int16{0, 0, 0, 0, 0, 0, 0, -1, -5, -12, -16, -16, -13, -14, -16, -18}
	for i := 0; i < len(first16); i++ {
		if soundFont.WaveData[i] != first16[i] {
			t.Fail()
		}
	}

	last40 := []int16{-4674, 4669, 8114, -374, -618, -2700, -8628, 2168, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	for i := 0; i < len(last40); i++ {
		if soundFont.WaveData[len(soundFont.WaveData)-len(last40)+i] != last40[i] {
			t.Fail()
		}
	}
}
