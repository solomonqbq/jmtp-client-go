package jmtp_client_go

type JMTPClient interface {
    Connect() error
    ConnectByUrl(url string) error
    Close() error

}
