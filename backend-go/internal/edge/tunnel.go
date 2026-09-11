package edge

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// TunnelBridge allows exposing local services across NAT via raw TCP socket multiplexing
type TunnelBridge struct {
	mu         sync.Mutex
	publicPort int
	localPort  int
	listener   net.Listener
	stopChan   chan struct{}
}

// NewTunnelBridge creates a new scratch TCP tunnel bridge
func NewTunnelBridge(publicPort, localPort int) *TunnelBridge {
	return &TunnelBridge{
		publicPort: publicPort,
		localPort:  localPort,
		stopChan:   make(chan struct{}),
	}
}

// Start launches the TCP tunnel listener
func (tb *TunnelBridge) Start() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", tb.publicPort))
	if err != nil {
		return fmt.Errorf("failed to bind tunnel public port %d: %w", tb.publicPort, err)
	}
	tb.listener = ln

	log.Printf("🚇 [Scratch Tunnel] TCP Multiplexed Tunnel Bridge active on :%d -> :%d", tb.publicPort, tb.localPort)

	go func() {
		for {
			select {
			case <-tb.stopChan:
				return
			default:
				publicConn, err := ln.Accept()
				if err != nil {
					return
				}
				go tb.handleTunnelConnection(publicConn)
			}
		}
	}()

	return nil
}

// Stop terminates the tunnel
func (tb *TunnelBridge) Stop() {
	close(tb.stopChan)
	if tb.listener != nil {
		tb.listener.Close()
	}
}

func (tb *TunnelBridge) handleTunnelConnection(publicConn net.Conn) {
	defer publicConn.Close()

	// Connect to local service
	localConn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", tb.localPort), 2*time.Second)
	if err != nil {
		log.Printf("⚠️ Tunnel failed to connect to local target port %d: %v", tb.localPort, err)
		return
	}
	defer localConn.Close()

	// Bidirectional full-duplex pipe
	errChan := make(chan error, 2)
	go func() {
		_, err := io.Copy(localConn, publicConn)
		errChan <- err
	}()
	go func() {
		_, err := io.Copy(publicConn, localConn)
		errChan <- err
	}()

	<-errChan
}
