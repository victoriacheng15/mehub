package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"mehub/internal/config"
	"mehub/internal/content"
	"mehub/internal/generator"
	"mehub/internal/utils"
)

func main() {
	if err := run("dist", "configs", "blog", "public"); err != nil {
		log.Fatalf("Build failed: %v", err)
	}
}

func run(distDir, configDir, blogDir, publicDir string) error {
	start := time.Now()

	if err := os.RemoveAll(distDir); err != nil {
		return fmt.Errorf("failed to clean dist dir: %w", err)
	}

	cfg, err := config.LoadConfig(configDir)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	gen := generator.New(cfg)

	if err := os.MkdirAll(distDir, 0755); err != nil {
		return fmt.Errorf("failed to create dist dir: %w", err)
	}

	if _, err := os.Stat(publicDir); err == nil {
		if err := utils.CopyDir(publicDir, distDir); err != nil {
			log.Printf("Warning: Failed to copy public assets: %v", err)
		}
	}

	rawPosts, err := content.GetPosts(blogDir)
	if err != nil {
		return fmt.Errorf("failed to load posts: %w", err)
	}

	data := content.ProcessPosts(rawPosts)

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
	if err := gen.GenerateBlogRegistry(distDir, data); err != nil {
		return fmt.Errorf("failed to generate blog registry: %w", err)
	}
	if err := gen.GenerateProjectsRegistry(distDir); err != nil {
		return fmt.Errorf("failed to generate projects registry: %w", err)
	}
	if err := gen.GenerateSkillsRegistry(distDir); err != nil {
		return fmt.Errorf("failed to generate skills registry: %w", err)
	}
	if err := gen.GenerateProfileRegistry(distDir); err != nil {
		return fmt.Errorf("failed to generate profile registry: %w", err)
	}
	if err := gen.GenerateCatalogRegistry(distDir); err != nil {
		return fmt.Errorf("failed to generate catalog registry: %w", err)
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
