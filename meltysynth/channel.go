package meltysynth

type channel struct {
	synthesizer         *Synthesizer
	isPercussionChannel bool

	bankNumber  int32
	patchNumber int32

	modulation int16
	volume     int16
	pan        int16
	expression int16
	holdPedal  bool

	reverbSend byte
	chorusSend byte

	rpn            int16
	pitchBendRange int16
	coarseTune     int16
	fineTune       int16

	pitchBend float32
}

func newChannel(s *Synthesizer, isPercussionChannel bool) *channel {
	result := new(channel)

	result.synthesizer = s
	result.isPercussionChannel = isPercussionChannel

	result.reset()

	return result
}

func (ch *channel) reset() {
	if ch.isPercussionChannel {
		ch.bankNumber = 128
	} else {
		ch.bankNumber = 0
	}

	ch.patchNumber = 0

	ch.modulation = 0
	ch.volume = 100 << 7
	ch.pan = 64 << 7
	ch.expression = 127 << 7
	ch.holdPedal = false

	ch.reverbSend = 40
	ch.chorusSend = 0

	ch.rpn = -1
	ch.pitchBendRange = 2 << 7
	ch.coarseTune = 0
	ch.fineTune = 8192

	ch.pitchBend = 0
}

func (ch *channel) resetAllControllers() {
	ch.modulation = 0
	ch.expression = 127 << 7
	ch.holdPedal = false

	ch.rpn = -1

	ch.pitchBend = 0
}

func (ch *channel) setBank(value int32) {
	ch.bankNumber = value

	if ch.isPercussionChannel {
		ch.bankNumber += 128
	}
}

func (ch *channel) setPatch(value int32) {
	ch.patchNumber = value
}

func (ch *channel) setModulationCoarse(value int32) {
	ch.modulation = int16((int32(ch.modulation) & 0x7F) | (value << 7))
}

func (ch *channel) setModulationFine(value int32) {
	ch.modulation = int16((int32(ch.modulation) & 0xFF80) | value)
}

func (ch *channel) setVolumeCoarse(value int32) {
	ch.volume = int16((int32(ch.volume) & 0x7F) | (value << 7))
}

func (ch *channel) setVolumeFine(value int32) {
	ch.volume = int16((int32(ch.volume) & 0xFF80) | value)
}

func (ch *channel) setPanCoarse(value int32) {
	ch.pan = int16((int32(ch.pan) & 0x7F) | (value << 7))
}

func (ch *channel) setPanFine(value int32) {
	ch.pan = int16((int32(ch.pan) & 0xFF80) | value)
}

func (ch *channel) setExpressionCoarse(value int32) {
	ch.expression = int16((int32(ch.expression) & 0x7F) | (value << 7))
}

func (ch *channel) setExpressionFine(value int32) {
	ch.expression = int16((int32(ch.expression) & 0xFF80) | value)
}

func (ch *channel) setHoldPedal(value int32) {
	ch.holdPedal = value >= 64
}

func (ch *channel) setReverbSend(value int32) {
	ch.reverbSend = byte(value)
}

func (ch *channel) setChorusSend(value int32) {
	ch.chorusSend = byte(value)
}

func (ch *channel) setRpnCoarse(value int32) {
	ch.rpn = int16((int32(ch.rpn) & 0x7F) | (value << 7))
}

func (ch *channel) setRpnFine(value int32) {
	ch.rpn = int16((int32(ch.rpn) & 0xFF80) | value)
}

func (ch *channel) dataEntryCoarse(value int32) {
	switch ch.rpn {
	case 0:
		ch.pitchBendRange = int16((int32(ch.pitchBendRange) & 0x7F) | (value << 7))
	case 1:
		ch.fineTune = int16((int32(ch.fineTune) & 0x7F) | (value << 7))
	case 2:
		ch.coarseTune = int16(value - 64)
	}
}

func (ch *channel) dataEntryFine(value int32) {
	switch ch.rpn {
	case 0:
		ch.pitchBendRange = int16((int32(ch.pitchBendRange) & 0xFF80) | value)
	case 1:
		ch.fineTune = int16((int32(ch.fineTune) & 0xFF80) | value)
	}
}

func (ch *channel) setPitchBend(value1 int32, value2 int32) {
	ch.pitchBend = (float32(1) / float32(8192)) * float32((value1|(value2<<7))-8192)
}

func (ch *channel) getModulation() float32 {
	return (float32(50) / float32(16383)) * float32(ch.modulation)
}

func (ch *channel) getVolume() float32 {
	return (float32(1) / float32(16383)) * float32(ch.volume)
}

func (ch *channel) getPan() float32 {
	return (float32(100)/float32(16383))*float32(ch.pan) - 50
}

func (ch *channel) getExpression() float32 {
	return (float32(1) / float32(16383)) * float32(ch.expression)
}

func (ch *channel) getReverbSend() float32 {
	return (float32(1) / float32(127)) * float32(ch.reverbSend)
}

func (ch *channel) getChorusSend() float32 {
	return (float32(1) / float32(127)) * float32(ch.chorusSend)
}

func (ch *channel) getPitchBendRange() float32 {
	return float32(ch.pitchBendRange>>7) + 0.01*float32(ch.pitchBendRange&0x7F)
}

func (ch *channel) getTune() float32 {
	return float32(ch.coarseTune) + (float32(1)/float32(8192))*float32(ch.fineTune-8192)
}

func (ch *channel) getPitchBend() float32 {
	return ch.getPitchBendRange() * ch.pitchBend
}
