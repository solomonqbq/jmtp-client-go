package main

import (
    "fmt"
    "github.com/jmtp/jmtp-client-go"
    "github.com/jmtp/jmtp-client-go/protocol/v1"
    "github.com/jmtp/jmtp-client-go/jmtpclient"
    "time"
)

func main() {
    jmtpClientTest()
}

func test(j jmtp_client_go.JmtpPacket) {
    switch v := j.(type) {
    case *v1.Connect:
        fmt.Println(v)
    }
}

func jmtpClientTest() {
    config := &jmtpclient.Config {
        Url: "jmtp://localhost:20560",
        TimeoutSec: 2,
        HeartbeatSec: 1,
        SerializeType: 1,
        ApplicationId: 999,
        InstanceId: 999,
    }
    client, _ := jmtpclient.NewJmtpClient(config, func(packet jmtp_client_go.JmtpPacket, err error) {
        if err != nil {
            fmt.Println(err)
        } else {
            fmt.Printf("%v", packet)
        }

    })
    err := client.Connect()
    if err != nil {
        fmt.Println(err)
    }
    time.Sleep(time.Duration(30) * time.Second)
    client.Close()
}


