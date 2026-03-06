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
	if err := run("dist", "internal/templates/contents", "blog", "internal/templates/static"); err != nil {
		log.Fatalf("Build failed: %v", err)
	}
}

func run(distDir, configDir, blogDir, publicDir string) error {
	start := time.Now()

	if err := os.RemoveAll(distDir); err != nil {
		return fmt.Errorf("failed to clean dist dir: %w", err)
	}

	cfg, err := contents.LoadConfig(configDir)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	gen := internal.New(cfg)

	if err := os.MkdirAll(distDir, 0755); err != nil {
		return fmt.Errorf("failed to create dist dir: %w", err)
	}

	if _, err := os.Stat(publicDir); err == nil {
		if err := utils.CopyDir(publicDir, distDir); err != nil {
			log.Printf("Warning: Failed to copy public assets: %v", err)
		}
	}

	rawPosts, err := post.GetPosts(blogDir)
	if err != nil {
		return fmt.Errorf("failed to load posts: %w", err)
	}

	data := post.ProcessPosts(rawPosts)

	// Site Generation
	if err := gen.GenerateStaticPages(distDir, data); err != nil {
		return fmt.Errorf("failed to generate static pages: %w", err)
	}
	if err := gen.GenerateBlogPagination(distDir, data, 10); err != nil {
		return fmt.Errorf("failed to generate blog pagination: %w", err)
	}
	if err := gen.GenerateTagPages(distDir, data); err != nil {
		return fmt.Errorf("failed to generate tag pages: %w", err)
	}
	if err := gen.GeneratePostPages(distDir, data); err != nil {
		return fmt.Errorf("failed to generate post pages: %w", err)
	}
	if err := gen.GenerateSearchIndex(distDir, data); err != nil {
		return fmt.Errorf("failed to generate search index: %w", err)
	}
	if err := gen.GenerateRegistries(distDir, data); err != nil {
		return fmt.Errorf("failed to generate registries: %w", err)
	}
	if err := gen.GenerateLLMsTxt(distDir); err != nil {
		return fmt.Errorf("failed to generate llms.txt: %w", err)
	}
	if err := gen.GenerateRSS(distDir, data.Posts); err != nil {
		return fmt.Errorf("failed to generate RSS: %w", err)
	}
	if err := gen.GenerateSitemap(distDir, data.Posts); err != nil {
		return fmt.Errorf("failed to generate sitemap: %w", err)
	}

	fmt.Printf("✅ Build completed: generated %d posts in %v\n", len(data.Posts), time.Since(start))
	return nil
}
