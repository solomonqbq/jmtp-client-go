package jmtp_client_go

import "io"

type JmtpPacket interface {
    Define() JmtpPacketDefine
    HasAck() bool
}

type ConnectPacket interface {
    JmtpPacket
    GetProtocolName() string
    GetProtocolVersion() int16
    GetHeartbeatSeconds() int16
    GetSerializeType() int16
    GetApplicationId() int
    GetInstanceId() int
    GetTags()   map[string]string
}

type JmtpPacketCodec interface {
    EncodeBody(packet *JmtpPacket) []byte
    Decode(flagBits byte, input io.Reader)
    GetFixedHeader(packet *JmtpPacket) byte
}

type JmtpPacketDefine interface {
    PacketType() *PacketType
    Code() byte
    CheckFlag(flagBits byte) bool
    CreatePacket() JmtpPacket
    Codec() JmtpPacketCodec
    ProtocolDefine()
}

type JmtpProtocolDefine interface {
    Name() string
    Version() int16
    PacketDefine(code int) JmtpPacketDefine

}

var Connect = NewPacketType(byte(0x1))
var ConnectAck = NewPacketType(byte(0x2))
var Ping = NewPacketType(byte(0x3))
var Pong = NewPacketType(byte(0x4))
var Disconnect = NewPacketType(byte(0x5))
var Report = NewPacketType(byte(0x6))
var ReportAck = NewPacketType(byte(0x7))
var Command = NewPacketType(byte(0x8))
var CommandAck = NewPacketType(byte(0x9))


type PacketType struct {
    code byte
    headerBits byte
}

func (p *PacketType) Check(t PacketType) bool {
    return *p == t
}

func (p *PacketType) Code() byte {
    return p.code
}

func (p *PacketType) buildHeader(flagBits ...byte) byte{
    header := p.headerBits
    for _, flagBit := range flagBits {
        header |= flagBit
    }
    return header
}

func NewPacketType(t byte) *PacketType {
    pt := &PacketType{
        code: t,
    }
    pt.headerBits = t << 4
    return pt
}

type ConnectOption struct {
    HeartbeatSeconds    int16
    SerializeType   int16
    ApplicationId   int
    InstanceId  int
    Tags    map[string]string
}
