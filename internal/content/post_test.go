package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParsePost(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  string
		validate func(*testing.T, *Post)
		wantErr  bool
	}{
		{
			name:     "Valid Post",
			filename: "valid.md",
			content: `---
title: "Hello World"
date: 2023-01-01T00:00:00Z
tags: ["go", "testing"]
description: "A test post"
---
# Hello
This is a test.
`,
			validate: func(t *testing.T, post *Post) {
				if post.Title != "Hello World" {
					t.Errorf("Expected title 'Hello World', got '%s'", post.Title)
				}
				if post.Slug != "valid" {
					t.Errorf("Expected slug 'valid', got '%s'", post.Slug)
				}
				if len(post.Tags) != 2 {
					t.Errorf("Expected 2 tags, got %d", len(post.Tags))
				}
				if !strings.Contains(post.Content, "<h1>Hello</h1>") {
					t.Errorf("Expected HTML content to contain <h1>Hello</h1>, got %s", post.Content)
				}
			},
			wantErr: false,
		},
		{
			name:     "Mermaid Diagram Replacement",
			filename: "mermaid.md",
			content: `---
title: "Mermaid"
---
<pre><code class="language-mermaid">graph TD; A-->B;</code></pre>
`,
			validate: func(t *testing.T, post *Post) {
				expected := `<div class="mermaid">graph TD; A-->B;</div>`
				if !strings.Contains(post.Content, expected) {
					t.Errorf("Expected mermaid div, got: %s", post.Content)
				}
			},
			wantErr: false,
		},
		{
			name:     "Invalid Frontmatter",
			filename: "invalid_fm.md",
			content: `---
title: [Broken
---
Content
`,
			validate: nil,
			wantErr:  true,
		},
		{
			name:     "Missing Frontmatter",
			filename: "no_fm.md",
			content:  `Just some content without delimiters`,
			validate: func(t *testing.T, post *Post) {
				if post != nil {
					t.Error("Expected nil post for missing frontmatter, got struct")
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, tt.filename)
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to write post: %v", err)
			}

			post, err := ParsePost(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, post)
			}
		})
	}

	t.Run("File Read Error", func(t *testing.T) {
		tmpDir := t.TempDir()
		_, err := ParsePost(filepath.Join(tmpDir, "nonexistent.md"))
		if err == nil {
			t.Error("Expected error reading nonexistent file, got nil")
		}
	})
}

func TestProcessPosts(t *testing.T) {
	date1, _ := time.Parse("2006-01-02", "2023-05-01")
	date2, _ := time.Parse("2006-01-02", "2022-12-01")
	date3, _ := time.Parse("2006-01-02", "2023-01-10")

	tests := []struct {
		name     string
		posts    []Post
		validate func(*testing.T, *ContentData)
	}{
		{
			name: "Group by Tags and Year",
			posts: []Post{
				{
					Frontmatter: Frontmatter{Title: "Post 1", Date: date1, Tags: []string{"go", "web"}},
					Slug:        "post-1",
				},
				{
					Frontmatter: Frontmatter{Title: "Post 2", Date: date2, Tags: []string{"web"}},
					Slug:        "post-2",
				},
				{
					Frontmatter: Frontmatter{Title: "Post 3", Date: date3, Tags: []string{"go"}},
					Slug:        "post-3",
				},
			},
			validate: func(t *testing.T, data *ContentData) {
				// Check Tags
				if len(data.Tags) != 2 {
					t.Errorf("Expected 2 unique tags, got %d", len(data.Tags))
				}
				if data.Tags[0] != "go" || data.Tags[1] != "web" {
					t.Errorf("Tags not sorted correctly: %v", data.Tags)
				}

				// Check Tag Counts
				if data.TagCounts["go"] != 2 {
					t.Errorf("Expected 2 posts for 'go', got %d", data.TagCounts["go"])
				}
				if data.TagCounts["web"] != 2 {
					t.Errorf("Expected 2 posts for 'web', got %d", data.TagCounts["web"])
				}

				// Check PostsByYear
				if len(data.PostsByYear[2023]) != 2 {
					t.Errorf("Expected 2 posts in 2023, got %d", len(data.PostsByYear[2023]))
				}
				if len(data.PostsByYear[2022]) != 1 {
					t.Errorf("Expected 1 post in 2022, got %d", len(data.PostsByYear[2022]))
				}

				// Check ArchiveYears order (descending)
				if len(data.ArchiveYears) != 2 {
					t.Errorf("Expected 2 archive years, got %d", len(data.ArchiveYears))
				}
				if data.ArchiveYears[0] != 2023 || data.ArchiveYears[1] != 2022 {
					t.Errorf("Archive years not sorted descending: %v", data.ArchiveYears)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := ProcessPosts(tt.posts)
			tt.validate(t, data)
		})
	}
}

func TestGetPosts(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string // filename -> content
		validate func(*testing.T, []Post)
		wantErr  bool
	}{
		{
			name: "Valid Posts and Ignore non-md",
			files: map[string]string{
				"post1.md": `---
title: "Valid"
date: 2023-01-01T00:00:00Z
tags: ["a"]
description: "d"
---
Content`,
				"ignore.txt": "ignore me",
			},
			validate: func(t *testing.T, posts []Post) {
				if len(posts) != 1 {
					t.Errorf("Expected 1 post, got %d", len(posts))
				}
				if len(posts) > 0 && posts[0].Title != "Valid" {
					t.Errorf("Expected title 'Valid', got '%s'", posts[0].Title)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for filename, content := range tt.files {
				if err := os.WriteFile(filepath.Join(tmpDir, filename), []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create file %s: %v", filename, err)
				}
			}

			posts, err := GetPosts(tmpDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, posts)
			}
		})
	}

	t.Run("Directory Not Found", func(t *testing.T) {
		tmpDir := t.TempDir()
		_, err := GetPosts(filepath.Join(tmpDir, "nonexistent"))
		if err == nil {
			t.Error("Expected error for nonexistent directory, got nil")
		}
	})
}
