package internal

import (
	"fmt"
	"log"
	"os"

	"mehub/internal/contents"
	"mehub/internal/post"
	"mehub/internal/utils"
)

// RunPipeline orchestrates the entire site generation flow.
func RunPipeline(distDir, configDir, templatesDir, blogDir, publicDir string) (int, error) {
	// 1. Clean and Prepare Dist Directory
	if err := os.RemoveAll(distDir); err != nil {
		return 0, fmt.Errorf("failed to clean dist dir: %w", err)
	}
	if err := os.MkdirAll(distDir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create dist dir: %w", err)
	}

	// 2. Load Configuration and Initialize Generator
	cfg, err := contents.LoadConfig(configDir)
	if err != nil {
		return 0, fmt.Errorf("failed to load config: %w", err)
	}
	gen := New(cfg, templatesDir)

	// 3. Copy Static Assets
	if _, err := os.Stat(publicDir); err == nil {
		if err := utils.CopyDir(publicDir, distDir); err != nil {
			log.Printf("Warning: Failed to copy public assets: %v", err)
		}
	}

	// 4. Load and Process Content
	rawPosts, err := post.GetPosts(blogDir)
	if err != nil {
		return 0, fmt.Errorf("failed to load posts: %w", err)
	}
	data := post.ProcessPosts(rawPosts)

	// 5. Build Site
	if err := gen.Build(distDir, data); err != nil {
		return 0, fmt.Errorf("build failed: %w", err)
	}

	return len(data.Posts), nil
}
