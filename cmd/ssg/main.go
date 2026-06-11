package main

import (
	"fmt"
	"log"
	"time"

	"mehub/internal"
)

func main() {
	start := time.Now()

	distDir := "dist"
	configDir := "internal/templates/contents"
	templatesDir := "internal/templates"
	blogDir := "blog"
	publicDir := "internal/templates/static"

	count, err := internal.RunPipeline(distDir, configDir, templatesDir, blogDir, publicDir)
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	fmt.Printf("✅ Build completed: generated %d posts in %v\n", count, time.Since(start))
}
