package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	version      string
	commit       string
	buildTime    string
	isDevVersion string // Build-time variable to indicate if it's a dev version
)

func main() {
	// Define command-line flags
	var showVersion bool
	var port int

	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.IntVar(&port, "port", 3000, "Port to run the server on")

	// Parse command-line flags
	flag.Parse()

	// Show version if requested
	if showVersion {
		if isDevVersion == "true" {
			fmt.Printf("OpenFMan development version (%s) %s\n", commit, buildTime)
		} else {
			fmt.Printf("OpenFMan v%s (%s)\n", version, commit)
		}
		return
	}

	loadConfig()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	setupRoutes(r)

	fmt.Printf("OpenFMan Server listening on port %d...\n", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
