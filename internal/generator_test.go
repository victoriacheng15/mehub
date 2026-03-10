package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"mehub/internal/contents"
	"mehub/internal/post"
)

// Helper to create a dummy config
func createConfig() *contents.SiteConfig {
	return &contents.SiteConfig{
		Site: contents.SiteMetadata{
			Title:       "Test Site",
			URL:         "http://example.com/",
			Description: "Test Description",
		},
		Projects: []contents.Project{
			{
				Title:            "Test Project",
				ShortDescription: "Desc",
				Link:             "http://link",
				Techs:            []string{"Go", "Test"},
			},
		},
		Skills: []contents.Skill{
			{Name: "Go", Icon: "go.svg"},
		},
		Specialties: []string{"Testing"},
	}
}

// Helper to create dummy templates
func createTemplates(t *testing.T, dir string) {
	t.Helper()
	baseTmpl := `{{ define "base.html" }}<html><body>{{ template "content" . }}</body></html>{{ end }}`
	pageTmpl := `{{ define "content" }}<h1>{{ .Title }}</h1>{{ end }}`

	// Create internal/templates structure
	tmplDir := filepath.Join(dir, "internal", "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(tmplDir, "base.html"), []byte(baseTmpl), 0644); err != nil {
		t.Fatal(err)
	}

	files := []string{"index.html", "about.html", "now.html", "404.html", "tags.html", "archive.html", "blog.html", "post.html"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tmplDir, f), []byte(pageTmpl), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestFuncMap(t *testing.T) {
	gen := New(createConfig())

	tests := []struct {
		name     string
		fn       string
		input    interface{}
		input2   interface{} // For add/sub
		expected interface{}
	}{
		{
			name:     "cleanYAMLList - String List",
			fn:       "cleanYAMLList",
			input:    "- Item 1\n- Item 2",
			expected: 2, // Length check
		},
		{
			name:     "cleanYAMLList - Slice",
			fn:       "cleanYAMLList",
			input:    []string{"A", "B"},
			expected: 2,
		},
		{
			name:     "cleanYAMLList - Invalid",
			fn:       "cleanYAMLList",
			input:    123,
			expected: 0,
		},
		{
			name:     "Math - Add",
			fn:       "add",
			input:    2,
			input2:   3,
			expected: 5,
		},
		{
			name:     "Math - Sub",
			fn:       "sub",
			input:    5,
			input2:   3,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.fn {
			case "cleanYAMLList":
				fn := gen.FuncMap[tt.fn].(func(interface{}) []string)
				res := fn(tt.input)
				if len(res) != tt.expected.(int) {
					t.Errorf("Expected length %d, got %d", tt.expected, len(res))
				}
			case "add", "sub":
				fn := gen.FuncMap[tt.fn].(func(int, int) int)
				got := fn(tt.input.(int), tt.input2.(int))
				if got != tt.expected.(int) {
					t.Errorf("%s(%d, %d) = %d; want %d", tt.fn, tt.input, tt.input2, got, tt.expected)
				}
			}
		})
	}
}

