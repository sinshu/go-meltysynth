package meltysynth_test

import (
	"math"
	"testing"

	"github.com/sinshu/go-meltysynth/meltysynth"
)

func areEqual(t *testing.T, x float64, y float64) {
	max := math.Max(math.Abs(x), math.Abs(y))
	limit := max / float64(1000)
	delta := math.Abs(x - y)
	if delta > limit {
		t.Fail()
	}
}

func checkInstrumentRegion(t *testing.T, region *meltysynth.InstrumentRegion, values []float64) {
	areEqual(t, float64(meltysynth.GetInstrumentSampleStart(region)), values[0])
	areEqual(t, float64(meltysynth.GetInstrumentSampleEnd(region)), values[1])
	areEqual(t, float64(meltysynth.GetInstrumentSampleStartLoop(region)), values[2])
	areEqual(t, float64(meltysynth.GetInstrumentSampleEndLoop(region)), values[3])
	areEqual(t, float64(meltysynth.GetInstrumentStartAddressOffset(region)), values[4])
	areEqual(t, float64(meltysynth.GetInstrumentEndAddressOffset(region)), values[5])
	areEqual(t, float64(meltysynth.GetInstrumentStartLoopAddressOffset(region)), values[6])
	areEqual(t, float64(meltysynth.GetInstrumentEndLoopAddressOffset(region)), values[7])
	areEqual(t, float64(meltysynth.GetInstrumentModulationLfoToPitch(region)), values[8])
	areEqual(t, float64(meltysynth.GetInstrumentVibratoLfoToPitch(region)), values[9])
	areEqual(t, float64(meltysynth.GetInstrumentModulationEnvelopeToPitch(region)), values[10])
	areEqual(t, float64(meltysynth.GetInstrumentInitialFilterCutoffFrequency(region)), values[11])
	areEqual(t, float64(meltysynth.GetInstrumentInitialFilterQ(region)), values[12])
	areEqual(t, float64(meltysynth.GetInstrumentModulationLfoToFilterCutoffFrequency(region)), values[13])
	areEqual(t, float64(meltysynth.GetInstrumentModulationEnvelopeToFilterCutoffFrequency(region)), values[14])
	areEqual(t, float64(meltysynth.GetInstrumentModulationLfoToVolume(region)), values[15])
	areEqual(t, float64(meltysynth.GetInstrumentChorusEffectsSend(region)), values[16])
	areEqual(t, float64(meltysynth.GetInstrumentReverbEffectsSend(region)), values[17])
	areEqual(t, float64(meltysynth.GetInstrumentPan(region)), values[18])
	areEqual(t, float64(meltysynth.GetInstrumentDelayModulationLfo(region)), values[19])
	areEqual(t, float64(meltysynth.GetInstrumentFrequencyModulationLfo(region)), values[20])
	areEqual(t, float64(meltysynth.GetInstrumentDelayVibratoLfo(region)), values[21])
	areEqual(t, float64(meltysynth.GetInstrumentFrequencyVibratoLfo(region)), values[22])
	areEqual(t, float64(meltysynth.GetInstrumentDelayModulationEnvelope(region)), values[23])
	areEqual(t, float64(meltysynth.GetInstrumentAttackModulationEnvelope(region)), values[24])
	areEqual(t, float64(meltysynth.GetInstrumentHoldModulationEnvelope(region)), values[25])
	areEqual(t, float64(meltysynth.GetInstrumentDecayModulationEnvelope(region)), values[26])
	areEqual(t, float64(meltysynth.GetInstrumentSustainModulationEnvelope(region)), values[27])
	areEqual(t, float64(meltysynth.GetInstrumentReleaseModulationEnvelope(region)), values[28])
	areEqual(t, float64(meltysynth.GetInstrumentKeyNumberToModulationEnvelopeHold(region)), values[29])
	areEqual(t, float64(meltysynth.GetInstrumentKeyNumberToModulationEnvelopeDecay(region)), values[30])
	areEqual(t, float64(meltysynth.GetInstrumentDelayVolumeEnvelope(region)), values[31])
	areEqual(t, float64(meltysynth.GetInstrumentAttackVolumeEnvelope(region)), values[32])
	areEqual(t, float64(meltysynth.GetInstrumentHoldVolumeEnvelope(region)), values[33])
	areEqual(t, float64(meltysynth.GetInstrumentDecayVolumeEnvelope(region)), values[34])
	areEqual(t, float64(meltysynth.GetInstrumentSustainVolumeEnvelope(region)), values[35])
	areEqual(t, float64(meltysynth.GetInstrumentReleaseVolumeEnvelope(region)), values[36])
	areEqual(t, float64(meltysynth.GetInstrumentKeyNumberToVolumeEnvelopeHold(region)), values[37])
	areEqual(t, float64(meltysynth.GetInstrumentKeyNumberToVolumeEnvelopeDecay(region)), values[38])
	areEqual(t, float64(meltysynth.GetInstrumentKeyRangeStart(region)), values[39])
	areEqual(t, float64(meltysynth.GetInstrumentKeyRangeEnd(region)), values[40])
	areEqual(t, float64(meltysynth.GetInstrumentVelocityRangeStart(region)), values[41])
	areEqual(t, float64(meltysynth.GetInstrumentVelocityRangeEnd(region)), values[42])
	areEqual(t, float64(meltysynth.GetInstrumentInitialAttenuation(region)), values[43])
	areEqual(t, float64(meltysynth.GetInstrumentCoarseTune(region)), values[44])
	areEqual(t, float64(meltysynth.GetInstrumentFineTune(region)), values[45])
	areEqual(t, float64(meltysynth.GetInstrumentSampleModes(region)), values[46])
	areEqual(t, float64(meltysynth.GetInstrumentScaleTuning(region)), values[47])
	areEqual(t, float64(meltysynth.GetInstrumentExclusiveClass(region)), values[48])
	areEqual(t, float64(meltysynth.GetInstrumentRootKey(region)), values[49])
}
