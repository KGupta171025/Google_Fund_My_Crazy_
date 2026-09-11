package edge

import (
	"bytes"
	"encoding/binary"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDomainEncoding(t *testing.T) {
	var buf bytes.Buffer
	encodeDomainName(&buf, "prana-grid.live")

	expected := []byte{10, 'p', 'r', 'a', 'n', 'a', '-', 'g', 'r', 'i', 'd', 4, 'l', 'i', 'v', 'e', 0}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Fatalf("Encoded domain mismatch: got %v, expected %v", buf.Bytes(), expected)
	}
}

func TestDNSServerResolution(t *testing.T) {
	dnsPort := 15353
	server := NewDNSServer(dnsPort)
	server.RegisterRecord("test.local", net.ParseIP("192.168.1.100"))

	err := server.Start()
	if err != nil {
		t.Fatalf("Failed to start DNS server: %v", err)
	}
	defer server.Stop()

	// Wait for socket to bind
	time.Sleep(50 * time.Millisecond)

	// Send raw DNS query packet
	conn, err := net.Dial("udp", "127.0.0.1:15353")
	if err != nil {
		t.Fatalf("Failed to connect to DNS server: %v", err)
	}
	defer conn.Close()

	// Build raw DNS Query packet
	var qBuf bytes.Buffer
	header := DNSHeader{
		ID:      0x1234,
		Flags:   0x0100, // Standard query
		QDCount: 1,
	}
	binary.Write(&qBuf, binary.BigEndian, header)
	encodeDomainName(&qBuf, "test.local")
	binary.Write(&qBuf, binary.BigEndian, uint16(1)) // Type A
	binary.Write(&qBuf, binary.BigEndian, uint16(1)) // Class IN

	_, err = conn.Write(qBuf.Bytes())
	if err != nil {
		t.Fatalf("Failed to send DNS query: %v", err)
	}

	resp := make([]byte, 512)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(resp)
	if err != nil {
		t.Fatalf("Failed to read DNS response: %v", err)
	}

	if n < 12 {
		t.Fatalf("DNS response too short: %d bytes", n)
	}

	var respHeader DNSHeader
	binary.Read(bytes.NewReader(resp[:12]), binary.BigEndian, &respHeader)

	if respHeader.ID != 0x1234 {
		t.Fatalf("Response ID mismatch: got %x, expected 1234", respHeader.ID)
	}

	if respHeader.ANCount < 1 {
		t.Fatalf("Expected at least 1 answer record, got %d", respHeader.ANCount)
	}
}

func TestStaticServerETag(t *testing.T) {
	srv := NewStaticServer("./non-existent-dir")
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}
