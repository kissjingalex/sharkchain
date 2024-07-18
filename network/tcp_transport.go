package network

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

type TCPPeer struct {
	conn     net.Conn
	Outgoing bool
}

func (p *TCPPeer) Send(bs []byte) error {
	_, err := p.conn.Write(bs)
	return err
}

func (p *TCPPeer) readLoop(rpcCh chan RPC) {
	buf := make([]byte, 4096)
	for {
		n, err := p.conn.Read(buf)
		if err == io.EOF {
			//conn closed
			continue
		}
		if err != nil {
			fmt.Printf("read error: %s\n", err)
			// TODO should handle conn closed, such as
			//if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
			//	fmt.Println("The connection was already closed.")
			continue
		}

		msg := buf[:n]
		rpcCh <- RPC{
			From:    p.conn.RemoteAddr(),
			Payload: bytes.NewReader(msg),
		}
	}
}

type TCPTransport struct {
	peerCh     chan *TCPPeer
	listenAddr string
	listener   net.Listener
}

func NewTCPTransport(addr string, peerCh chan *TCPPeer) *TCPTransport {
	return &TCPTransport{
		peerCh:     peerCh,
		listenAddr: addr,
	}
}

func (t *TCPTransport) Start() error {
	ln, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	t.listener = ln
	go t.acceptLoop()

	return nil
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("accept error from %+v\n", conn)
			continue
		}

		fmt.Printf("============ accepted connection from %s\n", conn.RemoteAddr().String())

		t.peerCh <- &TCPPeer{
			conn:     conn,
			Outgoing: true,
		}
	}
}
