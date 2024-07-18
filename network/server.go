package network

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-kit/log"
	"net"
	"os"
	"sharkchain/core"
	"sharkchain/crypto"
	"time"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	Transports []Transport
	BlockTime  time.Duration
	ID         string
	Logger     log.Logger
	PrivateKey *crypto.PrivateKey

	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
}

type Server struct {
	ServerOpts
	memPool     *TxPool
	isValidator bool // depends on weather has private key
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "addr", opts.ID)
	}

	s := &Server{
		ServerOpts: opts,
		memPool:    NewTxPool(0),
		rpcCh:      make(chan RPC),
		quitCh:     make(chan struct{}),
	}

	// If we dont got any processor from the server options, we going to use
	// the server as default.
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	return s
}

func (s *Server) Start() {
	s.initTransports()

	ticker := time.NewTicker(s.BlockTime)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("RPC error", err)
				continue
			}

			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				if !errors.Is(err, core.ErrBlockKnown) {
					s.Logger.Log("error", err)
				}
			}
		case <-s.quitCh:
			break free
		case <-ticker.C:
			if s.isValidator {
				s.createNewBlock()
			}
		default:
		}
	}

	fmt.Println("Server stopped")
}

func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	fmt.Printf("process message from %s\n", msg.From.String())

	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	case *core.Block:
		return s.processBlock(t)
	case *GetStatusMessage:
		return s.processGetStatusMessage(msg.From, t)
	case *StatusMessage:
		return s.processStatusMessage(msg.From, t)
	case *GetBlocksMessage:
		return s.processGetBlocksMessage(msg.From, t)
	case *BlocksMessage:
		return s.processBlocksMessage(msg.From, t)
	}

	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {
	fmt.Printf("processing transaction from %s\n", tx.From.Address().String())

	if err := tx.Verify(); err != nil {
		return err
	}

	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Contains(hash) {
		return nil
	}

	// TODO need to broadcast

	s.memPool.Add(tx)
	return nil
}

func (s *Server) broadcast(payload []byte) error {
	//s.mu.RLock()
	//defer s.mu.RUnlock()
	//for netAddr, peer := range s.peerMap {
	//	if err := peer.Send(payload); err != nil {
	//		fmt.Printf("peer send error => addr %s [err: %s]\n", netAddr, err)
	//	}
	//}

	return nil
}

func (s *Server) broadcastBlock(b *core.Block) error {
	buf := &bytes.Buffer{}
	if err := b.Encode(core.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeBlock, buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)
	}
}

func (s *Server) processBlock(t *core.Block) error {
	return nil
}

func (s *Server) processGetStatusMessage(from net.Addr, t *GetStatusMessage) error {
	return nil
}

func (s *Server) processStatusMessage(from net.Addr, t *StatusMessage) error {
	return nil
}

func (s *Server) processGetBlocksMessage(from net.Addr, t *GetBlocksMessage) error {
	return nil
}

func (s *Server) processBlocksMessage(from net.Addr, t *BlocksMessage) error {
	return nil
}
