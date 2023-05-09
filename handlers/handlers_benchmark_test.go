package handlers

import (
	"encoding/json"
	"example/web-service-gin/hll"
	"example/web-service-gin/metrics"
	"fmt"
	"testing"
)

func BenchmarkProcessTextData(b *testing.B) {
	hll.HLLChan = make(chan metrics.Data) // Set the desired buffer size
	hll.HLLDone = make(chan bool)
	go CalculateCardinalityPeriodically(hll.HLLChan, hll.HLLDone)
	data := `up{job="node_exporter",instance="localhost:9105"} 0`
	for i := 0; i < b.N; i++ {
		ProcessOpenMetricsData(data)
	}
	close(hll.HLLChan)
	hll.HLLDone <- true

}

func BenchmarkProcessJSONData(b *testing.B) {
	hll.HLLChan = make(chan metrics.Data) // Set the desired buffer size
	hll.HLLDone = make(chan bool)
	go CalculateCardinalityPeriodically(hll.HLLChan, hll.HLLDone)
	data := `{"metric":{"__name__":"up","job":"node_exporter","instance":"localhost:9104"},"values":[1]}`
	// Create a map to hold the converted data
	jsonData := make(map[string]interface{})

	// Unmarshal the string into the map
	err := json.Unmarshal([]byte(data), &jsonData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for i := 0; i < b.N; i++ {
		ProcessJSONData(jsonData)
	}
	close(hll.HLLChan)
	hll.HLLDone <- true

}
