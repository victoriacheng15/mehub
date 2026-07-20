package internal

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

type SiteGenerator struct {
	Config            *SiteConfig
	FuncMap           template.FuncMap
	TemplatesDir      string
	minifier          *minify.M
	totalOriginalSize int64
	totalMinifiedSize int64
}

func New(cfg *SiteConfig, templatesDir string) *SiteGenerator {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)

	return &SiteGenerator{
		Config:       cfg,
		TemplatesDir: templatesDir,
		minifier:     m,
		FuncMap: template.FuncMap{
			"split":             strings.Split,
			"replace":           strings.ReplaceAll,
			"trimSpace":         strings.TrimSpace,
			"stringsHasPrefix":  strings.HasPrefix,
			"stringsTrimPrefix": strings.TrimPrefix,
			"add":               func(a, b int) int { return a + b },
			"sub":               func(a, b int) int { return a - b },
			"safeHTML":          func(s string) template.HTML { return template.HTML(s) },
			"cleanYAMLList": func(data interface{}) []string {
				var input string
				switch v := data.(type) {
				case string:
					input = v
				case []string:
					return v
				default:
					return nil
				}

				lines := strings.Split(strings.TrimSpace(input), "\n")
				var result []string
				for _, line := range lines {
					cleaned := strings.TrimSpace(line)
					cleaned = strings.TrimPrefix(cleaned, "- ")
					cleaned = strings.Trim(cleaned, "\"")
					if cleaned != "" {
						result = append(result, cleaned)
					}
				}
				return result
			},
		},
	}
}

func (g *SiteGenerator) RenderPage(dir, filename, tmplPath string, titlePrefix string, data PageData) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create dir %s: %w", dir, err)
	}

	basePath := filepath.Join(g.TemplatesDir, "base.html")
	fullTmplPath := filepath.Join(g.TemplatesDir, tmplPath)
	tmpl, err := template.New("base.html").Funcs(g.FuncMap).ParseFiles(basePath, fullTmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse templates for %s: %w", fullTmplPath, err)
	}

	outputFile, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", filename, err)
	}
	defer outputFile.Close()

	title := g.Config.Landing.Title
	if titlePrefix != "" {
		title = titlePrefix + " | " + title
	}

	data.Config = g.Config
	data.CurrentYear = time.Now().Year()
	data.Title = title

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base.html", data); err != nil {
		return fmt.Errorf("failed to execute template for %s: %w", fullTmplPath, err)
	}

	originalSize := buf.Len()
	minifiedHTML, err := g.minifier.Bytes("text/html", buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to minify HTML for %s: %w", filename, err)
	}
	minifiedSize := len(minifiedHTML)

	g.totalOriginalSize += int64(originalSize)
	g.totalMinifiedSize += int64(minifiedSize)

	if _, err := outputFile.Write(minifiedHTML); err != nil {
		return fmt.Errorf("failed to write minified HTML to %s: %w", filename, err)
	}
	return nil
}

func (g *SiteGenerator) GenerateStaticPages(distDir string, data *ContentData) error {
	pages := []struct {
		filename    string
		tmplPath    string
		titlePrefix string
		data        PageData
	}{
		{"index.html", "index.html", "", PageData{}},
		{"about.html", "about.html", "About", PageData{}},
		{"404.html", "404.html", "404 - Not Found", PageData{}},
		{"archive.html", "archive.html", "Archive", PageData{Archive: data.PostsByYear, ArchiveYears: data.ArchiveYears}},
	}

	for _, p := range pages {
		if err := g.RenderPage(distDir, p.filename, p.tmplPath, p.titlePrefix, p.data); err != nil {
			return err
		}
	}
	return nil
}

