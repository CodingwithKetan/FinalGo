package main

import (
	"encoding/json"
	"log"
	"os"

	"go_discovery_plugin/discovery"
	"go_discovery_plugin/metrics"
	"go_discovery_plugin/utils"
)

type RequestType string

const (
	Discovery RequestType = "discovery"
	Metrics   RequestType = "metrics"
)

type GeneralRequest struct {
	Type      RequestType       `json:"type"`
	Discovery discovery.Request `json:"discovery"`
	Metrics   metrics.Request   `json:"metrics"`
}

func main() {
	utils.InitLogger()
	defer utils.CloseLogger()

	log.Println("Plugin started. Waiting for input...")

	var request GeneralRequest
	decoder := json.NewDecoder(os.Stdin)
	if err := decoder.Decode(&request); err != nil {
		log.Fatalf("Failed to decode input: %v", err)
	}

	switch request.Type {
	case Discovery:
		result := discovery.RunDiscovery(request.Discovery)
		utils.OutputResponse(result)
	case Metrics:
		result := metrics.RunMetricsCollection(request.Metrics)
		utils.OutputResponse(result)
	default:
		log.Fatalf("Invalid request type: %s", request.Type)
	}
}
