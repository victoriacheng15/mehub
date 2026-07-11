package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"mehub/internal"
)

type testContext struct {
	tmpDir       string
	configDir    string
	blogDir      string
	distDir      string
	templatesDir string
	publicDir    string
	postCount    int
	err          error
}

// ============================================================================
// E2E Test Runner & Scenario Initialization
// ============================================================================

// TestFeatures acts as the primary runner for Gherkin E2E integration test scenarios.
// It initializes and runs the Cucumber-compatible Godog test suite, linking its output
// telemetry and assertion results directly to Go's standard testing runner framework.
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	tc := &testContext{}

	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		tmpDir, err := os.MkdirTemp("", "mehub-e2e-")
		if err != nil {
			return ctx, err
		}
		tc.tmpDir = tmpDir
		return ctx, nil
	})

	sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if tc.tmpDir != "" {
			os.RemoveAll(tc.tmpDir)
		}
		return ctx, nil
	})

	sc.Step(`^a configuration directory with a valid profile$`, tc.setupValidConfig)
	sc.Step(`^a configuration directory with a missing profile$`, tc.setupMissingConfig)
	sc.Step(`^a blog directory containing (\d+) published post$`, tc.setupSinglePost)
	sc.Step(`^a blog directory containing (\d+) published post and (\d+) draft post$`, tc.setupMixedPosts)
	sc.Step(`^a static assets directory containing a file "([^"]*)"$`, tc.setupStaticAsset)
	sc.Step(`^the build pipeline is executed$`, tc.runPipeline)
	sc.Step(`^the output directory should contain "([^"]*)"$`, tc.checkFileExists)
	sc.Step(`^the output directory should not contain "([^"]*)"$`, tc.checkFileMissing)
	sc.Step(`^the build pipeline execution should fail$`, tc.checkPipelineFailed)
	sc.Step(`^the output file "([^"]*)" should contain "([^"]*)"$`, tc.checkFileContains)
	sc.Step(`^the output file "([^"]*)" should not contain "([^"]*)"$`, tc.checkFileDoesNotContain)
}

// ============================================================================
// Pipeline Feature Helpers
// ============================================================================

// setupValidConfig creates a temporary contents directory containing a complete set of valid YAML configs.
func (tc *testContext) setupValidConfig() error {
	tc.configDir = filepath.Join(tc.tmpDir, "contents")
	if err := os.MkdirAll(tc.configDir, 0755); err != nil {
		return err
	}

	configYAML := `
landing:
  title: "E2E Test Site"
  url: "https://example.com/"
about:
  paragraphs:
    - "hello e2e"
navigation:
  header: []
skills: []
socials: []
`
	configs := map[string]string{
		"config.yaml":   configYAML,
		"projects.yaml": `projects: []`,
	}
	for file, content := range configs {
		if err := os.WriteFile(filepath.Join(tc.configDir, file), []byte(content), 0644); err != nil {
			return err
		}
	}

	// Create dummy templates
	tc.templatesDir = filepath.Join(tc.tmpDir, "templates")
	if err := os.MkdirAll(tc.templatesDir, 0755); err != nil {
		return err
	}
	baseHTML := `{{ define "base.html" }}<html><body>{{ template "content" . }}</body></html>{{ end }}`
	pageHTML := `{{ define "content" }}<h1>{{ .Title }}</h1>{{ end }}`

	if err := os.WriteFile(filepath.Join(tc.templatesDir, "base.html"), []byte(baseHTML), 0644); err != nil {
		return err
	}
	templateFiles := []string{"index.html", "about.html", "now.html", "404.html", "tags.html", "archive.html", "blog.html", "post.html"}
	for _, f := range templateFiles {
		if err := os.WriteFile(filepath.Join(tc.templatesDir, f), []byte(pageHTML), 0644); err != nil {
			return err
		}
	}

	return nil
}

