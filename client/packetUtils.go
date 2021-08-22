package main

var (
	STATE_PACKET_ID byte = 1
)

type PacketBuilder struct {
	data []byte
}

func PacketBuilderNew(packetType byte) PacketBuilder {
	pb := PacketBuilder{}
	pb.data = make([]byte, 1)
	pb.data[0] = packetType
	return pb
}

func (packetBuilder *PacketBuilder) add_data(bytes []byte) {
	packetBuilder.data = append(packetBuilder.data, bytes...)
}

func (packetBuilder *PacketBuilder) add_byte(singleByte byte) {
	packetBuilder.add_data([]byte{singleByte})
}

func (packetBuilder *PacketBuilder) build() []byte {
	return packetBuilder.data
}

type PacketDecoder struct {
	data  []byte
	index int
}

func PacketDecoderNew(data []byte) PacketDecoder {
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