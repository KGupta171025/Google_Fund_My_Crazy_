package edge

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

// DNSHeader represents RFC 1035 standard 12-byte DNS message header
type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16 // Question count
	ANCount uint16 // Answer count
	NSCount uint16 // Authority count
	ARCount uint16 // Additional count
}

// DNSQuestion represents a parsed query question
type DNSQuestion struct {
	Name  string
	Type  uint16 // Type A = 1
	Class uint16 // Class IN = 1
}

// DNSServer is a scratch-built UDP DNS nameserver
type DNSServer struct {
	mu         sync.RWMutex
	port       int
	records    map[string]net.IP // domain -> IPv4 address
	conn       *net.UDPConn
	stopChan   chan struct{}
	queryCount uint64
}

// NewDNSServer initializes a custom RFC 1035 DNS server
func NewDNSServer(port int) *DNSServer {
	if port <= 0 {
		port = 5353 // Default to 5353 or 53
	}
	s := &DNSServer{
		port:     port,
		records:  make(map[string]net.IP),
		stopChan: make(chan struct{}),
	}

	// Register default custom domains
	s.RegisterRecord("prana-grid.live", net.ParseIP("127.0.0.1"))
	s.RegisterRecord("metro-synapse.local", net.ParseIP("127.0.0.1"))
	s.RegisterRecord("fundmycrazy.local", net.ParseIP("127.0.0.1"))
	s.RegisterRecord("localhost", net.ParseIP("127.0.0.1"))

	return s
}

// RegisterRecord registers an A record for a custom domain
func (s *DNSServer) RegisterRecord(domain string, ip net.IP) {
	s.mu.Lock()
	defer s.mu.Unlock()
	domain = strings.ToLower(strings.TrimSuffix(domain, "."))
	s.records[domain] = ip.To4()
}

// Start launches the UDP DNS server loop
func (s *DNSServer) Start() error {
	addr := net.UDPAddr{
		Port: s.port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return fmt.Errorf("failed to bind UDP DNS server on port %d: %w", s.port, err)
	}
	s.conn = conn

	log.Printf("🌐 [Scratch DNS] Custom RFC 1035 DNS Server listening on UDP :%d", s.port)

	go s.serveLoop()
	return nil
}

// Stop closes the DNS server
func (s *DNSServer) Stop() {
	close(s.stopChan)
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *DNSServer) serveLoop() {
	buf := make([]byte, 512)
	for {
		select {
		case <-s.stopChan:
			return
		default:
			n, clientAddr, err := s.conn.ReadFromUDP(buf)
			if err != nil {
				return
			}
			go s.handleQuery(buf[:n], clientAddr)
		}
	}
}

func (s *DNSServer) handleQuery(reqBytes []byte, clientAddr *net.UDPAddr) {
	if len(reqBytes) < 12 {
		return
	}

	var header DNSHeader
	reader := bytes.NewReader(reqBytes)
	binary.Read(reader, binary.BigEndian, &header)

	// Parse Questions
	var questions []DNSQuestion
	for i := 0; i < int(header.QDCount); i++ {
		q, err := parseQuestion(reader)
		if err != nil {
			return
		}
		questions = append(questions, q)
	}

	// Prepare Response Header
	// Flags: QR=1 (response), AA=1 (authoritative), RCODE=0 (no error) -> 0x8180
	respHeader := DNSHeader{
		ID:      header.ID,
		Flags:   0x8180,
		QDCount: header.QDCount,
		ANCount: 0,
		NSCount: 0,
		ARCount: 0,
	}

	var answersBuf bytes.Buffer
	for _, q := range questions {
		s.mu.RLock()
		ip, exists := s.records[strings.ToLower(q.Name)]
		if !exists {
			// Fallback wildcard match
			ip = s.records["prana-grid.live"]
		}
		s.mu.RUnlock()

		if ip != nil && q.Type == 1 { // Type A record query
			respHeader.ANCount++
			// Encode Question Name into Answer
			encodeDomainName(&answersBuf, q.Name)
			binary.Write(&answersBuf, binary.BigEndian, uint16(1))      // Type A
			binary.Write(&answersBuf, binary.BigEndian, uint16(1))      // Class IN
			binary.Write(&answersBuf, binary.BigEndian, uint32(60))     // TTL = 60s
			binary.Write(&answersBuf, binary.BigEndian, uint16(4))      // Data length = 4 bytes (IPv4)
			answersBuf.Write(ip.To4())
		}
	}

	// Build full response packet
	var respBuf bytes.Buffer
	binary.Write(&respBuf, binary.BigEndian, respHeader)
	// Echo back question section (RFC 1035 requirement)
	respBuf.Write(reqBytes[12:reader.Size()-int64(reader.Len())])
	// Append answer resource records
	respBuf.Write(answersBuf.Bytes())

	s.conn.WriteToUDP(respBuf.Bytes(), clientAddr)
}

func parseQuestion(r *bytes.Reader) (DNSQuestion, error) {
	var q DNSQuestion
	var domainParts []string

	for {
		lenByte, err := r.ReadByte()
		if err != nil {
			return q, err
		}
		if lenByte == 0 {
			break
		}
		part := make([]byte, lenByte)
		_, err = r.Read(part)
		if err != nil {
			return q, err
		}
		domainParts = append(domainParts, string(part))
	}

	q.Name = strings.Join(domainParts, ".")
	binary.Read(r, binary.BigEndian, &q.Type)
	binary.Read(r, binary.BigEndian, &q.Class)
	return q, nil
}

func encodeDomainName(buf *bytes.Buffer, name string) {
	parts := strings.Split(name, ".")
	for _, p := range parts {
		if len(p) > 0 {
			buf.WriteByte(byte(len(p)))
			buf.WriteString(p)
		}
	}
	buf.WriteByte(0)
}
