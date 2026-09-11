package main

import (
	"context"
	"log"
	"metro-synapse/backend/internal/broker"
	"metro-synapse/backend/internal/edge"
	"metro-synapse/backend/internal/handlers"
	"metro-synapse/backend/internal/simulation"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	dnsPort := 5353 // Custom scratch UDP DNS port
	if os.Getenv("DNS_PORT") != "" {
		dnsPort = 53
	}

	log.Println("==========================================================================")
	log.Println("⚡ PROJECT METRO-SYNAPSE (PRANA-GRID) - AUTONOMOUS GO EDGE DEPLOYMENT ENGINE")
	log.Println("🚀 Scratch-Built RFC 1035 DNS Server • L7 Reverse Proxy • In-Memory Web Core")
	log.Println("==========================================================================")

	// 1. Initialize Custom Scratch RFC 1035 DNS Server
	dnsServer := edge.NewDNSServer(dnsPort)
	if err := dnsServer.Start(); err != nil {
		log.Printf("⚠️ DNS bind warning (port %d): %v (continuing with HTTP gateway)", dnsPort, err)
	}
	defer dnsServer.Stop()

	// 2. Initialize Core Go Simulation & Pub/Sub Broker
	eventBroker := broker.NewEventBroker(8192)
	cityEngine := simulation.NewCityEngine("city-blr-01", "Bengaluru Cyber-Metropolis", eventBroker)
	cityEngine.Start()
	defer cityEngine.Stop()

	// 3. Initialize Microservice Handlers
	apiHandler := handlers.NewAPIHandler(cityEngine)
	wsHandler := handlers.NewWSHandler(eventBroker)

	coreMux := http.NewServeMux()
	coreMux.HandleFunc("/api/health", apiHandler.HandleHealth)
	coreMux.HandleFunc("/api/city/state", apiHandler.HandleGetState)
	coreMux.HandleFunc("/api/corridor/dispatch", apiHandler.HandleDispatchCorridor)
	coreMux.HandleFunc("/api/energy/stats", apiHandler.HandleEnergyStats)
	coreMux.HandleFunc("/api/spatial/query", apiHandler.HandleSpatialQuery)

	wsMux := http.NewServeMux()
	wsMux.HandleFunc("/ws/city/stream", wsHandler.HandleWS)

	// 4. Initialize In-Memory Static Web Server
	staticDir := filepath.Join("..", "frontend-next", "out")
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		staticDir = filepath.Join("..", "frontend-next", "public")
	}
	staticServer := edge.NewStaticServer(staticDir)

	// 5. Initialize Layer-7 Edge Gateway / Reverse Proxy
	gateway := edge.NewEdgeGateway(edge.ReverseProxyConfig{
		GoBackendURL: "http://127.0.0.1:" + httpPort,
		PythonAIURL:  "http://127.0.0.1:8000",
		StaticServer: staticServer,
		CoreHandler:  coreMux,
		WSHandler:    wsMux,
	})

	// 6. Start Unified HTTP Server
	server := &http.Server{
		Addr:         ":" + httpPort,
		Handler:      gateway,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🌐 [Edge Server] Live HTTP Gateway listening on http://0.0.0.0:%s", httpPort)
		log.Printf("🌐 [Custom Domain] http://prana-grid.live:%s", httpPort)
		log.Printf("🌐 [Custom Domain] http://metro-synapse.local:%s", httpPort)
		log.Printf("⚡ [WebSocket] ws://localhost:%s/ws/city/stream", httpPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen failed: %v", err)
		}
	}()

	<-stop
	log.Println("\n🛑 Gracefully stopping Metro-Synapse Edge Deployment Engine...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	log.Println("✅ Metro-Synapse Edge Engine exited cleanly.")
}
