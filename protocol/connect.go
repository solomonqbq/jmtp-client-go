package protocol

import (
    jmtpClient "github.com/jmtp/jmtp-client-go"
    "io"
)

var ConnectCodec  = &ConnectPacketCodec{}

type Connect struct {
    ProtocolName    string
    ProtocolVersion int16
    HeartbeatSec    int16
    SerializeType   int16
    ApplicationId   int
    InstanceId      int
    Tags    map[string]string
}

func (c *Connect) GetProtocolName() string {
    return c.ProtocolName
}

func (c *Connect) GetProtocolVersion() int16 {
    return c.ProtocolVersion
}

func (c *Connect) GetHeartbeatSeconds() int16 {
    return c.HeartbeatSec
}

func (c *Connect) GetSerializeType() int16 {
    return c.SerializeType
}

func (c *Connect) GetApplicationId() int {
    return c.ApplicationId
}

func (c *Connect) GetInstanceId() int {
    return c.InstanceId
}

func (c *Connect) GetTags() map[string]string {
    return c.Tags
}

func (*Connect) Define() jmtpClient.JmtpPacketDefine {
    return nil
}

func (*Connect) HasAck() bool {
    return false
}

type ConnectPacketDefine struct {

}

func (c *ConnectPacketDefine) Codec() jmtpClient.JmtpPacketCodec {
    return nil
}

func (c *ConnectPacketDefine) ProtocolDefine() {
    panic("implement me")
}

func (c *ConnectPacketDefine) PacketType() *jmtpClient.PacketType {
    return jmtpClient.Connect
}

func (c *ConnectPacketDefine) Code() byte {
    return c.PacketType().Code()
}

func (c *ConnectPacketDefine) CheckFlag(flagBits byte) bool {
    return flagBits == 0
}

func (c *ConnectPacketDefine) CreatePacket() jmtpClient.JmtpPacket {
    return &Connect{}
}

type ConnectPacketCodec struct {

}

func (cpc *ConnectPacketCodec) EncodeBody(packet *jmtpClient.JmtpPacket) []byte {
    panic("implement me")
}

func (cpc *ConnectPacketCodec) Decode(flagBits byte, input io.Reader) {
    panic("implement me")
}

func (cpc *ConnectPacketCodec) GetFixedHeader(packet *jmtpClient.JmtpPacket) byte {
    panic("implement me")
}




