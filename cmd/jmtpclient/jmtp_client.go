package main

import (
    "fmt"
    "github.com/jmtp/jmtp-client-go"
    "github.com/jmtp/jmtp-client-go/protocol/v1"
    "net"
    "github.com/jmtp/jmtp-client-go/protocol"
    "os"
)

func main() {
    b := []byte{'a'}
    commandAck := v1.NewCommandAck(b, 1, "Test", nil)
    fmt.Printf("%v\n", commandAck)
    fmt.Printf("%v\n", commandAck)
    //jmtpClientTest()
}

func test(j jmtp_client_go.JmtpPacket) {
    switch v := j.(type) {
    case *v1.Connect:
        fmt.Println(v)
    }
}

func jmtpClientTest() {
    conn, err := net.Dial("tcp", "localhost:20560")
    fmt.Printf("%v\n", conn)
    if err != nil {
        fmt.Printf("err: %v", err)
        os.Exit(1)
    }
    defer conn.Close()
    fmt.Printf("Create connect packet\n")
    connectPacket := &v1.Connect{}
    connectPacket.ApplicationId = 999
    connectPacket.InstanceId = 999
    connectPacket.Tags = map[string]interface{} {
        "appName": "test",
        "site": "test",
    }
    connectPacket.SerializeType = 1
    connectPacket.HeartbeatSec = 10
    connectPacket.ProtocolName = "JMTP"
    connectPacket.ProtocolVersion = 1
    fmt.Printf("Create connect packet success\n")
    out, err := protocol.PacketEncoder(connectPacket)
    fmt.Printf("message length: %d\n", len(out))
    if err != nil {
        panic(err)
    }
    num, err := conn.Write(out)
    fmt.Printf("writer num %d, err: %v", num, err)
}


