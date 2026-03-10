package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"mehub/internal"
	"mehub/internal/contents"
	"mehub/internal/post"
	"mehub/internal/utils"
)

func main() {
	start := time.Now()

	distDir := "dist"
	configDir := "internal/templates/contents"
	blogDir := "blog"
	publicDir := "internal/templates/static"

	// 1. Clean and Prepare Dist Directory
	if err := os.RemoveAll(distDir); err != nil {
		log.Fatalf("failed to clean dist dir: %v", err)
	}
	if err := os.MkdirAll(distDir, 0755); err != nil {
		log.Fatalf("failed to create dist dir: %v", err)
	}

	// 2. Load Configuration and Initialize Generator
	cfg, err := contents.LoadConfig(configDir)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	gen := internal.New(cfg)

	// 3. Copy Static Assets
	if _, err := os.Stat(publicDir); err == nil {
		if err := utils.CopyDir(publicDir, distDir); err != nil {
			log.Printf("Warning: Failed to copy public assets: %v", err)
		}
	}

	// 4. Load and Process Content
	rawPosts, err := post.GetPosts(blogDir)
	if err != nil {
		log.Fatalf("failed to load posts: %v", err)
	}
	data := post.ProcessPosts(rawPosts)

	// 5. Build Site
	if err := gen.Build(distDir, data); err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	fmt.Printf("✅ Build completed: generated %d posts in %v\n", len(data.Posts), time.Since(start))
}
