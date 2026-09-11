package main

import (
	"context"
	"log"
	"metro-synapse/backend/internal/broker"
	"metro-synapse/backend/internal/handlers"
	"metro-synapse/backend/internal/simulation"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("==================================================================")
	log.Println("⚡ PROJECT METRO-SYNAPSE (PRANA-GRID) - URBAN KINETIC NERVOUS SYSTEM")
	log.Println("🚀 Google Gemini 'Fund My Crazy' Moonshot Core Server Starting...")
	log.Println("==================================================================")

	// Initialize event broker and simulation engine
	eventBroker := broker.NewEventBroker(8192)
	cityEngine := simulation.NewCityEngine("city-blr-01", "Bengaluru Cyber-Metropolis", eventBroker)

	// Start live simulation loop
	cityEngine.Start()
	defer cityEngine.Stop()

	// Initialize handlers
	apiHandler := handlers.NewAPIHandler(cityEngine)
	wsHandler := handlers.NewWSHandler(eventBroker)

	mux := http.NewServeMux()

	// REST Routes
	mux.HandleFunc("/api/health", apiHandler.HandleHealth)
	mux.HandleFunc("/api/city/state", apiHandler.HandleGetState)
	mux.HandleFunc("/api/corridor/dispatch", apiHandler.HandleDispatchCorridor)
	mux.HandleFunc("/api/energy/stats", apiHandler.HandleEnergyStats)
	mux.HandleFunc("/api/spatial/query", apiHandler.HandleSpatialQuery)

	// WebSocket Real-time Telemetry Stream
	mux.HandleFunc("/ws/city/stream", wsHandler.HandleWS)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Graceful shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("📡 Metro-Synapse Go Core listening on http://localhost:%s", port)
		log.Printf("⚡ Live Telemetry Stream available at ws://localhost:%s/ws/city/stream", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen failed: %v", err)
		}
	}()

	<-stop
	log.Println("\n🛑 Shutting down Metro-Synapse Go Core gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Metro-Synapse Go Core exited cleanly.")
}
