package v1

import (
    jmtpClient "github.com/jmtp/jmtp-client-go"
)

var JMTPV1PacketDefineInstance = newJMTPV1PacketDefine()

type jmtpV1PacketDefine struct {
    packetDefines []jmtpClient.JmtpPacketDefine
}

func newJMTPV1PacketDefine() *jmtpV1PacketDefine{
    jpd := &jmtpV1PacketDefine{}
    //jpd.packetDefines = make([]jmtpClient.JmtpPacketDefine, 0, 16)
    //fmt.Println(ConnectPacketDefineInstance.Code())
    //jpd.packetDefines[ConnectPacketDefineInstance.Code()] = ConnectPacketDefineInstance
    //jpd.packetDefines[PingPacketDefineIns.Code()] = PingPacketDefineIns
    //jpd.packetDefines[PongPacketDefineIns.Code()] = PongPacketDefineIns
    return jpd
}

func (j *jmtpV1PacketDefine) Get(code int) jmtpClient.JmtpPacketDefine{

    if code < 0 || code > 15 {
        return nil
    }
    return j.packetDefines[code]
}
