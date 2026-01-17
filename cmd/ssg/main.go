package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"mehub-ssg/internal/config"
	"mehub-ssg/internal/content"
	"mehub-ssg/internal/generator"
	"mehub-ssg/internal/utils"
)

func main() {
	start := time.Now()

	distDir := "dist"
	if err := os.RemoveAll(distDir); err != nil {
		log.Fatalf("Failed to clean dist dir: %v", err)
	}

	cfg, err := config.LoadConfig("configs")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	gen := generator.New(cfg)

	if err := os.MkdirAll(distDir, 0755); err != nil {
		log.Fatalf("Failed to create dist dir: %v", err)
	}

	if err := utils.CopyDir("public", distDir); err != nil {
		log.Printf("Warning: Failed to copy public assets: %v", err)
	}

	rawPosts, err := content.GetPosts("blog")
	if err != nil {
		log.Fatalf("Failed to load posts: %v", err)
	}

	data := content.ProcessPosts(rawPosts)

	// Site Generation
	gen.GenerateStaticPages(distDir, data)
	gen.GenerateBlogPagination(distDir, data, 10)
	gen.GenerateTagPages(distDir, data)
	gen.GeneratePostPages(distDir, data)
	gen.GenerateSearchIndex(distDir, data)
	gen.GenerateRSS(distDir, data.Posts)
	gen.GenerateSitemap(distDir, data.Posts)

	fmt.Printf("✅ Build completed: generated %d posts in %v\n", len(data.Posts), time.Since(start))
}
