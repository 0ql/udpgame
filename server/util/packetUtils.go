package util

type PacketDecoder struct {
	data  []byte
	index int
}

func NewPacketDecoder(data []byte) PacketDecoder {
	return PacketDecoder{
		data:  data,
		index: 1,
	}
}

func (packetDecoder *PacketDecoder) GetPacketType() byte {
	return packetDecoder.data[0]
}

func (packetDecoder *PacketDecoder) ExtractData(length int) []byte {
	data := packetDecoder.data[packetDecoder.index : packetDecoder.index+length]
	packetDecoder.index += length
	return data
}

func (packetDecoder *PacketDecoder) ExtractByte() byte {
	return packetDecoder.ExtractData(1)[0]
}

func (packetDecoder *PacketDecoder) SetIndex(i int) {
	packetDecoder.index = i
}