// setupMissingConfig creates a contents directory with navigation/metadata but missing the critical profile.yaml.
func (tc *testContext) setupMissingConfig() error {
	tc.configDir = filepath.Join(tc.tmpDir, "contents")
	if err := os.MkdirAll(tc.configDir, 0755); err != nil {
		return err
	}
	configs := map[string]string{
		"projects.yaml": `projects: []`,
	}
	for file, content := range configs {
		if err := os.WriteFile(filepath.Join(tc.configDir, file), []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

// setupStaticAsset creates a static directory and writes a dummy asset file into it.
func (tc *testContext) setupStaticAsset(filename string) error {
	tc.publicDir = filepath.Join(tc.tmpDir, "static")
	if err := os.MkdirAll(tc.publicDir, 0755); err != nil {
		return err
	}
	content := "User-agent: *\nDisallow: /"
	if err := os.WriteFile(filepath.Join(tc.publicDir, filename), []byte(content), 0644); err != nil {
		return err
	}
	return nil
}

// setupSinglePost writes a single dummy published blog post to the blog content directory.
func (tc *testContext) setupSinglePost(count int) error {
	tc.blogDir = filepath.Join(tc.tmpDir, "blog")
	if err := os.MkdirAll(tc.blogDir, 0755); err != nil {
		return err
	}

	postMarkdown := `---
title: "E2E Post"
date: 2026-06-11T00:00:00Z
tags: ["e2e"]
description: "A test post"
---
# Hello E2E
`
	if err := os.WriteFile(filepath.Join(tc.blogDir, "test.md"), []byte(postMarkdown), 0644); err != nil {
		return err
	}
	return nil
}

// runPipeline executes the generator build orchestrator, capturing the output state and any error returned.
func (tc *testContext) runPipeline() error {
	tc.distDir = filepath.Join(tc.tmpDir, "dist")
	tc.publicDir = filepath.Join(tc.tmpDir, "static")
	if err := os.MkdirAll(tc.publicDir, 0755); err != nil {
		return err
	}

	tc.postCount, tc.err = internal.RunPipeline(tc.distDir, tc.configDir, tc.templatesDir, tc.blogDir, tc.publicDir)
	return nil
}

// checkFileExists asserts that the target file path exists inside the generated output directory.
func (tc *testContext) checkFileExists(filename string) error {
	if tc.err != nil {
		return fmt.Errorf("pipeline run failed: %w", tc.err)
	}
	path := filepath.Join(tc.distDir, filename)
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("expected file %s does not exist: %w", filename, err)
	}
	return nil
}

// checkPipelineFailed asserts that the pipeline run failed and returned a non-nil error.
func (tc *testContext) checkPipelineFailed() error {
	if tc.err == nil {
		return fmt.Errorf("expected pipeline execution to fail, but it succeeded")
	}
	return nil
}

// ============================================================================
// Draft Filtering Feature Helpers
// ============================================================================

// setupMixedPosts generates a set of published and draft markdown posts in the blog content directory.
func (tc *testContext) setupMixedPosts(published, drafts int) error {
	tc.blogDir = filepath.Join(tc.tmpDir, "blog")
	if err := os.MkdirAll(tc.blogDir, 0755); err != nil {
		return err
	}

	for i := 1; i <= published; i++ {
		content := fmt.Sprintf(`---
title: "Published Post %d"
date: 2026-06-11T00:00:00Z
tags: ["e2e"]
description: "A test post"
draft: false
---
# Published %d
`, i, i)
		filename := fmt.Sprintf("published-%d.md", i)
		if err := os.WriteFile(filepath.Join(tc.blogDir, filename), []byte(content), 0644); err != nil {
			return err
		}
	}

	for i := 1; i <= drafts; i++ {
		content := fmt.Sprintf(`---
title: "Draft Post %d"
date: 2026-06-12T00:00:00Z
tags: ["e2e"]
description: "A draft post"
draft: true
---
# Draft %d
`, i, i)
		filename := fmt.Sprintf("draft-%d.md", i)
		if err := os.WriteFile(filepath.Join(tc.blogDir, filename), []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// checkFileMissing asserts that the target file path is not present in the generated output directory.
func (tc *testContext) checkFileMissing(filename string) error {
	if tc.err != nil {
		return fmt.Errorf("pipeline run failed: %w", tc.err)
	}
	path := filepath.Join(tc.distDir, filename)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("expected file %s to not exist, but it does", filename)
	}
	return nil
}

// checkFileContains asserts that the target file exists and contains the specified search string.
func (tc *testContext) checkFileContains(filename, searchStr string) error {
	if tc.err != nil {
		return fmt.Errorf("pipeline run failed: %w", tc.err)
	}
	path := filepath.Join(tc.distDir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	if !strings.Contains(string(data), searchStr) {
		return fmt.Errorf("expected file %s to contain %q, but it did not", filename, searchStr)
	}
	return nil
}

// checkFileDoesNotContain asserts that if the target file exists, it does not contain the specified search string.
func (tc *testContext) checkFileDoesNotContain(filename, searchStr string) error {
	if tc.err != nil {
		return fmt.Errorf("pipeline run failed: %w", tc.err)
	}
	path := filepath.Join(tc.distDir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	if strings.Contains(string(data), searchStr) {
		return fmt.Errorf("expected file %s to not contain %q, but it did", filename, searchStr)
	}
	return nil
}
