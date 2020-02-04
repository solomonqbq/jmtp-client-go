package protocol

import (
    jmtpClinet "github.com/jmtp/jmtp-client-go"
    "bytes"
    "github.com/jmtp/jmtp-client-go/protocol/v1"
)

func PacketEncoder(packet jmtpClinet.JmtpPacket) ([]byte, error) {
    var out bytes.Buffer
    codec := packet.Define().Codec()
    head, err := codec.GetFixedHeader(packet)
    if err != nil {
        return nil, err
    }
    out.WriteByte(head)
    out.WriteByte(head ^ 0xFF)
    byteBody, err := codec.EncodeBody(packet)
    if byteBody == nil || len(byteBody) == 0 {
        out.WriteByte(0x00)
    } else {
        if err := v1.EncodeRemainingLength(len(byteBody), &out);err != nil {
            return nil, err
        }
        out.Write(byteBody)
    }
    return out.Bytes(), nil
}
