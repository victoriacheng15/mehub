package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunPipeline(t *testing.T) {
	tmpDir := t.TempDir()

	distDir := filepath.Join(tmpDir, "dist")
	configDir := filepath.Join(tmpDir, "contents")
	templatesDir := filepath.Join(tmpDir, "templates")
	blogDir := filepath.Join(tmpDir, "blog")
	publicDir := filepath.Join(tmpDir, "static")

	// 1. Create dummy configs
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	configYAML := `
landing:
  title: "Integration Test Site"
  url: "https://example.com/"
about:
  paragraphs:
    - "hello integration"
navigation:
  header: []
skills: []
socials: []
`
	projectsYAML := `projects: []`

	configs := map[string]string{
		"config.yaml":   configYAML,
		"projects.yaml": projectsYAML,
	}
	for file, content := range configs {
		if err := os.WriteFile(filepath.Join(configDir, file), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// 2. Create dummy templates
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatal(err)
	}
	baseHTML := `{{ define "base.html" }}<html><body>{{ template "content" . }}</body></html>{{ end }}`
	pageHTML := `{{ define "content" }}<h1>{{ .Title }}</h1>{{ end }}`

	if err := os.WriteFile(filepath.Join(templatesDir, "base.html"), []byte(baseHTML), 0644); err != nil {
		t.Fatal(err)
	}
	templateFiles := []string{"index.html", "about.html", "404.html", "tags.html", "archive.html", "blog.html", "post.html"}
	for _, f := range templateFiles {
		if err := os.WriteFile(filepath.Join(templatesDir, f), []byte(pageHTML), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// 3. Create dummy blog posts
	if err := os.MkdirAll(blogDir, 0755); err != nil {
		t.Fatal(err)
	}
	postMarkdown := `---
title: "Integration Post"
date: 2026-06-11T00:00:00Z
tags: ["integration"]
description: "A test post"
---
# Hello Integration
`
	if err := os.WriteFile(filepath.Join(blogDir, "test.md"), []byte(postMarkdown), 0644); err != nil {
		t.Fatal(err)
	}

	// 4. Create public static assets
	if err := os.MkdirAll(publicDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(publicDir, "test.txt"), []byte("assets"), 0644); err != nil {
		t.Fatal(err)
	}

	// 5. Run the pipeline
	count, err := RunPipeline(distDir, configDir, templatesDir, blogDir, publicDir)
	if err != nil {
		t.Fatalf("RunPipeline failed: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 post to be built, got %d", count)
	}

	// 6. Verify outputs exist
	expectedOutputs := []string{
		"index.html",
		"about.html",
		"404.html",
		"archive.html",
		"blog.html",
		"test.txt",
		"sitemap.xml",
		"rss.xml",
		"search-index.json",
		"llms.txt",
		filepath.Join("blog", "test.html"),
		filepath.Join("tags", "integration.html"),
		filepath.Join("api", "manifest.json"),
	}

	for _, output := range expectedOutputs {
		path := filepath.Join(distDir, output)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("Expected output file %s not found: %v", output, err)
		}
	}
}
