package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
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

	flag.BoolVar(&showVersion, "version", false, "Show version information")

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
	setupRoutes(r)

	fmt.Println("OpenFMan Server listening on port 3000...")
	if err := r.Run(":3000"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
