package generator

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mehub/internal/config"
	"mehub/internal/content"
)

type PageData struct {
	Config       *config.SiteConfig
	CurrentYear  int
	Title        string
	Posts        []content.Post
	Post         *content.Post
	Tags         []string
	TagCounts    map[string]int
	Archive      map[int][]content.Post
	ArchiveYears []int
	CurrentPage  int
	TotalPages   int
	PathPrefix   string
}

type SiteGenerator struct {
	Config  *config.SiteConfig
	FuncMap template.FuncMap
}

func New(cfg *config.SiteConfig) *SiteGenerator {
	return &SiteGenerator{
		Config: cfg,
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

	tmpl, err := template.New("base.html").Funcs(g.FuncMap).ParseFiles("internal/templates/base.html", tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse templates for %s: %w", tmplPath, err)
	}

	outputFile, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", filename, err)
	}
	defer outputFile.Close()

	title := g.Config.Site.Title
	if titlePrefix != "" {
		title = titlePrefix + " | " + title
	}

	data.Config = g.Config
	data.CurrentYear = time.Now().Year()
	data.Title = title

	if err := tmpl.ExecuteTemplate(outputFile, "base.html", data); err != nil {
		return fmt.Errorf("failed to execute template for %s: %w", tmplPath, err)
	}
	return nil
}

func (g *SiteGenerator) GenerateStaticPages(distDir string, data *content.ContentData) error {
	pages := []struct {
		filename    string
		tmplPath    string
		titlePrefix string
		data        PageData
	}{
		{"index.html", "internal/templates/index.html", "", PageData{}},
		{"about.html", "internal/templates/about.html", "About", PageData{}},
		{"now.html", "internal/templates/now.html", "Now", PageData{}},
		{"404.html", "internal/templates/404.html", "404 - Not Found", PageData{}},
		{"tags.html", "internal/templates/tags.html", "Tags", PageData{Tags: data.Tags, TagCounts: data.TagCounts}},
		{"archive.html", "internal/templates/archive.html", "Archive", PageData{Archive: data.PostsByYear, ArchiveYears: data.ArchiveYears}},
	}

	for _, p := range pages {
		if err := g.RenderPage(distDir, p.filename, p.tmplPath, p.titlePrefix, p.data); err != nil {
			return err
		}
	}
	return nil
}

func (g *SiteGenerator) GenerateBlogPagination(distDir string, data *content.ContentData, pageSize int) error {
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
			if err := g.RenderPage(distDir, "blog.html", "internal/templates/blog.html", "Blog", PageData{
				Posts:       pagePosts,
				CurrentPage: pageNumber,
				TotalPages:  totalPages,
			}); err != nil {
				return err
			}
		} else {
			pageDir := filepath.Join(distDir, "blog")
			if err := g.RenderPage(pageDir, fmt.Sprintf("%d.html", pageNumber), "internal/templates/blog.html", fmt.Sprintf("Blog - Page %d", pageNumber), PageData{
				Posts:       pagePosts,
				CurrentPage: pageNumber,
				TotalPages:  totalPages,
				PathPrefix:  "../",
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *SiteGenerator) GenerateTagPages(distDir string, data *content.ContentData) error {
	tagsDistDir := filepath.Join(distDir, "tags")
	for tag, tagPosts := range data.PostsByTag {
		if err := g.RenderPage(tagsDistDir, tag+".html", "internal/templates/blog.html", "#"+tag, PageData{
			Posts:      tagPosts,
			PathPrefix: "../",
		}); err != nil {
			return err
		}
	}
	return nil
}

func (g *SiteGenerator) GeneratePostPages(distDir string, data *content.ContentData) error {
	blogDistDir := filepath.Join(distDir, "blog")
	for _, post := range data.Posts {
		p := post
		if err := g.RenderPage(blogDistDir, post.Slug+".html", "internal/templates/post.html", post.Title, PageData{
			Post:       &p,
			PathPrefix: "../",
		}); err != nil {
			return err
		}
	}
	return nil
}

func (g *SiteGenerator) GenerateSearchIndex(distDir string, data *content.ContentData) error {
	type SearchItem struct {
		Title       string   `json:"title"`
		Slug        string   `json:"slug"`
		Description string   `json:"description"`
		Date        string   `json:"date"`
		Tags        []string `json:"tags"`
	}

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

func (g *SiteGenerator) GenerateRSS(distDir string, posts []content.Post) error {
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
  <title>`+escape(g.Config.Site.Title)+`</title>
  <link>`+g.Config.Site.URL+`</link>
  <description>`+escape(g.Config.Site.Description)+`</description>
  <language>en-us</language>
`); err != nil {
		return err
	}

	for _, post := range posts {
		link := g.Config.Site.URL + "blog/" + post.Slug + ".html"
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

func (g *SiteGenerator) GenerateSitemap(distDir string, posts []content.Post) error {
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
	pages := []string{"", "about.html", "now.html", "blog.html", "tags.html", "archive.html"}
	for _, page := range pages {
		if _, err := fmt.Fprintf(f, `  <url>
    <loc>%s%s</loc>
    <lastmod>%s</lastmod>
  </url>
`, g.Config.Site.URL, page, time.Now().Format("2006-01-02")); err != nil {
			return err
		}
	}

	// Blog Posts
	for _, post := range posts {
		if _, err := fmt.Fprintf(f, `  <url>
    <loc>%sblog/%s.html</loc>
    <lastmod>%s</lastmod>
  </url>
`, g.Config.Site.URL, post.Slug, post.Date.Format("2006-01-02")); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(f, `</urlset>`); err != nil {
		return err
	}
	return nil
}

func (g *SiteGenerator) GenerateRegistryEndpoint(distDir string, data *content.ContentData) error {
	apiDir := filepath.Join(distDir, "api")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return fmt.Errorf("failed to create api dir: %w", err)
	}

	type RegistryItem struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		URL         string   `json:"url"`
		Date        string   `json:"date_published"`
		Tags        []string `json:"skills"`
	}

	var items []RegistryItem
	for _, post := range data.Posts {
		items = append(items, RegistryItem{
			Title:       post.Title,
			Description: post.Description,
			URL:         g.Config.Site.URL + "blog/" + post.Slug + ".html",
			Date:        post.Date.Format(time.RFC3339),
			Tags:        post.Tags,
		})
	}

	jsonData, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	if err := os.WriteFile(filepath.Join(apiDir, "skills-registry.json"), jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write skills-registry.json: %w", err)
	}
	return nil
}
