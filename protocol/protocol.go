package protocol

const (
    Reserved = 0x0
    Conntect = 0x1
    ConnectAck = 0x2
    Ping = 0x3
    Pong = 0x4
    Disconnect = 0x5
    Report = 0x6
    ReportAck = 0x7
    Command = 0x8
    CommandAck = 0x9
)

type Protocol struct {
    header byte
}