package main

import (
	"github.com/sridhargude/levitate/handlers"
	"github.com/sridhargude/levitate/hll"
	"github.com/sridhargude/levitate/metrics"
	"github.com/gin-gonic/gin"
)

// Index page
func indexPage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to MetricsServer page",
	})
}

// Process JSON and Open Metrics
func handleMetrics(c *gin.Context) {
	// Determine the content type of the request
	contentType := c.GetHeader("Content-Type")

	// Handle JSON data
	if contentType == "application/json" {
		// Parse the JSON data
		var jsonData map[string]interface{}
		if err := c.BindJSON(&jsonData); err != nil {
			c.AbortWithError(400, err)
			return
		}

		// Process the JSON data
		if err := handlers.ProcessJSONData(jsonData); err != nil {
			c.AbortWithError(500, err)
			return
		}

		// Send a response
		c.JSON(200, gin.H{
			"message": "JSON data processed successfully",
		})
	}

	// Handle OpenMetrics data
	if contentType == "text/plain" {
		// Read the raw data from the request body
		rawData, err := c.GetRawData()
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		// Process the OpenMetrics data
		if err := handlers.ProcessOpenMetricsData(string(rawData)); err != nil {
			c.AbortWithError(500, err)
			return
		}

		// Send a response
		c.JSON(200, gin.H{
			"message": "OpenMetrics data processed successfully",
		})
	}
}

func main() {
	// Run the Cardinality Metrics
	hll.HLLChan = make(chan metrics.Data, 0)
	hll.HLLDone = make(chan bool)
	// TODO: Configurable time period to emit the Cardinality
	go handlers.CalculateCardinalityPeriodically(hll.HLLChan, hll.HLLDone)

	// Create a new Gin router
	r := gin.Default()

	// Route to handle JSON and OpenMetrics data
	r.POST("/metrics/json", handleMetrics)
	// Index Page
	r.GET("/", indexPage)

	// Start the server
	r.Run("localhost:8081")
}
