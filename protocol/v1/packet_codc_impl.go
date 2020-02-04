package v1

import (
    jmtpClient "github.com/jmtp/jmtp-client-go"
    "github.com/jmtp/jmtp-client-go/util"
    "github.com/jmtp/jmtp-client-go/util/fieldcodec"
    "errors"
    "fmt"
    "reflect"
)

func encodeBody(packet jmtpClient.JmtpPacket) ([]byte, error) {
    writer := util.NewJMTPEncodingWriter()
    var err error
    switch pack := packet.(type) {
    case *Connect:
        err = subpackageConnectBody(writer, pack)
    case *ConnectAck:
        err = subpacketConnectAckBody(writer, pack)
    case *Command:
        err = subpackageCommandBody(writer, pack)
    default:
        return nil, errors.New(
            fmt.Sprintf(
                "not implement type %s encoding writer func",
                reflect.TypeOf(packet).String()))
    }
    if err != nil {
        return nil, err
    }
    return writer.GetBytes(), nil
}

func subpackageConnectBody(writer *util.JMTPEncodingWriter, conn *Connect) error {
    if err := writer.WriteTinyField(conn.ProtocolName, fieldcodec.StringCodec);err != nil {
        return err
    }
    if err := writer.WriteUnsignedTiny(int(conn.ProtocolVersion));err != nil {
        return err
    }
    if err := writer.WriteVarUnsignedShort(int(conn.HeartbeatSec));err != nil {
        return err
    }
    if err := writer.WriteVarUnsignedShort(int(conn.SerializeType));err != nil {
        return err
    }
    if err := writer.WriteInt32(conn.ApplicationId);err != nil {
        return err
    }
    if err := writer.WriteInt32(conn.InstanceId);err != nil {
        return err
    }
    if err := writer.WriteTinyMap(conn.Tags, fieldcodec.StringCodec);err != nil {
        return err
    }
    return nil
}

func subpackageCommandBody(writer *util.JMTPEncodingWriter, command *Command) error {
    if err := writer.WriteTinyByte(command.PacketId);err != nil {
        return err
    }
    if err := writer.WriteShortField(command.Command, fieldcodec.StringCodec);err != nil {
        return err
    }
    if err := writer.WriteAllBytes(command.Payload);err != nil {
        return err
    }
    return nil
}

func subpacketConnectAckBody(writer *util.JMTPEncodingWriter, connectAck *ConnectAck) error {
    if err := writer.WriteVarUnsignedInt(connectAck.GetCode());err != nil {
        return err
    }
    if connectAck.GetCode() != 0 {
        if err := writer.WriteShortField(connectAck.GetMessage(), fieldcodec.StringCodec);err != nil {
            return err
        }
    }
    if connectAck.GetRetrySeconds() > 0 {
        if err := writer.WriteVarUnsignedShort(connectAck.GetRetrySeconds());err != nil {
            return err
        }
    }
    if connectAck.GetRedirectUrl() != "" {
        if err := writer.WriteTinyField(connectAck.GetRedirectUrl(), fieldcodec.StringCodec);err != nil {
            return err
        }
    }
    return nil
}
