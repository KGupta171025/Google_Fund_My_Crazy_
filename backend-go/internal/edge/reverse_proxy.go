package edge

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ReverseProxyConfig defines upstream routing configurations
type ReverseProxyConfig struct {
	GoBackendURL  string
	PythonAIURL   string
	StaticServer  *StaticServer
	CoreHandler   http.Handler
	WSHandler     http.Handler
}

// EdgeGateway is the unified Layer-7 Reverse Proxy & Router
type EdgeGateway struct {
	config    ReverseProxyConfig
	aiUpstream *url.URL
	client    *http.Client
}

// NewEdgeGateway creates a new edge reverse proxy
func NewEdgeGateway(cfg ReverseProxyConfig) *EdgeGateway {
	aiURL, _ := url.Parse(cfg.PythonAIURL)
	return &EdgeGateway{
		config:    cfg,
		aiUpstream: aiURL,
		client: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// ServeHTTP routes incoming traffic to the appropriate microservice or static asset
func (g *EdgeGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// CORS Preflight
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Powered-By", "Metro-Synapse-Go-Edge")

	path := r.URL.Path

	// 1. WebSocket Streaming Endpoint
	if strings.HasPrefix(path, "/ws") {
		if g.config.WSHandler != nil {
			g.config.WSHandler.ServeHTTP(w, r)
			return
		}
	}

	// 2. Core Go Telemetry & Spatial API
	if strings.HasPrefix(path, "/api") {
		if g.config.CoreHandler != nil {
			g.config.CoreHandler.ServeHTTP(w, r)
			return
		}
	}

	// 3. Python AI Microservice Proxy
	if strings.HasPrefix(path, "/ai") || strings.HasPrefix(path, "/predict") || strings.HasPrefix(path, "/gemini") {
		g.proxyToAI(w, r)
		return
	}

	// 4. Static Web Application
	if g.config.StaticServer != nil {
		g.config.StaticServer.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}

	_ = start
}

func (g *EdgeGateway) proxyToAI(w http.ResponseWriter, r *http.Request) {
	if g.aiUpstream == nil {
		http.Error(w, `{"error": "AI upstream not configured"}`, http.StatusBadGateway)
		return
	}

	targetURL := g.aiUpstream.String() + r.URL.RequestURI()
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, `{"error": "Failed to create proxy request"}`, http.StatusInternalServerError)
		return
	}

	for k, vv := range r.Header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	resp, err := g.client.Do(req)
	if err != nil {
		log.Printf("⚠️ [Edge Proxy] Python AI upstream connection failed: %v (serving local fallback)", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "fallback", "engine": "metro-synapse-edge-heuristic"}`))
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
