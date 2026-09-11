package edge

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// CachedAsset represents an in-memory compressed static file
type CachedAsset struct {
	ContentType  string
	Content      []byte
	GzipContent  []byte
	ETag         string
	LastModified string
}

// StaticServer is a high-performance in-memory asset server in pure Go
type StaticServer struct {
	mu       sync.RWMutex
	baseDir  string
	cache    map[string]*CachedAsset
	fallback []byte
}

// NewStaticServer initializes the in-memory static web server
func NewStaticServer(baseDir string) *StaticServer {
	srv := &StaticServer{
		baseDir: baseDir,
		cache:   make(map[string]*CachedAsset),
	}

	// Preload directory if exists
	if _, err := os.Stat(baseDir); err == nil {
		srv.preloadDirectory(baseDir)
	}

	return srv
}

func (s *StaticServer) preloadDirectory(root string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(root, path)
		urlPath := "/" + filepath.ToSlash(relPath)

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		ext := filepath.Ext(path)
		ctype := mime.TypeByExtension(ext)
		if ctype == "" {
			switch ext {
			case ".js", ".mjs":
				ctype = "application/javascript; charset=utf-8"
			case ".css":
				ctype = "text/css; charset=utf-8"
			case ".json":
				ctype = "application/json"
			case ".html":
				ctype = "text/html; charset=utf-8"
			default:
				ctype = "application/octet-stream"
			}
		}

		// Calculate SHA256 ETag
		h := sha256.Sum256(content)
		etag := `"` + hex.EncodeToString(h[:8]) + `"`

		// Pre-compress with GZIP
		var gzipBuf bytes.Buffer
		gw, _ := gzip.NewWriterLevel(&gzipBuf, gzip.BestCompression)
		gw.Write(content)
		gw.Close()

		asset := &CachedAsset{
			ContentType: ctype,
			Content:     content,
			GzipContent: gzipBuf.Bytes(),
			ETag:        etag,
		}

		s.cache[urlPath] = asset
		if urlPath == "/index.html" {
			s.cache["/"] = asset
			s.fallback = content
		}

		return nil
	})
}

// ServeHTTP serves the cached static file with ETag and Gzip support
func (s *StaticServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	reqPath := r.URL.Path
	if reqPath == "" {
		reqPath = "/"
	}

	asset, exists := s.cache[reqPath]
	if !exists {
		// HTML5 PushState routing fallback to root index.html
		asset = s.cache["/"]
	}

	if asset == nil {
		// Fallback embedded landing page if no static assets found
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Metro-Synapse Edge Gateway</title></head>
<body style="background:#080c14;color:#f1f5f9;font-family:sans-serif;padding:40px;text-align:center;">
<h1 style="color:#00f0ff;">⚡ METRO-SYNAPSE (PRANA-GRID)</h1>
<p style="color:#94a3b8;">Custom Go Edge Gateway & Autonomous Live Server Active</p>
<p style="color:#00ff88;">DNS: RFC 1035 UDP 5353 • Reverse Proxy L7 Active</p>
</body>
</html>`)
		return
	}

	// Check ETag for 304 Not Modified
	if match := r.Header.Get("If-None-Match"); match != "" && match == asset.ETag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", asset.ContentType)
	w.Header().Set("ETag", asset.ETag)
	w.Header().Set("Cache-Control", "public, max-age=3600")

	// Check if client supports GZIP
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && len(asset.GzipContent) < len(asset.Content) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Write(asset.GzipContent)
	} else {
		w.Write(asset.Content)
	}
}
