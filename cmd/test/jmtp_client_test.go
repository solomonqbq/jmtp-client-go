package test

import (
    "testing"
    "github.com/jmtp/jmtp-client-go/jmtpclient"
    "github.com/jmtp/jmtp-client-go"
    "fmt"
    "github.com/jmtp/jmtp-client-go/protocol/v1"
    "time"
)

func TestJmtpClient(t *testing.T) {
    config := &jmtpclient.Config {
        Url: "jmtp://localhost:20560",
        TimeoutSec: 2,
        HeartbeatSec: 10,
        SerializeType: 1,
        ApplicationId: 999,
        InstanceId: 999,
    }
    client, _ := jmtpclient.NewJmtpClient(config, func(packet jmtp_client_go.JmtpPacket, err error) {
        if err != nil {
            t.Error(err)
        } else {
            fmt.Printf("%v", packet)
        }

    })
    err := client.Connect()
    if err != nil {
        t.Error(err)
    }
    for i := 0; i < 5; i++ {
        report := &v1.Report{}
        report.PacketId = []byte{0x01}
        report.ReportType = 1
        report.Payload = []byte{byte(i)}
        if _, err := client.SendPacket(report); err != nil {
            t.Error(err)
            i--
        }
        if i == 3 {
            client.Reconnect()
        }
        time.Sleep(time.Duration(3000) * time.Millisecond)

    }
    client.Close()
}
