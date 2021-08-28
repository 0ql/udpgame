package util

const (
	TCP_CONNECT_REQUEST_PACKET    byte = 0
	TCP_CONNECT_PACKET            byte = 0
	TCP_CHUNK_REQUEST_PACKET      byte = 1
	TCP_PLAYERLIST_REQUEST_PACKET byte = 2
	TCP_PLAYERLIST_PACKET         byte = 2
	UDP_STATE_PACKET              byte = 0
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

func (packetBuilder *PacketBuilder) AddData(bytes []byte) {
	packetBuilder.data = append(packetBuilder.data, bytes...)
}

func (packetBuilder *PacketBuilder) AddByte(singleByte byte) {
	packetBuilder.AddData([]byte{singleByte})
}

func (packetBuilder *PacketBuilder) Build() []byte {
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
