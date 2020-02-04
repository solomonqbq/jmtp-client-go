package v1

import (
    jmtpClient "github.com/jmtp/jmtp-client-go"
    "bytes"
    "github.com/jmtp/jmtp-client-go/util"
    "github.com/jmtp/jmtp-client-go/util/fieldcodec"
)

var (
    CommandAckDefineIns = &CommandAckDefine{}
    CommandAckCodecIns = &CommandAckCodec{}
)

type CommandAck struct {
    PacketId []byte
    Code    int
    Message string
    Payload []byte
}

func (c *CommandAck) GetPacketId() []byte {
    return c.PacketId
}

func (c *CommandAck) GetCode() int {
    return c.Code
}

func (c *CommandAck) GetMessage() string {
    return c.Message
}

func (c *CommandAck) GetPayload() []byte {
    return c.Payload
}

func (c *CommandAck) Define() jmtpClient.JmtpPacketDefine {
    return CommandAckDefineIns
}

func (c *CommandAck) HasAck() bool {
    return false
}

func NewCommandAck(packetId []byte, code int, message string, payload []byte) *CommandAck {
    return &CommandAck{
       PacketId: packetId,
       Code: code,
       Message: message,
       Payload: payload,
    }
}

type CommandAckDefine struct {

}

func (c *CommandAckDefine) PacketType() *jmtpClient.PacketType {
    return jmtpClient.CommandAck
}

func (c *CommandAckDefine) Code() byte {
    return c.PacketType().Code()
}

func (c *CommandAckDefine) CheckFlag(flagBits byte) bool {
    return flagBits == 0
}

func (c *CommandAckDefine) CreatePacket() jmtpClient.JmtpPacket {
    return &CommandAck{}
}

func (c *CommandAckDefine) Codec() jmtpClient.JmtpPacketCodec {
    return CommandAckCodecIns
}

func (c *CommandAckDefine) ProtocolDefine() jmtpClient.JmtpProtocolDefine {
    return JMTPV1ProtocolDefineInstance
}

type CommandAckCodec struct {

}

func (c *CommandAckCodec) EncodeBody(packet jmtpClient.JmtpPacket) ([]byte, error) {
    return encodeBody(packet)
}

func (c *CommandAckCodec) Decode(flagBits byte, input *bytes.Reader) (jmtpClient.JmtpPacket, error) {
    reader := util.NewJMTPDecodingReader(input)
    commandAck := &CommandAck{}
    if packetId, err := reader.ReadTinyBytesField();err != nil {
        return nil, err
    } else {
        commandAck.PacketId = packetId
    }
    if code, err := reader.ReadVarUnsignedInt();err != nil {
        return nil, err
    } else {
        commandAck.Code = code
    }
    if commandAck.Code != 0 {
        if msg, err := reader.ReadVShortField(fieldcodec.StringCodec);err != nil {
            return nil, err
        } else {
            commandAck.Message = msg.(string)
        }
    }
    if payload, err := reader.ReadAllByte();err != nil {
        return nil, err
    } else {
        commandAck.Payload = payload
    }
    return commandAck, nil
}

func (c *CommandAckCodec) GetFixedHeader(packet jmtpClient.JmtpPacket) (byte, error) {
    return packet.Define().PacketType().BuildHeader(), nil
}


