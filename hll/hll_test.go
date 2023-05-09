package hll

import (
	"fmt"
	"testing"
)

func TestHLL(t *testing.T) {
	result := int(5)
	// Sample dictionaries
	dictionaries := []map[string]interface{}{
		{"__name__": "up", "job": "node_exporter", "instance": "localhost:9100"},
		{"__name__": "up", "job": "node_exporter", "instance": "localhost:9100"},
		{"__name__": "down", "job": "node_exporter", "instance": "localhost:9100"},
		{"__name__": "down", "job": "node_exporter", "instance": "localhost:9100"},
		{"__name__": "up", "job": "node_exporter", "instance": "localhost:9200"},
		{"__name__": "up", "job": "node_exporter", "instance": "localhost:9300"},
		{"__name__": "up", "job": "node_exporter", "instance": "localhost:9300"},
		{"job": "node_exporter", "__name__": "up", "instance": "localhost:9300"},
		{"instance": "localhost:9300", "job": "node_exporter", "__name__": "up"},
		{"__name__": "up1", "job": "node_exporter", "instance": "localhost:9300"},
		// Add more dictionaries as needed
	}

	//Init
	hllData := Init()

	// Calculate the distinct count of the dictionaries using HyperLogLog
	distinctCount, err := hllData.CalculateDistinctMetrics(dictionaries)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Distinct Count:", distinctCount)
	if fmt.Sprintf("%d", result) != fmt.Sprintf("%d", distinctCount) {
		t.Errorf("Expected %d but got %d", result, distinctCount)
	}
}

func TestHLLAddCount(t *testing.T) {
	result := int(5)
	dictionaries := []map[string]interface{}{
		{"__name__": "up", "job": "node_exporter", "instance": "localhost:9100"},
		{"__name__": "up", "job": "node_exporter", "instance": "localhost:9100"},
		{"__name__": "down", "job": "node_exporter", "instance": "localhost:9100"},
		{"__name__": "down", "job": "node_exporter", "instance": "localhost:9100"},
		{"__name__": "up", "job": "node_exporter", "instance": "localhost:9200"},
		{"__name__": "up", "job": "node_exporter", "instance": "localhost:9300"},
		{"__name__": "up", "job": "node_exporter", "instance": "localhost:9300"},
		{"job": "node_exporter", "__name__": "up", "instance": "localhost:9300"},
		{"instance": "localhost:9300", "job": "node_exporter", "__name__": "up"},
		{"__name__": "up1", "job": "node_exporter", "instance": "localhost:9300"},
		// Add more dictionaries as needed
	}

	//Init
	hllData := Init()

	for _, metric := range dictionaries {
		hllData.Add(metric)
	}
	if fmt.Sprintf("%d", result) != fmt.Sprintf("%d", hllData.Count) {
		t.Errorf("Expected %d but got %d", result, hllData.Count)
	}
}
