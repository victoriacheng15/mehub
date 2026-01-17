package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	wd, _ := os.Getwd()
	defer os.Chdir(wd)

	setupSuccess := func(t *testing.T, root string) (string, string, string, string) {
		dist := "dist"
		conf := "configs"
		blog := "blog"
		pub := "public"
		tmpl := "internal/templates"

		// Create directories relative to root
		for _, d := range []string{conf, blog, pub, tmpl} {
			if err := os.MkdirAll(filepath.Join(root, d), 0755); err != nil {
				t.Fatal(err)
			}
		}

		// Create minimal templates
		base := `{{ define "base.html" }}<html><body>{{ template "content" . }}</body></html>{{ end }}`
		content := `{{ define "content" }}<h1>{{ .Title }}</h1>{{ end }}`
		tmpls := []string{"base.html", "index.html", "about.html", "now.html", "404.html", "tags.html", "archive.html", "blog.html", "post.html"}
		for _, f := range tmpls {
			c := content
			if f == "base.html" {
				c = base
			}
			if err := os.WriteFile(filepath.Join(root, tmpl, f), []byte(c), 0644); err != nil {
				t.Fatal(err)
			}
		}

		// Create minimal configs
		configs := map[string]string{
			"site.yaml":       "site: { title: 'Test', url: 'http://x.com' }",
			"navigation.yaml": "navigation: { header: [], footer: [] }",
			"socials.yaml":    "socials: []",
			"projects.yaml":   "projects: []",
			"skills.yaml":     "skills: []",
		}
		for f, c := range configs {
			if err := os.WriteFile(filepath.Join(root, conf, f), []byte(c), 0644); err != nil {
				t.Fatal(err)
			}
		}

		// Create a blog post
		post := "---\ntitle: 'Hi'\ndate: 2023-01-01T00:00:00Z\n---\nBody"
		if err := os.WriteFile(filepath.Join(root, blog, "hi.md"), []byte(post), 0644); err != nil {
			t.Fatal(err)
		}

		// We need to be in 'root' because generator.go has hardcoded "internal/templates/..."
		if err := os.Chdir(root); err != nil {
			t.Fatal(err)
		}

		return dist, conf, blog, pub
	}

	tests := []struct {
		name    string
		setup   func(t *testing.T, root string) (dist, conf, blog, pub string)
		wantErr bool
	}{
		{
			name:    "Success - Minimal Build",
			setup:   setupSuccess,
			wantErr: false,
		},
		{
			name: "Failure - Missing Config",
			setup: func(t *testing.T, root string) (string, string, string, string) {
				if err := os.Chdir(root); err != nil {
					t.Fatal(err)
				}
				return "dist", "nonexistent-configs", "blog", "public"
			},
			wantErr: true,
		},
		{
			name: "Failure - Dist Dir Creation",
			setup: func(t *testing.T, root string) (string, string, string, string) {
				// Create a read-only directory
				locked := filepath.Join(root, "locked")
				if err := os.Mkdir(locked, 0555); err != nil {
					t.Fatal(err)
				}
				dist := filepath.Join(locked, "dist")

				// We need valid other dirs
				_, conf, blog, pub := setupSuccess(t, root)

				return dist, conf, blog, pub
			},
			wantErr: true,
		},
		{
			name: "Failure - Load Posts",
			setup: func(t *testing.T, root string) (string, string, string, string) {
				dist, conf, _, pub := setupSuccess(t, root)
				return dist, conf, "nonexistent-blog", pub
			},
			wantErr: true,
		},
		{
			name: "Failure - Generator (Static Pages)",
			setup: func(t *testing.T, root string) (string, string, string, string) {
				dist, conf, blog, pub := setupSuccess(t, root)
				// Remove a template to cause generation failure
				if err := os.Remove(filepath.Join(root, "internal", "templates", "about.html")); err != nil {
					t.Fatal(err)
				}
				return dist, conf, blog, pub
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			dist, conf, blog, pub := tt.setup(t, tmpDir)

			err := run(dist, conf, blog, pub)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify some output
				if _, err := os.Stat(filepath.Join(tmpDir, dist, "index.html")); err != nil {
					t.Errorf("index.html not generated: %v", err)
				}
			}
			// Revert to original WD for next test case setup
			os.Chdir(wd)
		})
	}
}
