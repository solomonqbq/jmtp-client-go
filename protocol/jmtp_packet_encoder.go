package protocol

import (
    jmtpClient "github.com/jmtp/jmtp-client-go"
    "bytes"
    "github.com/jmtp/jmtp-client-go/util"
    "bufio"
    "fmt"
    "github.com/jmtp/jmtp-client-go/protocol/v1"
    "time"
    "errors"
)

const packetMinSize = 3

func PacketEncoder(packet jmtpClient.JmtpPacket) ([]byte, error) {
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
        if err := util.EncodeRemainingLength(len(byteBody), &out);err != nil {
            return nil, err
        }
        out.Write(byteBody)
    }
    return out.Bytes(), nil
}

func PacketDecoder(reader *bufio.Reader, packetsChain chan jmtpClient.JmtpPacket,errorChain chan error) error {
    for {
        if readPkg, err := reader.Peek(3);err == nil && len(readPkg) >= packetMinSize {
            header, err := reader.ReadByte()
            if err != nil {
                continue
            }
            crc, err := reader.ReadByte()
            if err != nil {
                continue
            }
            crc = crc ^ 0xFF
            if header != crc {
                reader.Discard(reader.Buffered())
                continue
            }
            packetDefine := v1.JMTPV1ProtocolDefineInstance.PacketDefine((header >> 4) & 0x0F)
            if packetDefine == nil {
                continue
            }
            flagBits := header & 0x0F
            if !packetDefine.CheckFlag(flagBits) {
                // close & continue
            }
            remainingLength, err := util.DecodeRemainingLength(reader)
            if remainingLength < 0 {
                // close & continue
            }
            if remainingLength > 0 {
                retryTimes := 0
                for {
                    if retryTimes > 10 {
                        err := errors.New(
                            fmt.Sprintf(
                                "can't read read enough byte stream, expect %d", remainingLength))
                        return err
                    }
                    payload, err := reader.Peek(remainingLength)
                    if err == nil && len(payload) == remainingLength {
                        if discarded, err := reader.Discard(len(payload));err != nil {
                            return err
                        } else if discarded != len(payload) {
                            return errors.New(
                                fmt.Sprintf("discarded length %d not equal payload length %d",
                                    len(payload),
                                    discarded))
                        }
                        packet, err:= packetDefine.Codec().Decode(flagBits, bytes.NewReader(payload))
                        if err != nil {
                            errorChain <- err
                            break
                        }
                        packetsChain <- packet
                        break
                    } else {
                        retryTimes += 1
                    }
                }
            } else {
                packet, err := packetDefine.Codec().Decode(flagBits, nil)
                if err != nil {
                    errorChain <- err
                } else {
                    packetsChain <- packet
                }
            }
        } else {
            time.Sleep(time.Duration(1) * time.Millisecond)
        }
    }
}
