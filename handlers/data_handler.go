package handlers

import (
	"encoding/json"
	"example/web-service-gin/hll"
	"example/web-service-gin/metrics"
	"fmt"
	"sync"
	"time"
)

var mu sync.Mutex
var metricsData []metrics.Data
var done chan bool

func ProcessJSONData(jsonData map[string]interface{}) error {
	// Parse the JSON data
	metricBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	var data metrics.Data
	// Convert it into Data struct
	err = json.Unmarshal(metricBytes, &data)
	if err != nil {
		return err
	}

	// Append the metric data to the global metrics data
	SaveMetrics(data)

	return nil
}

func ProcessOpenMetricsData(data string) error {
	// Extract the metric name and values
	dataObj, err := metrics.ProcessDataText(data)
	if err != nil {
		return err
	}

	// Append the metric data to the global metrics data
	SaveMetrics(*dataObj)
	return nil
}

func SaveMetrics(d metrics.Data) {

	// Add the new metric data received to the Channel for processing
	hll.HLLChan <- d
	// fmt.Printf("Added element: %v\n", d)

	// TODO : Flush the metrics to DB
	// Temporarily save it to an array to flush
	//mu.Lock()
	//metricsData = append(metricsData, d)
	//mu.Unlock()

}

func CalculateCardinalityPeriodically(dataChan <-chan metrics.Data, done chan bool) {
	cardinality := hll.Init()
	var mu sync.Mutex // Mutex for synchronizing access to cardinality

	// Start a goroutine for periodic processing
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Time's up, process the accumulated elements
				mu.Lock()
				count := cardinality.Count
				mu.Unlock()

				currentTime := time.Now()
				fmt.Printf("%v Cardinality: %v\n", currentTime.Format("2006-01-02 15:04:05"), count)

			case <-done:
				// Processing finished, exit the goroutine
				return
			}
		}
	}()

	// Process incoming elements
	for element := range dataChan {
		// Process the element
		dic, err := metrics.StructToMap(element.Metric)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		mu.Lock()
		cardinality.Add(dic)
		mu.Unlock()
	}
	// Wait for processing to finish
	defer func() {
		done <- true
	}()

	<-done
}
