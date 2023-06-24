package requests

import (
	"net"
	"net/http"
	"time"
)

type CustomTransport struct {
	CTimeout  time.Duration
	RWTimeout time.Duration
	LocalAddr string
}

func (t CustomTransport) Dial(network, addr string) (net.Conn, error) {
	dialer := net.Dialer{
		Timeout:   t.CTimeout,
		DualStack: true,
	}

	if t.LocalAddr != "" {
		dialer.LocalAddr = &net.TCPAddr{IP: net.ParseIP(t.LocalAddr)}
	}

	conn, err := dialer.Dial(network, addr)

	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Now().Add(t.RWTimeout))

	return conn, nil
}

func (t CustomTransport) Transport() *http.Transport {
	return &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		ForceAttemptHTTP2: true,
		Dial:              t.Dial,
	}
}
