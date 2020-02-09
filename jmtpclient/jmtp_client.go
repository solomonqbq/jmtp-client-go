package jmtpclient

import (
    jc "github.com/jmtp/jmtp-client-go"
    urlParser "net/url"
    "strings"
    "errors"
    "strconv"
    "fmt"
    "net"
    "bufio"
    "github.com/jmtp/jmtp-client-go/protocol"
    "github.com/jmtp/jmtp-client-go/protocol/v1"
    "time"
)

const jmtpProtocol = "jmtp"

type Config struct {
    Url string
    TimeoutSec  int
    HeartbeatSec int
    SerializeType int
    ApplicationId int
    InstanceId int
}

type jmtpClient struct {
    url string
    jmtpUrl *jmtpUrl
    connection  *net.TCPConn
    hawkServer  *net.TCPAddr
    callBack    jc.Callback
    connectSuccess  bool
    packetChain chan jc.JmtpPacket
    errorChain  chan error
    clientConfig *Config
    shutdownSignal  chan bool
}

func NewJmtpClient(config *Config, callback jc.Callback) (*jmtpClient, error) {
    urlParser, err := NewUrlParser(config.Url)
    if err != nil {
        return nil, err
    }
    return &jmtpClient {
        url: config.Url,
        callBack: callback,
        packetChain: make(chan jc.JmtpPacket),
        errorChain: make(chan error),
        clientConfig: config,
        shutdownSignal: make(chan bool, 1),
        jmtpUrl: urlParser,
    }, nil
}

func (c *jmtpClient) Connect() error {
    if c.url == "" {
        return errors.New(fmt.Sprintf("invalidate jmtp connect url: %s", c.url))
    }
    hawkServer, err := net.ResolveTCPAddr("tcp", c.jmtpUrl.GetHost())
    if err != nil {
        return err
    }
    conn, err := net.DialTCP("tcp", nil, hawkServer)
    if err != nil {
        return err
    }
    c.hawkServer = hawkServer
    c.connection = conn
    err = c.sendConnectReq()
    go c.receivePackets()
    go c.chanListener()
    go c.ping()

    return err
}

func (c *jmtpClient) SetUrl(url string) {
    c.url = url
}

func (c *jmtpClient) Close() error {
    if c.connection != nil {
        c.disconnectReq()
        c.shutdownSignal <- true
        return c.connection.Close()
    }
    return nil
}

func (c *jmtpClient) Destroy() error {
    panic("implement me")
}

func (c *jmtpClient) SendPacket(packet jc.JmtpPacket) (int, error) {
    if c.connection != nil {
        out, err := protocol.PacketEncoder(packet)
        if err != nil {
            return 0, err
        }
        return c.connection.Write(out)
    } else {
        return 0, errors.New("connection has been closed")
    }
}

func (c *jmtpClient) Reset() error {
    return nil
}

func (c *jmtpClient) receivePackets() {
    reader := bufio.NewReader(c.connection)
    err := protocol.PacketDecoder(reader, c.packetChain, c.errorChain)
    if err != nil {
        fmt.Println(err)
        c.connection.Close()
    }
}

func (c *jmtpClient) sendConnectReq() error{
    option := &jc.ConnectOption{
        HeartbeatSeconds: int16(c.clientConfig.HeartbeatSec),
        SerializeType: int16(c.clientConfig.SerializeType),
        ApplicationId: c.clientConfig.ApplicationId,
        InstanceId: c.clientConfig.InstanceId,
    }
    connectPack := v1.JMTPV1ProtocolDefineInstance.ConnectPacket(option)
    _, err := c.SendPacket(connectPack)
    return err
}

func (c *jmtpClient) ping() {
    ticker := time.NewTicker(
        time.Duration(c.clientConfig.HeartbeatSec) * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            fmt.Println("Send ping packet")
            _, err := c.SendPacket(v1.JMTPV1ProtocolDefineInstance.PingPacket())
            if err != nil {
                fmt.Printf("err: %v\n", err)
                ticker.Stop()
            }
        case <-c.shutdownSignal:
            break
        }
    }
}

func (c * jmtpClient) disconnectReq() error{
    disconnect := &v1.Disconnect{
        RedirectUrl: c.jmtpUrl.ToUrlString(),
    }
    _, err := c.SendPacket(disconnect)
    return err
}

func (c *jmtpClient) chanListener() {
    for {
        select {
        case packet := <- c.packetChain:
            switch pack := packet.(type) {
            case *v1.ConnectAck:
                fmt.Printf("%v", pack)
                if pack.Code != 0 {
                    err := errors.New("connect to server error, connect has been closed")
                    c.callBack(pack, err)
                    c.Close()
                }
            case *v1.Command:
                // do command
            case *v1.Pong:
                fmt.Println("Get pont packet")
            default:
                c.callBack(packet, nil)
            }
        case err := <- c.errorChain:
            c.callBack(nil, err)
        case <- c.shutdownSignal:
            break
        }
    }
}

type jmtpUrl struct {
    hostname string
    host    string
    port    int
}

func (j *jmtpUrl) parseUrl(urlString string) error {
    url, err := urlParser.Parse(urlString)
    if err != nil {
        return err
    }
    if strings.ToLower(url.Scheme) != jmtpProtocol {
        return errors.New("invalidate protocol name")
    }
    j.hostname = url.Hostname()
    j.host = url.Host
    j.port, err = strconv.Atoi(url.Port())
    if err != nil {
        return err
    }
    return nil
}

func (j *jmtpUrl) GetHost() string{
    return j.host
}

func (j *jmtpUrl) ToUrlString() string{
    return fmt.Sprintf("%s://%s", jmtpProtocol, j.host)
}

func (j *jmtpUrl) GetPort() int {
    return j.port
}

func (j *jmtpUrl) GetHostname() string {
    return j.hostname
}

func NewUrlParser(urlStr string) (*jmtpUrl, error) {
    jmtpUrl := &jmtpUrl{}
    return jmtpUrl, jmtpUrl.parseUrl(urlStr)
}