func (g *SiteGenerator) GenerateBlogPagination(distDir string, data *ContentData, pageSize int) error {
	totalPages := (len(data.Posts) + pageSize - 1) / pageSize
	for i := 0; i < totalPages; i++ {
		startIdx := i * pageSize
		endIdx := startIdx + pageSize
		if endIdx > len(data.Posts) {
			endIdx = len(data.Posts)
		}
		pagePosts := data.Posts[startIdx:endIdx]
		pageNumber := i + 1

		if pageNumber == 1 {
			if err := g.RenderPage(distDir, "blog.html", "blog.html", "Blog", PageData{
				Posts:       pagePosts,
				CurrentPage: pageNumber,
				TotalPages:  totalPages,
				Tags:        data.Tags,
				TagCounts:   data.TagCounts,
			}); err != nil {
				return err
			}
		} else {
			pageDir := filepath.Join(distDir, "blog")
			if err := g.RenderPage(pageDir, fmt.Sprintf("%d.html", pageNumber), "blog.html", fmt.Sprintf("Blog - Page %d", pageNumber), PageData{
				Posts:       pagePosts,
				CurrentPage: pageNumber,
				TotalPages:  totalPages,
				PathPrefix:  "../",
				Tags:        data.Tags,
				TagCounts:   data.TagCounts,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *SiteGenerator) GenerateTagPages(distDir string, data *ContentData) error {
	tagsDistDir := filepath.Join(distDir, "tags")
	for tag, tagPosts := range data.PostsByTag {
		if err := g.RenderPage(tagsDistDir, tag+".html", "blog.html", "#"+tag, PageData{
			Posts:      tagPosts,
			PathPrefix: "../",
			Tags:       data.Tags,
			TagCounts:  data.TagCounts,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (g *SiteGenerator) GeneratePostPages(distDir string, data *ContentData) error {
	blogDistDir := filepath.Join(distDir, "blog")
	for _, post := range data.Posts {
		p := post
		if err := g.RenderPage(blogDistDir, post.Slug+".html", "post.html", post.Title, PageData{
			Post:       &p,
			PathPrefix: "../",
		}); err != nil {
			return err
		}
	}
	return nil
}

func (g *SiteGenerator) GenerateSearchIndex(distDir string, data *ContentData) error {
	var items []SearchItem
	for _, post := range data.Posts {
		items = append(items, SearchItem{
			Title:       post.Title,
			Slug:        post.Slug,
			Description: post.Description,
			Date:        post.Date.Format("January 02, 2006"),
			Tags:        post.Tags,
		})
	}

	jsonData, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("failed to marshal search index: %w", err)
	}

	if err := os.WriteFile(filepath.Join(distDir, "search-index.json"), jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write search-index.json: %w", err)
	}
	return nil
}

func (g *SiteGenerator) GenerateRSS(distDir string, posts []Post) error {
	f, err := os.Create(filepath.Join(distDir, "rss.xml"))
	if err != nil {
		return fmt.Errorf("failed to create rss.xml: %w", err)
	}
	defer f.Close()

	escape := func(s string) string {
		var b strings.Builder
		if err := xml.EscapeText(&b, []byte(s)); err != nil {
			return s
		}
		return b.String()
	}

	if _, err := fmt.Fprint(f, `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
<channel>
  <title>`+escape(g.Config.Landing.Title)+`</title>
  <link>`+g.Config.Landing.URL+`</link>
  <description>`+escape(g.Config.Landing.Slogan)+`</description>
  <language>en-us</language>
`); err != nil {
		return err
	}

	for _, post := range posts {
		link := g.Config.Landing.URL + "blog/" + post.Slug + ".html"
		if _, err := fmt.Fprintf(f, `  <item>
    <title>%s</title>
    <link>%s</link>
    <description>%s</description>
    <pubDate>%s</pubDate>
    <guid>%s</guid>
  </item>
`, escape(post.Title), link, escape(post.Description), post.Date.Format(time.RFC1123), link); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(f, `</channel>
</rss>`); err != nil {
		return err
	}
	return nil
}

func (g *SiteGenerator) GenerateSitemap(distDir string, posts []Post) error {
	f, err := os.Create(filepath.Join(distDir, "sitemap.xml"))
	if err != nil {
		return fmt.Errorf("failed to create sitemap.xml: %w", err)
	}
	defer f.Close()

	if _, err := fmt.Fprint(f, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
`); err != nil {
		return err
	}

	// Static Pages
	pages := []string{"", "about.html", "blog.html", "archive.html"}
	for _, page := range pages {
		if _, err := fmt.Fprintf(f, `  <url>
    <loc>%s%s</loc>
    <lastmod>%s</lastmod>
  </url>
`, g.Config.Landing.URL, page, time.Now().Format("2006-01-02")); err != nil {
			return err
		}
	}

	// Blog Posts
	for _, post := range posts {
		if _, err := fmt.Fprintf(f, `  <url>
    <loc>%sblog/%s.html</loc>
    <lastmod>%s</lastmod>
  </url>
`, g.Config.Landing.URL, post.Slug, post.Date.Format("2006-01-02")); err != nil {
			return err
		}
	}

	// API Manifest
	if _, err := fmt.Fprintf(f, `  <url>
    <loc>%sapi/manifest.json</loc>
    <lastmod>%s</lastmod>
  </url>
`, g.Config.Landing.URL, time.Now().Format("2006-01-02")); err != nil {
		return err
	}

	if _, err := fmt.Fprint(f, `</urlset>`); err != nil {
		return err
	}
	return nil
}

func (g *SiteGenerator) GenerateRegistries(distDir string, data *ContentData) error {
	apiDir := filepath.Join(distDir, "api")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return fmt.Errorf("failed to create api dir %s: %w", apiDir, err)
	}

	// Blog Items
	var blogItems []BlogItem
	for _, post := range data.Posts {
		blogItems = append(blogItems, BlogItem{
			Title:       post.Title,
			Description: post.Description,
			URL:         g.Config.Landing.URL + "blog/" + post.Slug + ".html",
			Date:        post.Date.Format(time.RFC3339),
			Tags:        post.Tags,
		})
	}

	// Projects
	var projectItems []ProjectItem
	for _, p := range g.Config.Projects {
		techs := g.FuncMap["cleanYAMLList"].(func(interface{}) []string)(p.Techs)
		projectItems = append(projectItems, ProjectItem{
			Title:            p.Title,
			ShortDescription: p.ShortDescription,
			Link:             p.Link,
			Techs:            techs,
		})
	}

	// Skills
	var allSkills []string
	for _, s := range g.Config.Skills {
		allSkills = append(allSkills, s.Name)
	}

	// Unified MCP Manifest
	manifest := Manifest{
		MCPVersion: "1.0",
		Name:       g.Config.Landing.Title,
		URL:        g.Config.Landing.URL,
		UpdatedAt:  time.Now().Format(time.RFC3339),
		Profile: ProfileRegistry{
			URL:        g.Config.Landing.URL,
			Title:      g.Config.Landing.Title,
			Name:       g.Config.Landing.Name,
			Slogan:     g.Config.Landing.Slogan,
			Experience: g.Config.Landing.Experience,
			Status:     g.Config.Landing.Status,
			FocusAreas: g.Config.Landing.FocusAreas,
			About: ProfileAbout{
				Timeline:    g.Config.About.Timeline,
				LastUpdated: g.Config.About.LastUpdated,
				Currently:   g.Config.About.Currently,
			},
		},
		Skills:   allSkills,
		Projects: projectItems,
		Blog: BlogRegistry{
			TotalPosts: len(data.Posts),
			Posts:      blogItems,
		},
	}
	return g.writeJSON(filepath.Join(apiDir, "manifest.json"), manifest)
}

func (g *SiteGenerator) Build(distDir string, data *ContentData) error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"static pages", func() error { return g.GenerateStaticPages(distDir, data) }},
		{"blog pagination", func() error { return g.GenerateBlogPagination(distDir, data, 10) }},
		{"tag pages", func() error { return g.GenerateTagPages(distDir, data) }},
		{"post pages", func() error { return g.GeneratePostPages(distDir, data) }},
		{"search index", func() error { return g.GenerateSearchIndex(distDir, data) }},
		{"registries", func() error { return g.GenerateRegistries(distDir, data) }},
		{"llms.txt", func() error { return g.GenerateLLMsTxt(distDir) }},
		{"RSS", func() error { return g.GenerateRSS(distDir, data.Posts) }},
		{"sitemap", func() error { return g.GenerateSitemap(distDir, data.Posts) }},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			return fmt.Errorf("failed to generate %s: %w", step.name, err)
		}
	}

	if g.totalOriginalSize > 0 {
		beforeMB := float64(g.totalOriginalSize) / 1000000.0
		afterMB := float64(g.totalMinifiedSize) / 1000000.0
		savings := float64(g.totalOriginalSize-g.totalMinifiedSize) / float64(g.totalOriginalSize) * 100
		fmt.Printf("HTML minification:\nBefore: %d bytes (%.2f MB)\nAfter: %d bytes (%.2f MB)\nSavings: %.1f%%\n", g.totalOriginalSize, beforeMB, g.totalMinifiedSize, afterMB, savings)
	}

	return nil
}

func (g *SiteGenerator) writeJSON(path string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for %s: %w", path, err)
	}
	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON to %s: %w", path, err)
	}
	return nil
}

func (g *SiteGenerator) GenerateLLMsTxt(distDir string) error {
	var sb strings.Builder

	// Role & Identity
	sb.WriteString("# " + g.Config.Landing.Name + " - Technical Portfolio\n\n")
	sb.WriteString("## Role & Identity\n\n")
	sb.WriteString(g.Config.Landing.Slogan + "\n\n")

	// Recruiting Signals
	sb.WriteString("## Recruiting Signals\n\n")
	if g.Config.Landing.Status != "" {
		sb.WriteString("- **Status**: " + g.Config.Landing.Status + "\n")
	}
	sb.WriteString("- **Experience**: " + g.Config.Landing.Experience + "\n")
	sb.WriteString("- **Focus Areas**: " + strings.Join(g.Config.Landing.FocusAreas, ", ") + "\n\n")

	// Technical Skills
	sb.WriteString("## Technical Skills\n\n")
	var allSkills []string
	for _, s := range g.Config.Skills {
		allSkills = append(allSkills, s.Name)
	}
	sb.WriteString(strings.Join(allSkills, ", ") + "\n\n")

	// Project Index
	sb.WriteString("## Project Index (Discovery)\n\n")
	for _, p := range g.Config.Projects {
		sb.WriteString("- **" + p.Title + "**: " + p.ShortDescription + "\n")
	}
	sb.WriteString("\n")

	// Discovery Registry
	sb.WriteString("## Discovery Registry\n\n")
	sb.WriteString("The following endpoint provides unified technical context for AI agents (Model Context Protocol):\n\n")
	sb.WriteString("- **Unified Manifest**: " + g.Config.Landing.URL + "api/manifest.json\n\n")

	// Contact
	sb.WriteString("## Contact\n\n")
	for _, s := range g.Config.Socials {
		if strings.EqualFold(s.Name, "GitHub") || strings.EqualFold(s.Name, "LinkedIn") {
			sb.WriteString("- **" + s.Name + "**: " + s.Href + "\n")
		}
	}

	content := []byte(sb.String())

	// Write to root for LLM discovery
	if err := os.WriteFile(filepath.Join(distDir, "llms.txt"), content, 0644); err != nil {
		return fmt.Errorf("failed to write root llms.txt: %w", err)
	}

	return nil
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
func CopyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return CopyFile(path, targetPath)
	})
}

// CopyFile copies a single file from src to dst, preserving file permissions.
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}