func TestRenderPage(t *testing.T) {
	tmpDir := t.TempDir()
	wd, _ := os.Getwd()
	defer os.Chdir(wd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	createTemplates(t, tmpDir)

	gen := New(createConfig())
	distDir := filepath.Join(tmpDir, "dist")

	tests := []struct {
		name        string
		filename    string
		tmplPath    string
		titlePrefix string
		data        PageData
		setup       func()
		wantErr     bool
	}{
		{
			name:        "Render Index",
			filename:    "index.html",
			tmplPath:    "internal/templates/index.html",
			titlePrefix: "Home",
			data:        PageData{},
			wantErr:     false,
		},
		{
			name:     "Invalid Template Path",
			filename: "fail.html",
			tmplPath: "nonexistent.html",
			wantErr:  true,
		},
		{
			name:     "Output Dir Creation Failure",
			filename: "subdir/fail.html",
			tmplPath: "internal/templates/index.html",
			setup: func() {
				os.MkdirAll(distDir, 0755)
				os.WriteFile(filepath.Join(distDir, "subdir"), []byte("block"), 0644)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			err := gen.RenderPage(distDir, tt.filename, tt.tmplPath, tt.titlePrefix, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderPage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if _, err := os.Stat(filepath.Join(distDir, tt.filename)); err != nil {
					t.Errorf("Output file %s not found: %v", tt.filename, err)
				}
			}
		})
	}
}

func TestGenerators(t *testing.T) {
	tmpDir := t.TempDir()
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	createTemplates(t, tmpDir)

	gen := New(createConfig())
	distDir := filepath.Join(tmpDir, "dist")

	posts := []post.Post{
		{
			Frontmatter: post.Frontmatter{
				Title: "Test Post",
				Date:  time.Now(),
			},
			Slug: "test-post",
		},
	}
	data := &post.ContentData{
		Posts: posts,
		PostsByTag: map[string][]post.Post{
			"go": posts,
		},
	}

	tests := []struct {
		name    string
		fn      func() error
		check   func() error
		wantErr bool
	}{
		{
			name: "Static Pages",
			fn:   func() error { return gen.GenerateStaticPages(distDir, data) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "about.html"))
				return err
			},
		},
		{
			name: "RSS",
			fn:   func() error { return gen.GenerateRSS(distDir, posts) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "rss.xml"))
				return err
			},
		},
		{
			name: "Sitemap",
			fn:   func() error { return gen.GenerateSitemap(distDir, posts) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "sitemap.xml"))
				return err
			},
		},
		{
			name: "Search Index",
			fn:   func() error { return gen.GenerateSearchIndex(distDir, data) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "search-index.json"))
				return err
			},
		},
		{
			name: "Blog Pagination",
			fn:   func() error { return gen.GenerateBlogPagination(distDir, data, 10) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "blog.html"))
				return err
			},
		},
		{
			name: "Tag Pages",
			fn:   func() error { return gen.GenerateTagPages(distDir, data) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "tags", "go.html"))
				return err
			},
		},
		{
			name: "Post Pages",
			fn:   func() error { return gen.GeneratePostPages(distDir, data) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "blog", "test-post.html"))
				return err
			},
		},
		{
			name: "Unified Registries",
			fn:   func() error { return gen.GenerateRegistries(distDir, data) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "api", "manifest.json"))
				return err
			},
		},
		{
			name: "LLMs Txt",
			fn:   func() error { return gen.GenerateLLMsTxt(distDir) },
			check: func() error {
				content, err := os.ReadFile(filepath.Join(distDir, "llms.txt"))
				if err != nil {
					return err
				}
				if !strings.Contains(string(content), "http://example.com/api/manifest.json") {
					return fmt.Errorf("llms.txt missing manifest URL")
				}
				return nil
			},
		},
		{
			name: "Failure - Missing Template",
			fn: func() error {
				if err := os.Remove(filepath.Join(tmpDir, "internal", "templates", "about.html")); err != nil {
					return err
				}
				defer createTemplates(t, tmpDir) // Restore for next tests
				return gen.GenerateStaticPages(distDir, data)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Got error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.check != nil {
				if err := tt.check(); err != nil {
					t.Errorf("Check failed: %v", err)
				}
			}
		})
	}
}

func TestBuild(t *testing.T) {
	setup := func(t *testing.T) (string, *SiteGenerator, *post.ContentData) {
		tmpDir := t.TempDir()
		createTemplates(t, tmpDir)
		gen := New(createConfig())
		posts := []post.Post{{Slug: "test", Frontmatter: post.Frontmatter{Title: "T", Date: time.Now()}}}
		data := &post.ContentData{Posts: posts}
		return tmpDir, gen, data
	}

	tests := []struct {
		name    string
		setup   func(t *testing.T, tmpDir string)
		wantErr bool
	}{
		{
			name:    "Success",
			setup:   func(t *testing.T, tmpDir string) {},
			wantErr: false,
		},
		{
			name: "Failure - Missing Template",
			setup: func(t *testing.T, tmpDir string) {
				os.Remove(filepath.Join(tmpDir, "internal", "templates", "index.html"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, gen, data := setup(t)

			// We need to change to tmpDir because RenderPage has hardcoded internal/templates
			wd, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(wd)

			tt.setup(t, tmpDir)

			distDir := filepath.Join(tmpDir, "dist")
			err := gen.Build(distDir, data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
