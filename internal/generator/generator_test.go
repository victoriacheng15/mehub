package generator

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"mehub/internal/config"
	"mehub/internal/content"
)

// Helper to create a dummy config
func createConfig() *config.SiteConfig {
	return &config.SiteConfig{
		Site: config.SiteMetadata{
			Title:       "Test Site",
			URL:         "http://example.com/",
			Description: "Test Description",
		},
		Projects: []config.Project{
			{
				Title:            "Test Project",
				ShortDescription: "Desc",
				Link:             "http://link",
				Techs:            "- Go\n- Test",
			},
		},
		Skills: []config.Skill{
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

	t.Run("cleanYAMLList", func(t *testing.T) {
		tests := []struct {
			name     string
			input    interface{}
			expected int
		}{
			{"String List", "- Item 1\n- Item 2", 2},
			{"Slice", []string{"A", "B"}, 2},
			{"Invalid", 123, 0}, // Should return nil/empty
		}

		cleanFunc := gen.FuncMap["cleanYAMLList"].(func(interface{}) []string)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res := cleanFunc(tt.input)
				if len(res) != tt.expected {
					t.Errorf("Expected length %d, got %d", tt.expected, len(res))
				}
			})
		}
	})

	t.Run("Math", func(t *testing.T) {
		tests := []struct {
			name     string
			fn       string
			a, b     int
			expected int
		}{
			{"add", "add", 2, 3, 5},
			{"sub", "sub", 5, 3, 2},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				fn := gen.FuncMap[tt.fn].(func(int, int) int)
				if got := fn(tt.a, tt.b); got != tt.expected {
					t.Errorf("%s(%d, %d) = %d; want %d", tt.fn, tt.a, tt.b, got, tt.expected)
				}
			})
		}
	})
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
		setup       func() // Optional setup for failure cases
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
				// Create a file named 'dist/subdir' to block directory creation
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
				content, err := os.ReadFile(filepath.Join(distDir, tt.filename))
				if err != nil {
					t.Fatalf("Failed to read output file: %v", err)
				}
				if len(content) == 0 {
					t.Error("Output file is empty")
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

	posts := []content.Post{
		{
			Frontmatter: content.Frontmatter{
				Title: "Test Post",
				Date:  time.Now(),
			},
			Slug: "test-post",
		},
	}
	data := &content.ContentData{
		Posts: posts,
		PostsByTag: map[string][]content.Post{
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
			wantErr: false,
		},
		{
			name: "RSS",
			fn:   func() error { return gen.GenerateRSS(distDir, posts) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "rss.xml"))
				return err
			},
			wantErr: false,
		},
		{
			name: "Sitemap",
			fn:   func() error { return gen.GenerateSitemap(distDir, posts) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "sitemap.xml"))
				return err
			},
			wantErr: false,
		},
		{
			name: "Search Index",
			fn:   func() error { return gen.GenerateSearchIndex(distDir, data) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "search-index.json"))
				return err
			},
			wantErr: false,
		},
		{
			name: "Blog Pagination",
			fn:   func() error { return gen.GenerateBlogPagination(distDir, data, 10) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "blog.html"))
				return err
			},
			wantErr: false,
		},
		{
			name: "Tag Pages",
			fn:   func() error { return gen.GenerateTagPages(distDir, data) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "tags", "go.html"))
				return err
			},
			wantErr: false,
		},
		{
			name: "Post Pages",
			fn:   func() error { return gen.GeneratePostPages(distDir, data) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "blog", "test-post.html"))
				return err
			},
			wantErr: false,
		},
		{
			name: "Failure - Static Pages (Missing Template)",
			fn: func() error {
				// Remove a template to cause failure
				if err := os.Remove(filepath.Join(tmpDir, "internal", "templates", "about.html")); err != nil {
					return err
				}
				return gen.GenerateStaticPages(distDir, data)
			},
			check:   nil,
			wantErr: true,
		},
		{
			name: "Failure - Blog Pagination (Missing Template)",
			fn: func() error {
				// Restore about.html first if needed (though we use tmpDir per test run in separate subtests? No, same tmpDir for all TestGenerators cases in the loop? No wait)
				// TestGenerators uses `tmpDir := t.TempDir()` at the top.
				// The tests loop runs sequentially in `t.Run`.
				// So if I delete a template in one test, it's gone for the next.
				// I should restore it or ensure independence.
				// Currently, `t.Run` uses the same `gen` instance and `distDir`.
				// But `t.TempDir()` is unique per test function call.
				// Wait, `TestGenerators` calls `t.TempDir()` ONCE.
				// So the state persists across subtests.
				// This is bad for "Failure" tests that destroy state.

				// I should probably re-create the template or use a separate setup.
				// But let's just create it back.
				createTemplates(t, tmpDir)
				if err := os.Remove(filepath.Join(tmpDir, "internal", "templates", "blog.html")); err != nil {
					return err
				}
				return gen.GenerateBlogPagination(distDir, data, 10)
			},
			check:   nil,
			wantErr: true,
		},
		{
			name: "Failure - Tag Pages (Missing Template)",
			fn: func() error {
				createTemplates(t, tmpDir)
				if err := os.Remove(filepath.Join(tmpDir, "internal", "templates", "blog.html")); err != nil {
					return err
				}
				return gen.GenerateTagPages(distDir, data)
			},
			check:   nil,
			wantErr: true,
		},
		{
			name: "Failure - Post Pages (Missing Template)",
			fn: func() error {
				createTemplates(t, tmpDir)
				if err := os.Remove(filepath.Join(tmpDir, "internal", "templates", "post.html")); err != nil {
					return err
				}
				return gen.GeneratePostPages(distDir, data)
			},
			check:   nil,
			wantErr: true,
		},
		{
			name: "Blog Registry",
			fn:   func() error { return gen.GenerateBlogRegistry(distDir, data) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "api", "blog-registry.json"))
				return err
			},
			wantErr: false,
		},
		{
			name: "Projects Registry",
			fn:   func() error { return gen.GenerateProjectsRegistry(distDir) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "api", "projects-registry.json"))
				return err
			},
			wantErr: false,
		},
		{
			name: "Skills Registry",
			fn:   func() error { return gen.GenerateSkillsRegistry(distDir) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "api", "skills-registry.json"))
				return err
			},
			wantErr: false,
		},
		{
			name: "Profile Registry",
			fn:   func() error { return gen.GenerateProfileRegistry(distDir) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "api", "profile-registry.json"))
				return err
			},
			wantErr: false,
		},
		{
			name: "LLMs Txt",
			fn:   func() error { return gen.GenerateLLMsTxt(distDir) },
			check: func() error {
				_, err := os.Stat(filepath.Join(distDir, "llms.txt"))
				return err
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Generator error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.check != nil {
				if err := tt.check(); err != nil {
					t.Errorf("Check failed: %v", err)
				}
			}
		})
	}
}
