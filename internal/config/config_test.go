package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string // filename -> content
		validate func(*testing.T, *SiteConfig)
		wantErr  bool
	}{
		{
			name: "Happy Path - All files valid",
			files: map[string]string{
				"site.yaml": `
site:
  title: "My Site"
  description: "A cool site"
  about:
    paragraphs: |
      - Hello world
`,
				"navigation.yaml": `
navigation:
  header:
    - href: "/"
      text: "Home"
`,
				"socials.yaml": `
socials:
  - name: "GitHub"
    href: "https://github.com"
`,
				"projects.yaml": `
projects:
  - title: "Project A"
`,
				"skills.yaml": `
skills:
  - name: "Go"
`,
			},
			validate: func(t *testing.T, cfg *SiteConfig) {
				if cfg.Site.Title != "My Site" {
					t.Errorf("Expected Site.Title to be 'My Site', got '%s'", cfg.Site.Title)
				}
				if len(cfg.Navigation.Header) != 1 {
					t.Errorf("Expected 1 header nav item, got %d", len(cfg.Navigation.Header))
				}
				if cfg.Navigation.Header[0].Text != "Home" {
					t.Errorf("Expected first nav text 'Home', got '%s'", cfg.Navigation.Header[0].Text)
				}
			},
			wantErr: false,
		},
		{
			name: "Missing File",
			files: map[string]string{
				"site.yaml": `site: { title: "Test" }`,
			},
			validate: nil,
			wantErr:  true,
		},
		{
			name: "Invalid YAML",
			files: map[string]string{
				"site.yaml": `
site:
	title: "Invalid Tab"
`,
			},
			validate: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			for filename, content := range tt.files {
				path := filepath.Join(tmpDir, filename)
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create config file %s: %v", filename, err)
				}
			}

			cfg, err := LoadConfig(tmpDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}
