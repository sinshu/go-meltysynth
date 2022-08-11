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

func newChannel(synthesizer *Synthesizer, isPercussionChannel bool) *channel {

	result := new(channel)

	result.synthesizer = synthesizer
	result.isPercussionChannel = isPercussionChannel

	result.reset()

	return result
}

func (channel *channel) reset() {

	if channel.isPercussionChannel {
		channel.bankNumber = 128
	} else {
		channel.bankNumber = 0
	}

	channel.patchNumber = 0

	channel.modulation = 0
	channel.volume = 100 << 7
	channel.pan = 64 << 7
	channel.expression = 127 << 7
	channel.holdPedal = false

	channel.reverbSend = 40
	channel.chorusSend = 0

	channel.rpn = -1
	channel.pitchBendRange = 2 << 7
	channel.coarseTune = 0
	channel.fineTune = 8192

	channel.pitchBend = 0
}

func (channel *channel) resetAllControllers() {

	channel.modulation = 0
	channel.expression = 127 << 7
	channel.holdPedal = false

	channel.rpn = -1

	channel.pitchBend = 0
}

func (channel *channel) setBank(value int32) {

	channel.bankNumber = value

	if channel.isPercussionChannel {
		channel.bankNumber += 128
	}
}

func (channel *channel) setPatch(value int32) {
	channel.patchNumber = value
}

func (channel *channel) setModulationCoarse(value int32) {
	channel.modulation = int16((int32(channel.modulation) & 0x7F) | (value << 7))
}

func (channel *channel) setModulationFine(value int32) {
	channel.modulation = int16((int32(channel.modulation) & 0xFF80) | value)
}

func (channel *channel) setVolumeCoarse(value int32) {
	channel.volume = int16((int32(channel.volume) & 0x7F) | (value << 7))
}

func (channel *channel) setVolumeFine(value int32) {
	channel.volume = int16((int32(channel.volume) & 0xFF80) | value)
}

func (channel *channel) setPanCoarse(value int32) {
	channel.pan = int16((int32(channel.pan) & 0x7F) | (value << 7))
}

func (channel *channel) setPanFine(value int32) {
	channel.pan = int16((int32(channel.pan) & 0xFF80) | value)
}

func (channel *channel) setExpressionCoarse(value int32) {
	channel.expression = int16((int32(channel.expression) & 0x7F) | (value << 7))
}

func (channel *channel) setExpressionFine(value int32) {
	channel.expression = int16((int32(channel.expression) & 0xFF80) | value)
}

func (channel *channel) setHoldPedal(value int32) {
	channel.holdPedal = value >= 64
}

func (channel *channel) setReverbSend(value int32) {
	channel.reverbSend = byte(value)
}

func (channel *channel) setChorusSend(value int32) {
	channel.chorusSend = byte(value)
}

func (channel *channel) setRpnCoarse(value int32) {
	channel.rpn = int16((int32(channel.rpn) & 0x7F) | (value << 7))
}

func (channel *channel) setRpnFine(value int32) {
	channel.rpn = int16((int32(channel.rpn) & 0xFF80) | value)
}

func (channel *channel) dataEntryCoarse(value int32) {
	switch channel.rpn {
	case 0:
		channel.pitchBendRange = int16((int32(channel.pitchBendRange) & 0x7F) | (value << 7))
	case 1:
		channel.fineTune = int16((int32(channel.fineTune) & 0x7F) | (value << 7))
	case 2:
		channel.coarseTune = int16(value - 64)
	}
}

func (channel *channel) dataEntryFine(value int32) {
	switch channel.rpn {
	case 0:
		channel.pitchBendRange = int16((int32(channel.pitchBendRange) & 0xFF80) | value)
	case 1:
		channel.fineTune = int16((int32(channel.fineTune) & 0xFF80) | value)
	}
}

func (channel *channel) setPitchBend(value1 int32, value2 int32) {
	channel.pitchBend = (float32(1) / float32(8192)) * float32((value1|(value2<<7))-8192)
}

func (channel *channel) getModulation() float32 {
	return (float32(50) / float32(16383)) * float32(channel.modulation)
}

func (channel *channel) getVolume() float32 {
	return (float32(1) / float32(16383)) * float32(channel.volume)
}

func (channel *channel) getPan() float32 {
	return (float32(100)/float32(16383))*float32(channel.pan) - 50
}

func (channel *channel) getExpression() float32 {
	return (float32(1) / float32(16383)) * float32(channel.expression)
}

func (channel *channel) getReverbSend() float32 {
	return (float32(1) / float32(127)) * float32(channel.reverbSend)
}

func (channel *channel) getChorusSend() float32 {
	return (float32(1) / float32(127)) * float32(channel.chorusSend)
}

func (channel *channel) getPitchBendRange() float32 {
	return float32(channel.pitchBendRange>>7) + 0.01*float32(channel.pitchBendRange&0x7F)
}

func (channel *channel) getTune() float32 {
	return float32(channel.coarseTune) + (float32(1)/float32(8192))*float32(channel.fineTune-8192)
}

func (channel *channel) getPitchBend() float32 {
	return channel.getPitchBendRange() * channel.pitchBend
}
