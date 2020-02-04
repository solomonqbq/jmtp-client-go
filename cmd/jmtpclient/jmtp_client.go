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

    //var a = 1
    //b := make([]byte, 1)
    //binary.PutVarint(b, int64(a))
    //reader := bytes.NewReader(b)
    //i, err := binary.ReadVarint(reader)
    //if err != nil {
    //    panic(err)
    //}
    //fmt.Println(i)
    //fmt.Println(byte(1))
    //c := []byte{'1', '2', '3'}
    //r := bytes.NewReader(c)
    //r.ReadByte()
    //r.ReadByte()
    //fmt.Println(r.Len())
    //a := &v1.Connect{}
    //test(a)
    jmtpClientTest()
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


